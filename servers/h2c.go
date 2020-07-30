package servers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/rikaaa0928/awsl/model"
	"github.com/rikaaa0928/awsl/tools"
)

// NewH2C NewH2C
func NewH2C(ctx context.Context, conf model.In, id int) *H2C {
	return &H2C{
		IP:        conf.Host,
		Port:      conf.Port,
		URI:       conf.Awsl.URI,
		Auth:      conf.Awsl.Auth,
		Cert:      conf.Awsl.Cert,
		Key:       conf.Awsl.Key,
		closeWait: tools.NewCloseWait(ctx),
		id:        id,
		tag:       conf.Tag,
		//CloseChan: make(chan int8),
		Conns: make(chan net.Conn, conf.Awsl.Chan)}
}

// H2C H2C
type H2C struct {
	IP   string
	Port string
	URI  string
	Auth string
	Cert string
	Key  string
	id   int
	tag  string
	//CloseChan chan int8
	closeWait *tools.CloseWait
	Conns     chan net.Conn
	Srv       http.Server
}

type rewrite struct {
	w http.ResponseWriter
}

func (w rewrite) Write(b []byte) (n int, err error) {
	n, err = w.w.Write(b)
	w.w.(http.Flusher).Flush()
	return
}

func (s *H2C) serve(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "PUT required.", http.StatusBadRequest)
		return
	}

	addr := model.ANetAddr{}
	addrCookie, err := r.Cookie("addr")
	if err != nil {
		http.Error(w, "address required.", http.StatusBadRequest)
		return
	}
	addrStr, err := url.QueryUnescape(addrCookie.Value)
	if err != nil {
		http.Error(w, "address required.", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal([]byte(addrStr), &addr)
	if err != nil {
		http.Error(w, "address required.", http.StatusBadRequest)
		return
	}
	pw, err := r.Cookie("pw")
	if err != nil {
		http.Error(w, "wtf?", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "flush required.", http.StatusHTTPVersionNotSupported)
		return
	}
	f.Flush()
	wr, ww := io.Pipe()
	rr, rw := io.Pipe()
	defer func() {
		wr.Close()
		ww.Close()
		rr.Close()
		rw.Close()
	}()
	c := &h2cConn{w: ww, r: rr, Pw: pw.Value, Addr: addr}

	s.Conns <- c

	go io.Copy(rewrite{w}, wr)
	io.Copy(rw, r.Body)
}

// Listen Listen
func (s *H2C) Listen() net.Listener {
	mux := http.NewServeMux()
	mux.HandleFunc("/"+s.URI+"/", s.serve)
	s.Srv = http.Server{Addr: s.IP + ":" + s.Port, Handler: mux}
	/*var srv http.Server
	srv.Handler = mux
	srv.Addr = s.IP + ":" + s.Port
	srv.ConnState = idleTimeoutHook()
	http2.ConfigureServer(&srv, &http2.Server{})*/
	//http.HandleFunc("/"+s.URI, s.serve)
	go func() {
		if len(s.Cert) == 0 || len(s.Key) == 0 {
			err := s.Srv.ListenAndServe()
			//err := http.ListenAndServe(s.IP+":"+s.Port, mux)
			if err != nil {
				panic("ListenAndServe: " + err.Error())
			}
		} else {
			err := s.Srv.ListenAndServeTLS(s.Cert, s.Key)
			//err := http.ListenAndServeTLS(s.IP+":"+s.Port, s.Cert, s.Key, mux)
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
		if ac.Pw != s.Auth {
			return model.ANetAddr{}, errors.New("wtf")
		}
		return ac.Addr, nil
	}
	return model.ANetAddr{}, errors.New("wrong conn type")
}

// Accept Accept
func (s *H2C) Accept() (net.Conn, error) {
	select {
	case conn := <-s.Conns:
		return conn, nil
	case <-s.closeWait.WaitClose():
	}
	return nil, errors.New("h2c server closed")
}

// IDTag id and tag
func (s *H2C) IDTag() (int, string) {
	return s.id, s.tag
}

// Close Close
func (s *H2C) Close() error {
	/*defer func() {
		recover()
	}()
	close(s.CloseChan)*/
	s.closeWait.Close()
	return s.Srv.Shutdown(context.Background())
}

// Addr Addr
func (s *H2C) Addr() net.Addr {
	return &net.IPAddr{
		IP:   net.ParseIP(s.IP),
		Zone: "",
	}
}

type h2cConn struct {
	w    io.WriteCloser
	r    io.ReadCloser
	Pw   string
	Addr model.ANetAddr
}

func (c *h2cConn) Read(b []byte) (n int, err error) {
	return c.r.Read(b)
}

func (c *h2cConn) Write(b []byte) (n int, err error) {
	//n, err = c.w.Write(b)
	//c.w.(http.Flusher).Flush()
	return c.w.Write(b)
}

func (c *h2cConn) Close() error {
	c.w.Close()
	return c.r.Close()
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
