package clients

import (
	"net"

	"github.com/Evi1/awsl/model"
)

// Client listen and handel incomming
type Client interface {
	Dial(model.ANetAddr) (net.Conn, error)
	Verify(net.Conn) error
	IDTag() (int, string)
}
