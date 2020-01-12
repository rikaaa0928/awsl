package servers

import "net"

// Server listen and handel incomming
type Server interface {
	Listen() net.Listener
	ReadRemote(net.Conn) (string, string, error)
}
