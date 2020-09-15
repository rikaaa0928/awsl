package aconn

import "net"

type AConn interface {
	net.Conn
	EndAddr() net.Addr
	SetEndAddr(eh string, ep int, en string)
}
