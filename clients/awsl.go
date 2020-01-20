package clients

import (
	"crypto/tls"
	"encoding/json"
	"github.com/Evi1/awsl/servers"
	"golang.org/x/net/websocket"
	"log"
	"net"
)

func NewAWSL(serverHost, serverPort, uri string) AWSL {
	return AWSL{
		ServerHost: serverHost,
		ServerPort: serverPort,
		Uri:        uri,
	}
}

// AWSL AWSL
type AWSL struct {
	ServerHost string
	ServerPort string
	Uri        string
}

// Dial Dial
func (c AWSL) Dial(addr servers.ANetAddr) (net.Conn, error) {
	config, err := websocket.NewConfig("wss://"+c.ServerHost+":"+c.ServerPort+"/"+c.Uri, "wss://"+c.ServerHost+":"+c.ServerPort+"/"+c.Uri)
	if err != nil {
		log.Println("conf:" + err.Error())
		return nil, err
	}
	config.TlsConfig = &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         addr.Host,
	}
	ws, err := websocket.DialConfig(config)
	// ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Println("dial:" + err.Error())
		return ws, err
	}
	addrBytes, err := json.Marshal(addr)
	if err != nil {
		log.Println("dial json:" + err.Error())
		return ws, err
	}
	_, err = ws.Write(addrBytes)
	return ws, err
}

// Verify Verify
func (c AWSL) Verify(conn net.Conn) error {
	return nil
}
