package clients

import (
	"crypto/tls"
	"log"
	"net"

	"golang.org/x/net/websocket"
)

// AWSL AWSL
type AWSL struct {
	Host     string
	HostName string
}

// Dial Dial
func (c AWSL) Dial(h string, p string) (net.Conn, error) {
	config, _ := websocket.NewConfig(c.Host, "")
	config.TlsConfig = &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         c.HostName,
	}
	ws, err := websocket.DialConfig(config)
	// ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	return ws, nil
}

// Verify Verify
func (c AWSL) Verify(_ net.Conn) error {
	return nil
}
