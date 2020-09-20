package aconn

import "net"

type AConn interface {
	net.Conn
	EndAddr() net.Addr
	SetEndAddr(addr net.Addr)
}

type AConnMidor interface {
	RegisterCloser(mid CLoserMid)
	RegisterReader(mid IOMid)
	RegisterWriter(mid IOMid)
}
