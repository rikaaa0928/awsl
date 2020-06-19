package servers

import (
	"errors"
	"net"

	"github.com/Evi1/awsl/model"
)

// ErrUDP ErrUDP
var ErrUDP = errors.New("udp error")

// Server listen and handel incomming
type Server interface {
	Listen() net.Listener
	ReadRemote(net.Conn) (model.ANetAddr, error)
	IDTag() (int, string)
}
