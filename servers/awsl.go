package servers

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/Evi1/awsl/config"
	"github.com/Evi1/awsl/model"
	"github.com/Evi1/awsl/tools"
	"golang.org/x/net/websocket"
)

// NewAWSL NewAWSL
func NewAWSL(listenHost, listenPort, uri, auth, key, cert string, connsSize int) *AWSL {
	a := &AWSL{
		IP:   listenHost,
		Port: listenPort,
		URI:  uri,
		Auth: auth,
		Listener: &AWSListener{
			Conns: make(chan net.Conn, connsSize),
			IP:    listenHost,
			Port:  listenPort,
		},
		Cert:    cert,
		Key:     key,
		ConnNum: make(chan int, 1),
	}
	a.ConnNum <- 0
	return a
}

// AWSL AWSL
type AWSL struct {
	IP       string
	Port     string
	URI      string
	Auth     string
	Listener *AWSListener
	Cert     string
	Key      string
	ConnNum  chan int
	Max      int
}

func (s *AWSL) awslHandler(conn *websocket.Conn) {
	ac := &awslConn{
		Conn: conn,
		CloseChan:    make(chan int8),
	}
	s.Listener.Conns <- ac
	if config.Debug {
		num := <-s.ConnNum
		num++
		if num > s.Max {
			s.Max = num
			log.Println("max conn: " + strconv.Itoa(num))
		}
		s.ConnNum <- num
	}

	<-ac.CloseChan

	if config.Debug {
		num := <-s.ConnNum
		num--
		s.ConnNum <- num
	}

}

// Listen server
func (s *AWSL) Listen() net.Listener {
	log.Println(s.IP+":"+s.Port, s.Cert, s.Key)
	http.Handle("/"+s.URI, websocket.Handler(s.awslHandler))
	go func() {
		if len(s.Cert) == 0 || len(s.Key) == 0 {
			err := http.ListenAndServe(s.IP+":"+s.Port, nil)
			if err != nil {
				panic("ListenAndServe: " + err.Error())
			}
		} else {
			err := http.ListenAndServeTLS(s.IP+":"+s.Port, s.Cert, s.Key, nil)
			if err != nil {
				panic("ListenAndServe: " + err.Error())
			}
		}

	}()
	return s.Listener
}

// ReadRemote server
func (s *AWSL) ReadRemote(c net.Conn) (model.ANetAddr, error) {
	buf := tools.MemPool.Get(65536)
	defer func() {
		tools.MemPool.Put(buf)
	}()
	n, jsonBytes, err := tools.Receive(c, buf)
	//jsonBytes := make([]byte, 1024)
	//n, err := c.Read(jsonBytes)
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

// AWSListener listener
type AWSListener struct {
	Conns chan net.Conn
	IP    string
	Port  string
}

// Accept Accept
func (l *AWSListener) Accept() (net.Conn, error) {
	c := <-l.Conns
	return c, nil
}

// Close Close
func (l *AWSListener) Close() error {
	return nil
}

// Addr Addr
func (l AWSListener) Addr() net.Addr {
	return &net.IPAddr{
		IP:   net.ParseIP(l.IP),
		Zone: "",
	}
}

type awslConn struct {
	*websocket.Conn
	CloseChan chan int8
}

func (c *awslConn) Close() error {
	err := c.Conn.Close()
	c.CloseChan <- 1
	return err
}
