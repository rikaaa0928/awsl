package adialer

import (
	"context"
	"log"
	"net"
	"strings"
	"time"

	"github.com/rikaaa0928/awsl/aconn"
	"github.com/rikaaa0928/awsl/global"
)

var Free = func(ctx context.Context, addr net.Addr) (context.Context, aconn.AConn, error) {
	ctx = context.WithValue(ctx, global.CTXOutType, "free")
	if strings.ToLower(addr.Network()) == "udp" {
		log.Println("dial udp : " + addr.String())
		uDst, err := net.ResolveUDPAddr("udp", addr.String())
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
		//ac := aconn.NewAConn(&udpConnWrapper{UDPConn: c, toAddr: uDst, src: src, dst: dst, lAddr: luAddr})
		ac := aconn.NewAConn(&udpConnWrapper{UDPConn: c, toAddr: uDst, lAddr: luAddr})
		ac.SetEndAddr(uDst)
		return ctx, ac, err
	}
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
		ac := aconn.NewAConn(&udpConnWrapper{UDPConn: c, lAddr: luAddr, toAddr: uDst})
		ac.SetEndAddr(uDst)
		return ctx, ac, err
	}
}

type udpConnWrapper struct {
	// sync.Mutex
	*net.UDPConn
	toAddr *net.UDPAddr
	lAddr  *net.UDPAddr
	// src    string
	// dst    string
}

func (c *udpConnWrapper) Read(b []byte) (n int, err error) {
	var rAddr *net.UDPAddr
	c.UDPConn.SetReadDeadline(time.Now().Add(time.Minute * 10))
	n, rAddr, err = c.UDPConn.ReadFromUDP(b)
	if rAddr.String() != c.toAddr.String() {
		log.Println("free readFromUDP addr not match", rAddr, c.toAddr)
	}
	return
}

func (c *udpConnWrapper) Write(b []byte) (n int, err error) {
	c.UDPConn.SetWriteDeadline(time.Now().Add(time.Minute))
	n, err = c.UDPConn.Write(b)
	return
}
