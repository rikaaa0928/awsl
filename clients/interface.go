package clients

import (
	"github.com/Evi1/awsl/model"
	"net"
)

// Client listen and handel incomming
type Client interface {
	Dial(model.ANetAddr) (net.Conn, error)
	Verify(net.Conn) error
}
