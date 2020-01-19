package clients

import (
	"crypto/tls"
	"log"
	"net"

	"golang.org/x/net/websocket"
)

// AWSL AWSL
type AWSL struct {
}

// Dial Dial
func (c AWSL) Dial(h string, p string) (net.Conn, error) {
	config, err := websocket.NewConfig("wss://"+h+":"+p+"/wss", "wss://"+h+":"+p+"/wss")
	if err != nil {
		log.Println("conf:" + err.Error())
		return nil, err
	}
	config.TlsConfig = &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         h,
	}
	ws, err := websocket.DialConfig(config)
	// ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Println("dial:" + err.Error())
	}
	return ws, err
}

// Verify Verify
func (c AWSL) Verify(_ net.Conn) error {
	return nil
}
