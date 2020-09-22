package aconn

import (
	"net"
	"strconv"
	"strings"
)

func NewAConn(c net.Conn) AConn {
	return &BaseConn{
		Conn: c,
		End:  nil,
	}
}

func NewAddr(eh string, ep int, en string) net.Addr {
	return AddrInfo{
		Host:    eh,
		Port:    ep,
		NetName: en,
	}
}

type AddrInfo struct {
	Host    string
	Port    int
	NetName string
}

func (a AddrInfo) Network() string {
	return a.NetName
}

func (a AddrInfo) String() string {
	return net.JoinHostPort(a.Host, strconv.Itoa(a.Port))
}

func (a *AddrInfo) Parse(network, str string) (err error) {
	a.NetName = network
	l := strings.Split(str, ":")
	a.Port, err = strconv.Atoi(l[len(l)-1])
	if err != nil {
		return
	}
	a.Host = strings.Join(l[:len(l)-1], "")
	return
}

type BaseConn struct {
	net.Conn
	End net.Addr
}

func (c *BaseConn) EndAddr() net.Addr {
	return c.End
}

func (c *BaseConn) SetEndAddr(addr net.Addr) {
	c.End = addr
}
