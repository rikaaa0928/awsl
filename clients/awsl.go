package clients

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"log"
	"net"

	"github.com/rikaaa0928/awsl/config"
	"github.com/rikaaa0928/awsl/model"
	"github.com/rikaaa0928/awsl/tools/dialer"
	"golang.org/x/net/websocket"
)

// NewAWSL NewAWSL
func NewAWSL(id int, conf model.Out) *AWSL {
	wsConfig, err := websocket.NewConfig("wss://"+conf.Awsl.Host+":"+conf.Awsl.Port+"/"+conf.Awsl.URI, "https://"+conf.Awsl.Host+":"+conf.Awsl.Port+"/")
	if err != nil {
		panic(err)
	}
	m := make(map[string][]string)
	hp := net.JoinHostPort(conf.Awsl.Host, conf.Awsl.Port)
	m[hp] = []string{hp}
	backup := conf.Awsl.BackUp
	if backup != nil {
		m[hp] = append(m[hp], backup...)
	}
	d := &dialer.MultiAddr{Hosts: m, HostInUse: make(map[string]uint)}
	return &AWSL{
		ServerHost: conf.Awsl.Host,
		ServerPort: conf.Awsl.Port,
		URI:        conf.Awsl.URI,
		Auth:       conf.Awsl.Auth,
		wsConfig:   wsConfig,
		mDialer:    d,
		id:         id,
		tag:        conf.Tag,
	}
}

// AWSL AWSL
type AWSL struct {
	ServerHost string
	ServerPort string
	URI        string
	Auth       string
	id         int
	tag        string
	wsConfig   *websocket.Config
	mDialer    *dialer.MultiAddr
}

// Dial Dial
func (c *AWSL) Dial(addr model.ANetAddr) (net.Conn, error) {
	bc, err := c.mDialer.Dial("tcp", net.JoinHostPort(c.ServerHost, c.ServerPort))
	if err != nil {
		log.Println("awsl client dial", err)
		return nil, err
	}
	tc := tls.Client(bc, &tls.Config{
		InsecureSkipVerify: config.GetConf().NoVerify,
		ServerName:         c.ServerHost,
	})
	ws, err := websocket.NewClient(c.wsConfig, tc)
	if err != nil {
		log.Println("awsl client new client", err)
		return nil, err
	}
	/*ws, err := websocket.DialConfig(wsConfig)
	// ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		//log.Println("dial:" + err.Error())
		return ws, err
	}*/
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

// IDTag id tag
func (c *AWSL) IDTag() (int, string) {
	return c.id, c.tag
}

type awslConn struct {
	*websocket.Conn
	Addr model.ANetAddr
}
