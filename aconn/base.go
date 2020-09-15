package aconn

import (
	"net"
	"strconv"
)

func NewAConn(c net.Conn) AConn {
	return &BaseConn{
		Conn: c,
		End:  nil,
	}
}

type addrInfo struct {
	host    string
	port    int
	network string
}

func (a addrInfo) Network() string {
	return a.network
}

func (a addrInfo) String() string {
	return net.JoinHostPort(a.host, strconv.Itoa(a.port))
}

type BaseConn struct {
	net.Conn
	End net.Addr
}

func (c *BaseConn) EndAddr() net.Addr {
	return c.End
}

func (c *BaseConn) SetEndAddr(eh string, ep int, en string) {
	c.End = addrInfo{
		host:    eh,
		port:    ep,
		network: en,
	}
}
