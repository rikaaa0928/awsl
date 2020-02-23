package servers

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/Evi1/awsl/model"
)

// NewH2C NewH2C
func NewH2C(listenHost, listenPort, uri, auth, key, cert string, connsSize int) *H2C {
	return &H2C{
		IP:        listenHost,
		Port:      listenPort,
		URI:       uri + "/",
		Auth:      auth,
		Key:       key,
		Cert:      cert,
		CloseChan: make(chan int8),
		Conns:     make(chan net.Conn, connsSize)}
}

// H2C H2C
type H2C struct {
	IP        string
	Port      string
	URI       string
	Auth      string
	Cert      string
	Key       string
	CloseChan chan int8
	Conns     chan net.Conn
}

func (s *H2C) serve(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "PUT required.", 400)
		return
	}
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
	c := &h2cConn{w: w, r: r, closeChan: make(chan int8)}

	s.Conns <- c

	<-c.closeChan
}

// Listen Listen
func (s *H2C) Listen() net.Listener {
	mux := http.NewServeMux()
	mux.HandleFunc("/"+s.URI, s.serve)
	/*var srv http.Server
	srv.Handler = mux
	srv.Addr = s.IP + ":" + s.Port
	srv.ConnState = idleTimeoutHook()
	http2.ConfigureServer(&srv, &http2.Server{})*/
	//http.HandleFunc("/"+s.URI, s.serve)
	go func() {
		if len(s.Cert) == 0 || len(s.Key) == 0 {
			//err := srv.ListenAndServe()
			err := http.ListenAndServe(s.IP+":"+s.Port, mux)
			if err != nil {
				panic("ListenAndServe: " + err.Error())
			}
		} else {
			//err := srv.ListenAndServeTLS(s.Cert, s.Key)
			err := http.ListenAndServeTLS(s.IP+":"+s.Port, s.Cert, s.Key, mux)
			if err != nil {
				panic("ListenAndServe: " + err.Error())
			}
		}
	}()
	return s
}

// ReadRemote ReadRemote
func (s *H2C) ReadRemote(c net.Conn) (model.ANetAddr, error) {
	if ac, ok := c.(*h2cConn); ok {
		addr := model.ANetAddr{}
		addrCookie, err := ac.r.Cookie("addr")
		if err != nil {
			return addr, err
		}
		addrStr, err := url.QueryUnescape(addrCookie.Value)
		if err != nil {
			return addr, err
		}
		err = json.Unmarshal([]byte(addrStr), &addr)
		if err != nil {
			return addr, err
		}
		pw, err := ac.r.Cookie("pw")
		if err != nil {
			return addr, err
		}
		if pw.Value != s.Auth {
			return addr, errors.New("wtf")
		}
		return addr, nil
	}
	return model.ANetAddr{}, errors.New("wrong conn type")
}

// Accept Accept
func (s *H2C) Accept() (net.Conn, error) {
	select {
	case conn := <-s.Conns:
		return conn, nil
	case <-s.CloseChan:
	}
	return nil, errors.New("closed")
}

// Close Close
func (s *H2C) Close() error {
	defer func() {
		recover()
	}()
	close(s.CloseChan)
	return nil
}

// Addr Addr
func (s *H2C) Addr() net.Addr {
	return &net.IPAddr{
		IP:   net.ParseIP(s.IP),
		Zone: "",
	}
}

type h2cConn struct {
	w         http.ResponseWriter
	r         *http.Request
	closeChan chan int8
}

func (c *h2cConn) Read(b []byte) (n int, err error) {
	return c.r.Body.Read(b)
}

func (c *h2cConn) Write(b []byte) (n int, err error) {
	n, err = c.w.Write(b)
	c.w.(http.Flusher).Flush()
	return
}

func (c *h2cConn) Close() error {
	defer func() {
		recover()
	}()
	close(c.closeChan)
	return c.r.Body.Close()
}
func (c *h2cConn) LocalAddr() net.Addr {
	return nil
}
func (c *h2cConn) RemoteAddr() net.Addr {
	return nil
}
func (c *h2cConn) SetDeadline(t time.Time) error {
	return nil
}
func (c *h2cConn) SetReadDeadline(t time.Time) error {
	return nil
}
func (c *h2cConn) SetWriteDeadline(t time.Time) error {
	return nil
}
