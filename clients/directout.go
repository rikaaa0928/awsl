package clients

import "net"

// DirectOut DirectOut
type DirectOut struct{}

// Dial Dial
func (c DirectOut) Dial(h string, p string) (net.Conn, error) {
	return net.Dial("tcp", h+":"+p)
}

// Verify Verify
func (c DirectOut) Verify(_ net.Conn) error {
	return nil
}
