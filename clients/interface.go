package clients

import (
	"github.com/Evi1/awsl/servers"
	"net"
)

// Client listen and handel incomming
type Client interface {
	Dial(servers.ANetAddr) (net.Conn, error)
	Verify(net.Conn) error
}
