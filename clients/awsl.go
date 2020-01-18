package clients

import "net"

// AWSL AWSL
type AWSL struct{}

// Dial Dial
func (c AWSL) Dial(h string, p string) (net.Conn, error) {
	return nil, nil
}

// Verify Verify
func (c AWSL) Verify(_ net.Conn) error {
	return nil
}
