package tools

import "net"

// Receive receive
func Receive(c net.Conn, buf []byte) (int, []byte, error) {
	n, err := c.Read(buf)
	return n, buf[:n], err
}
