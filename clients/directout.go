package clients

import (
	"net"
	"strconv"

	"github.com/rikaaa0928/awsl/model"
)

// DirectOut DirectOut
type DirectOut struct {
	id  int
	tag string
}

// Dial Dial
func (c DirectOut) Dial(addr model.ANetAddr) (net.Conn, error) {
	network := "tcp"
	if addr.CMD == model.UDP {
		network = "udp"
	}
	return net.Dial(network, addr.Host+":"+strconv.Itoa(addr.Port))
}

// Verify Verify
func (c DirectOut) Verify(_ net.Conn) error {
	return nil
}

// IDTag id tag
func (c DirectOut) IDTag() (int, string) {
	return c.id, c.tag
}
