package clients

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"net"

	"github.com/Evi1/awsl/config"
	"github.com/Evi1/awsl/model"
	"golang.org/x/net/websocket"
)

// NewAWSL NewAWSL
func NewAWSL(serverHost, serverPort, uri, auth string) *AWSL {
	return &AWSL{
		ServerHost: serverHost,
		ServerPort: serverPort,
		URI:        uri,
		Auth:       auth,
	}
}

// AWSL AWSL
type AWSL struct {
	ServerHost string
	ServerPort string
	URI        string
	Auth       string
}

// Dial Dial
func (c *AWSL) Dial(addr model.ANetAddr) (net.Conn, error) {
	wsConfig, err := websocket.NewConfig("wss://"+c.ServerHost+":"+c.ServerPort+"/"+c.URI, "https://"+c.ServerHost+":"+c.ServerPort+"/")
	if err != nil {
		//log.Println("conf:" + err.Error())
		return nil, err
	}
	wsConfig.TlsConfig = &tls.Config{
		InsecureSkipVerify: config.GetConf().NoVerify,
		ServerName:         c.ServerHost,
	}
	ws, err := websocket.DialConfig(wsConfig)
	// ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		//log.Println("dial:" + err.Error())
		return ws, err
	}
	conn := awslConn{Conn: ws, Addr: addr}
	return conn, err
}

// Verify Verify
func (c *AWSL) Verify(conn net.Conn) error {
	ws, ok := conn.(awslConn)
	if !ok {
		return errors.New("wrong type")
	}
	auth := model.AddrWithAuth{ANetAddr: ws.Addr, Auth: c.Auth}
	addrBytes, err := json.Marshal(auth)
	if err != nil {
		//log.Println("json marshal : " + err.Error())
		return err
	}
	_, err = ws.Write(addrBytes)
	return err
}

type awslConn struct {
	*websocket.Conn
	Addr model.ANetAddr
}
