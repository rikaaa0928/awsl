package tools

import (
	"log"
	"net"
)

// Send send
func Send(c net.Conn, d []byte) int {
	n, e := c.Write(d)
	if e != nil {
		log.Println(e)
		return n
	}
	return n
}

// Receive receive
func Receive(c net.Conn) (int, []byte) {
	buf := make([]byte, 1024)
	n, err := c.Read(buf)
	if err != nil {
		log.Println(err)
	}
	return n, buf[:n]
}
