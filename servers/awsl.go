package servers

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"

	"github.com/Evi1/awsl/model"
	"github.com/Evi1/awsl/tools"
	"golang.org/x/net/websocket"
)

// NewAWSL NewAWSL
func NewAWSL(ctx context.Context, conf model.In, id int) *AWSL {
	a := &AWSL{
		IP:    conf.Host,
		Port:  conf.Port,
		URI:   conf.Awsl.URI,
		Auth:  conf.Awsl.Auth,
		Conns: make(chan net.Conn, conf.Awsl.Chan),
		Cert:  conf.Awsl.Cert,
		Key:   conf.Awsl.Key,
		id:    id,
		tag:   conf.Tag,
		//CloseChan: make(chan int8),
		closeWait: tools.NewCloseWait(ctx),
	}
	return a
}

/*func NewAWSL(ctx context.Context, listenHost, listenPort, uri, auth, key, cert string, connsSize int) *AWSL {
	a := &AWSL{
		IP:    listenHost,
		Port:  listenPort,
		URI:   uri,
		Auth:  auth,
		Conns: make(chan net.Conn, connsSize),
		Cert:  cert,
		Key:   key,
		//CloseChan: make(chan int8),
		closeWait: tools.NewCloseWait(ctx),
	}
	return a
}*/

// AWSL AWSL
type AWSL struct {
	IP   string
	Port string
	URI  string
	Auth string
	// Listener *AWSListener
	Cert  string
	Key   string
	Conns chan net.Conn
	Srv   http.Server
	id    int
	tag   string
	//CloseChan chan int8
	closeWait *tools.CloseWait
}

func (s *AWSL) awslHandler(conn *websocket.Conn) {
	ac := &awslConn{
		Conn: conn,
		//CloseChan: make(chan int8),
		closeWait: tools.NewCloseWait(s.closeWait.Ctx),
	}

	s.Conns <- ac
	ac.closeWait.WaitClose()
}

// Listen server
func (s *AWSL) Listen() net.Listener {
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

// IDTag id and tag
func (s *AWSL) IDTag() (int, string) {
	return s.id, s.tag
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
	//case <-s.CloseChan:
	case <-s.closeWait.WaitClose():
		return nil, errors.New("awsl listenner closed")
	}
}

// Close Close
func (s *AWSL) Close() error {
	/*defer func() {
		recover()
	}()*/
	s.closeWait.Close()
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
	//CloseChan chan int8
	closeWait *tools.CloseWait
}

func (c *awslConn) Close() error {
	/*defer func() {
		recover()
	}()*/
	err := c.Conn.Close()
	c.closeWait.Close()
	return err
}
