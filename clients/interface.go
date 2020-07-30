package clients

import (
	"net"

	"github.com/rikaaa0928/awsl/model"
)

// Client listen and handel incomming
type Client interface {
	Dial(model.ANetAddr) (net.Conn, error)
	Verify(net.Conn) error
	IDTag() (int, string)
}
