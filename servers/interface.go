package servers

import "net"

type ANetAddr struct {
	Type int //4 6 1
	Host string
	Port int
}

// Server listen and handel incomming
type Server interface {
	Listen() net.Listener
	ReadRemote(net.Conn) (ANetAddr, error)
}
