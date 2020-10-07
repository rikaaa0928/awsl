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

func (a *AddrInfo) Parse(network, str string) error {
	a.NetName = network
	h, p, e := net.SplitHostPort(str)
	if e != nil {
		return e
	}
	a.Host = h
	a.Port, e = strconv.Atoi(p)
	if e != nil {
		return e
	}
	// l := strings.Split(str, ":")
	// if len(l) < 2 {
	// 	return errors.New("error host str format(host:port): " + str)
	// }
	// a.Port, err = strconv.Atoi(l[len(l)-1])
	// if err != nil {
	// 	return
	// }
	// if len(a.Host) <= 0 {
	// 	return errors.New("wrong host: " + a.Host + " str: " + str)
	// }
	// a.Host = strings.Join(l[:len(l)-1], "")
	// if strings.Index(a.Host, ":") >= 0 {
	// 	a.Host = strings.Trim(a.Host, "[]")
	// }
	return nil
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
