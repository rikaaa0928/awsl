package adialer

import (
	"context"
	"encoding/json"
	"net"
	"time"

	"github.com/rikaaa0928/awsl/aconn"
	"github.com/rikaaa0928/awsl/consts"
	"github.com/rikaaa0928/awsl/utils"
)

var FreeTCP = func(ctx context.Context, addr net.Addr) (context.Context, aconn.AConn, error) {
	c, err := net.Dial("tcp", addr.String())
	ac := aconn.NewAConn(c)
	ac.SetEndAddr(addr)
	return ctx, ac, err
}

var FreeUDP = func(ctx context.Context, src, dst string) (context.Context, aconn.AConn, error) {
	uDst, err := net.ResolveUDPAddr("udp", dst)
	if err != nil {
		return ctx, nil, err
	}
	uSrc, err := net.ResolveUDPAddr("udp", src)
	if err != nil {
		return ctx, nil, err
	}
	c, err := net.DialUDP("udp", nil, uDst)
	ac := aconn.NewAConn(udpConnWrapper{c: c, a: uSrc, src: src, dst: dst})
	ac.SetEndAddr(uDst)
	return ctx, ac, err
}

type udpConnWrapper struct {
	c   *net.UDPConn
	a   *net.UDPAddr
	src string
	dst string
}

func (c udpConnWrapper) Read(b []byte) (n int, err error) {
	buf := utils.GetMem(65536)
	defer utils.PutMem(buf)
	var dstAddr *net.UDPAddr
	n, dstAddr, err = c.c.ReadFromUDP(buf)
	if err != nil {
		return
	}
	udp := consts.UDPMSG{
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

func (c udpConnWrapper) Write(b []byte) (n int, err error) {

}

func (c udpConnWrapper) Close() error {
	return c.c.Close()
}

func (c udpConnWrapper) LocalAddr() net.Addr {

}

// RemoteAddr returns the remote network address.
func (c udpConnWrapper) RemoteAddr() net.Addr {

}

// SetDeadline sets the read and write deadlines associated
// with the connection. It is equivalent to calling both
// SetReadDeadline and SetWriteDeadline.
//
// A deadline is an absolute time after which I/O operations
// fail instead of blocking. The deadline applies to all future
// and pending I/O, not just the immediately following call to
// Read or Write. After a deadline has been exceeded, the
// connection can be refreshed by setting a deadline in the future.
//
// If the deadline is exceeded a call to Read or Write or to other
// I/O methods will return an error that wraps os.ErrDeadlineExceeded.
// This can be tested using errors.Is(err, os.ErrDeadlineExceeded).
// The error's Timeout method will return true, but note that there
// are other possible errors for which the Timeout method will
// return true even if the deadline has not been exceeded.
//
// An idle timeout can be implemented by repeatedly extending
// the deadline after successful Read or Write calls.
//
// A zero value for t means I/O operations will not time out.
func (c udpConnWrapper) SetDeadline(t time.Time) error {

}

// SetReadDeadline sets the deadline for future Read calls
// and any currently-blocked Read call.
// A zero value for t means Read will not time out.
func (c udpConnWrapper) SetReadDeadline(t time.Time) error {

}

// SetWriteDeadline sets the deadline for future Write calls
// and any currently-blocked Write call.
// Even if write times out, it may return n > 0, indicating that
// some of the data was successfully written.
// A zero value for t means Write will not time out.
func (c udpConnWrapper) SetWriteDeadline(t time.Time) error {

}
