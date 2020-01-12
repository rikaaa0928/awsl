package clients

import "net"

type DirectOut struct{}

func (c DirectOut) Dail(h string, p string) (net.Conn, error) {
	return net.Dial("tcp", h+":"+p)
}

func (c DirectOut) Verify(_ net.Conn) error {
	return nil
}
