package clients

import "net"

// Client listen and handel incomming
type Client interface {
	Dail(h string, p string) (net.Conn, error)
	Verify(net.Conn) error
}
