package rwconn

import (
	"io"
	"net"
	"time"
)

func NewRWConn(w io.WriteCloser, r io.ReadCloser) *RWConn {
	return &RWConn{
		w: w,
		r: r,
	}
}

type RWConn struct {
	w io.WriteCloser
	r io.ReadCloser
}

func (c *RWConn) Read(b []byte) (n int, err error) {
	return c.r.Read(b)
}

func (c *RWConn) Write(b []byte) (n int, err error) {
	return c.w.Write(b)
}

func (c *RWConn) Close() error {
	c.w.Close()
	return c.r.Close()
}
func (c *RWConn) LocalAddr() net.Addr {
	return nil
}
func (c *RWConn) RemoteAddr() net.Addr {
	return nil
}
func (c *RWConn) SetDeadline(t time.Time) error {
	return nil
}
func (c *RWConn) SetReadDeadline(t time.Time) error {
	return nil
}
func (c *RWConn) SetWriteDeadline(t time.Time) error {
	return nil
}
