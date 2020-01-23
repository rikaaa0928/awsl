package servers

import (
	"net"

	"github.com/Evi1/awsl/model"
)

// Server listen and handel incomming
type Server interface {
	Listen() net.Listener
	ReadRemote(net.Conn) (model.ANetAddr, error)
}
