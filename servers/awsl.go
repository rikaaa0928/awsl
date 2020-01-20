package servers

import (
	"encoding/json"
	"golang.org/x/net/websocket"
	"net"
	"net/http"
)

// NewAWSL NewAWSL
func NewAWSL(listenHost, listenPort, uri, key, cert string) *AWSL {
	return &AWSL{
		IP:   listenHost,
		Port: listenPort,
		URI:  uri,
		Listener: &AWSListener{
			C:    make(chan net.Conn),
			IP:   listenHost,
			Port: listenPort,
		},
		Cert: cert,
		Key:  key,
	}
}

// AWSL AWSL
type AWSL struct {
	IP       string
	Port     string
	URI      string
	Listener *AWSListener
	Cert     string
	Key      string
}

func (s *AWSL) awslHandler(conn *websocket.Conn) {
	ac := &awslConn{
		Conn: conn,
		C:    make(chan int),
	}
	s.Listener.C <- ac
	<-ac.C
}

// Listen server
func (s *AWSL) Listen() net.Listener {
	http.Handle("/"+s.URI, websocket.Handler(s.awslHandler))
	go func() {
		err := http.ListenAndServeTLS(s.IP+":"+s.Port, s.Cert, s.Key, nil)
		if err != nil {
			panic("ListenAndServe: " + err.Error())
		}
	}()
	return s.Listener
}

// ReadRemote server
func (s *AWSL) ReadRemote(c net.Conn) (ANetAddr, error) {
	jsonBytes := make([]byte, 1024)
	n, err := c.Read(jsonBytes)
	if err != nil {
		return ANetAddr{}, err
	}
	addr := ANetAddr{}
	err = json.Unmarshal(jsonBytes[:n], &addr)
	return addr, err
}

// AWSListener listener
type AWSListener struct {
	C    chan net.Conn
	IP   string
	Port string
}

func (l *AWSListener) Accept() (net.Conn, error) {
	c := <-l.C
	return c, nil
}

func (l *AWSListener) Close() error {
	return nil
}

func (l AWSListener) Addr() net.Addr {
	return &net.IPAddr{
		IP:   net.ParseIP(l.IP),
		Zone: "",
	}
}

type awslConn struct {
	*websocket.Conn
	C chan int
}

func (c *awslConn) Close() error {
	err := c.Conn.Close()
	c.C <- 1
	return err
}
