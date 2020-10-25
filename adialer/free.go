package adialer

import (
	"context"
	"encoding/json"
	"net"
	"sync"
	"time"

	"github.com/rikaaa0928/awsl/aconn"
	"github.com/rikaaa0928/awsl/utils"
	"github.com/rikaaa0928/awsl/utils/superlib"
)

var FreeTCP = func(ctx context.Context, addr net.Addr) (context.Context, aconn.AConn, error) {
	c, err := net.Dial("tcp", addr.String())
	ac := aconn.NewAConn(c)
	ac.SetEndAddr(addr)
	return ctx, ac, err
}

func NewFreeUDP(src, dst string) ADialer {
	return func(ctx context.Context, _ net.Addr) (context.Context, aconn.AConn, error) {
		uDst, err := net.ResolveUDPAddr("udp", dst)
		if err != nil {
			return ctx, nil, err
		}
		//uSrc, err := net.ResolveUDPAddr("udp", src)
		//if err != nil {
		//	return ctx, nil, err
		//}
		c, err := net.DialUDP("udp", nil, uDst)
		if err != nil {
			return ctx, nil, err
		}
		lAddr := c.LocalAddr()
		luAddr, _ := net.ResolveUDPAddr(lAddr.Network(), lAddr.String())
		ac := aconn.NewAConn(&udpConnWrapper{UDPConn: c, toAddr: uDst, src: src, dst: dst, lAddr: luAddr})
		ac.SetEndAddr(uDst)
		return ctx, ac, err
	}
}

type udpConnWrapper struct {
	sync.Mutex
	*net.UDPConn
	toAddr *net.UDPAddr
	lAddr  *net.UDPAddr
	src    string
	dst    string
}

func (c *udpConnWrapper) reDial() error {
	c.Lock()
	defer c.Unlock()
	c2, err := net.DialUDP("udp", c.lAddr, c.toAddr)
	if err != nil {
		return err
	}
	c.UDPConn = c2
	return nil
}

func (c *udpConnWrapper) Read(b []byte) (n int, err error) {
	buf := utils.GetMem(65536)
	defer utils.PutMem(buf)
	var dstAddr *net.UDPAddr
	c.UDPConn.SetReadDeadline(time.Now().Add(time.Minute))
	n, dstAddr, err = c.UDPConn.ReadFromUDP(buf)
	i := 0
	for err != nil && i < 3 {
		err = c.reDial()
		if err != nil {
			continue
		}
		c.UDPConn.SetReadDeadline(time.Now().Add(time.Minute))
		n, dstAddr, err = c.UDPConn.ReadFromUDP(buf)
		if err != nil {
			break
		}
		i++
	}
	if err != nil {
		return
	}
	udp := superlib.UDPMSG{
		DstStr: dstAddr.String(),
		SrcStr: c.src,
		Data:   buf[:n],
	}
	str, err := json.Marshal(udp)
	if err != nil {
		return
	}
	n = len(str)
	copy(b, str)
	return
}

func (c *udpConnWrapper) Write(b []byte) (n int, err error) {
	var udpMsg superlib.UDPMSG
	err = json.Unmarshal(b, &udpMsg)
	if err != nil {
		return -1, err
	}
	n, err = c.UDPConn.Write(udpMsg.Data)
	i := 0
	for err != nil && i < 3 {
		err = c.reDial()
		if err != nil {
			time.Sleep(time.Second * time.Duration(i))
			continue
		}
		n, err = c.UDPConn.Write(udpMsg.Data)
		if err != nil {
			break
		}
		i++
	}
	return
}
