package servers

import (
	"golang.org/x/net/websocket"
	"log"
	"net"
	"net/http"
	"time"
)

// AWSL AWSL
type AWSL struct {
	IP       string
	Port     string
	URI      string
	Listener *AWSListener
	Cert     string
	Key      string
}

func (s *AWSL) awslHandeler(conn *websocket.Conn) {
	log.Println(conn.LocalAddr(), conn.RemoteAddr())
	s.Listener.C <- conn
	log.Println(conn)
	conn.Write([]byte("asfdsaf"))
	for {
		time.Sleep(time.Second)
	}
}

// Listen server
func (s *AWSL) Listen() net.Listener {
	http.Handle("/"+s.URI, websocket.Handler(s.awslHandeler))
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
	return ANetAddr{}, nil
}

// AWSListener listener
type AWSListener struct {
	C    chan *websocket.Conn
	IP   string
	Port string
}

func (l *AWSListener) Accept() (net.Conn, error) {
	c := <-l.C
	log.Println(c)
	log.Println(c.Write([]byte("asfdsaf")))
	log.Println(c.LocalAddr(), c.RemoteAddr())
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
