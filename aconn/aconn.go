package aconn

import "net"

type AConn interface {
	net.Conn
	EndAddr() net.Addr
	SetEndAddr(addr net.Addr)
	Magic() *uint32
	SetMagic(uint32)
}

type MidsMgr interface {
	RegisterCloser(mid CLoserMid)
	RegisterReader(mid IOMid)
	RegisterWriter(mid IOMid)
}
