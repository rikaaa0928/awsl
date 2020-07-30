package servers

import (
	"encoding/json"
	"errors"
	"net"

	"github.com/rikaaa0928/awsl/model"
	"github.com/rikaaa0928/awsl/tools"
	"github.com/rikaaa0928/awsl/tools/crypt"
)

// NewTCP new tcp
func NewTCP(listenHost, listenPort, auth, tag string, id int) *TCP {
	return &TCP{IP: listenHost,
		Port: listenPort,
		Auth: auth,
		tag:  tag,
		id:   id}
}

// TCP TCP
type TCP struct {
	IP   string
	Port string
	Auth string
	tag  string
	id   int
}

// Listen Listen
func (s *TCP) Listen() net.Listener {
	l, err := net.Listen("tcp", net.JoinHostPort(s.IP, s.Port))
	if err != nil {
		panic(err)
	}
	return cryptListenner{l}
}

// ReadRemote ReadRemote
func (s *TCP) ReadRemote(c net.Conn) (model.ANetAddr, error) {
	buf := tools.MemPool.Get(65536)
	defer tools.MemPool.Put(buf)
	n, jsonBytes, err := tools.Receive(c, buf)
	if err != nil {
		return model.ANetAddr{}, err
	}
	a := model.AddrWithAuth{}
	err = json.Unmarshal(jsonBytes[:n], &a)
	if err != nil {
		return model.ANetAddr{}, err
	}
	if a.Auth != s.Auth {
		return model.ANetAddr{}, errors.New("Authentication failed : " + string(jsonBytes[:n]))
	}
	return a.ANetAddr, nil
}

// IDTag id and tag
func (s *TCP) IDTag() (int, string) {
	return s.id, s.tag
}

type cryptListenner struct {
	net.Listener
}

func (l cryptListenner) Accept() (net.Conn, error) {
	conn, err := l.Listener.Accept()
	if err != nil {
		return conn, err
	}
	return cryptConn{Conn: conn, cry: crypt.Simple(1)}, nil
}

type cryptConn struct {
	net.Conn
	cry crypt.Cryptor
}

func (c cryptConn) Read(b []byte) (n int, err error) {
	n, err = c.Conn.Read(b)
	c.cry.Decrypt(b, n)
	return
}

func (c cryptConn) Write(b []byte) (n int, err error) {
	c.cry.Encrypt(b, len(b))
	n, err = c.Conn.Write(b)
	return
}
