package clients

import (
	"github.com/Evi1/awsl/servers"
	"net"
	"strconv"
)

// DirectOut DirectOut
type DirectOut struct{}

// Dial Dial
func (c DirectOut) Dial(addr servers.ANetAddr) (net.Conn, error) {
	return net.Dial("tcp", addr.Host+":"+strconv.Itoa(addr.Port))
}

// Verify Verify
func (c DirectOut) Verify(_ net.Conn) error {
	return nil
}
