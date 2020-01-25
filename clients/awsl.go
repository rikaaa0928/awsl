package clients

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net"

	"github.com/Evi1/awsl/model"
	"golang.org/x/net/websocket"
)

// NewAWSL NewAWSL
func NewAWSL(serverHost, serverPort, uri string) AWSL {
	return AWSL{
		ServerHost: serverHost,
		ServerPort: serverPort,
		URI:        uri,
	}
}

// AWSL AWSL
type AWSL struct {
	ServerHost string
	ServerPort string
	URI        string
}

// Dial Dial
func (c AWSL) Dial(addr model.ANetAddr) (net.Conn, error) {
	config, err := websocket.NewConfig("wss://"+c.ServerHost+":"+c.ServerPort+"/"+c.URI, "https://"+c.ServerHost+":"+c.ServerPort+"/")
	if err != nil {
		log.Println("conf:" + err.Error())
		return nil, err
	}
	config.TlsConfig = &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         c.ServerHost,
	}
	ws, err := websocket.DialConfig(config)
	// ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Println("dial:" + err.Error())
		return ws, err
	}
	conn := awslConn{Conn: ws, Addr: addr}
	return conn, err
}

// Verify Verify
func (c AWSL) Verify(conn net.Conn) error {
	ws := conn.(awslConn)
	addrBytes, err := json.Marshal(ws.Addr)
	if err != nil {
		log.Println("json marshal:" + err.Error())
		return err
	}
	_, err = ws.Write(addrBytes)
	return err
}

type awslConn struct {
	*websocket.Conn
	Addr model.ANetAddr
}
