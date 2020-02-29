package servers

import (
	"context"
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
		IP:        listenHost,
		Port:      listenPort,
		URI:       uri,
		Auth:      auth,
		Conns:     make(chan net.Conn, connsSize),
		Cert:      cert,
		Key:       key,
		ConnNum:   make(chan int, 1),
		CloseChan: make(chan int8),
	}
	a.ConnNum <- 0
	return a
}

// AWSL AWSL
type AWSL struct {
	IP   string
	Port string
	URI  string
	Auth string
	// Listener *AWSListener
	Cert      string
	Key       string
	ConnNum   chan int
	Max       int
	Conns     chan net.Conn
	Srv       http.Server
	CloseChan chan int8
}

func (s *AWSL) awslHandler(conn *websocket.Conn) {
	ac := &awslConn{
		Conn:      conn,
		CloseChan: make(chan int8),
	}
	s.Conns <- ac
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
		log.Println("current conn: " + strconv.Itoa(num))
		if num == 0 {
			log.Println("Connection clear")
		}
		s.ConnNum <- num
	}

}

// Listen server
func (s *AWSL) Listen() net.Listener {
	log.Println(s.IP+":"+s.Port, s.Cert, s.Key)
	mux := http.NewServeMux()
	mux.Handle("/"+s.URI, websocket.Handler(s.awslHandler))
	//http.Handle("/"+s.URI, websocket.Handler(s.awslHandler))
	s.Srv = http.Server{Addr: s.IP + ":" + s.Port, Handler: mux}
	go func() {
		if len(s.Cert) == 0 || len(s.Key) == 0 {
			//err := http.ListenAndServe(s.IP+":"+s.Port, mux)
			err := s.Srv.ListenAndServe()
			if err != nil {
				panic("ListenAndServe: " + err.Error())
			}
		} else {
			//err := http.ListenAndServeTLS(s.IP+":"+s.Port, s.Cert, s.Key, mux)
			err := s.Srv.ListenAndServeTLS(s.Cert, s.Key)
			if err != nil {
				panic("ListenAndServe: " + err.Error())
			}
		}

	}()
	return s
}

// ReadRemote server
func (s *AWSL) ReadRemote(c net.Conn) (model.ANetAddr, error) {
	buf := tools.MemPool.Get(65536)
	defer tools.MemPool.Put(buf)
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
/*type AWSListener struct {
	Conns chan net.Conn
	IP    string
	Port  string
}*/

// Accept Accept
func (s *AWSL) Accept() (net.Conn, error) {
	select {
	case c := <-s.Conns:
		return c, nil
	case <-s.CloseChan:
		return nil, errors.New("listenner closed")
	}
}

// Close Close
func (s *AWSL) Close() error {
	defer func() {
		recover()
	}()
	close(s.CloseChan)
	return s.Srv.Shutdown(context.Background())
}

// Addr Addr
func (s *AWSL) Addr() net.Addr {
	return &net.IPAddr{
		IP:   net.ParseIP(s.IP),
		Zone: "",
	}
}

type awslConn struct {
	net.Conn
	CloseChan chan int8
}

func (c *awslConn) Close() error {
	defer func() {
		recover()
	}()
	err := c.Conn.Close()
	close(c.CloseChan)
	return err
}
