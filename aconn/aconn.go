package aconn

import "net"

type AConn interface {
	net.Conn
	EndAddr() net.Addr
	SetEndAddr(addr net.Addr)
}
