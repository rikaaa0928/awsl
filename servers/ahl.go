package servers

import (
	"errors"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/Evi1/awsl/model"
)

// AHL AHL
type AHL struct {
	IP       string
	Port     string
	URI      string
	Auth     string
	Cert     string
	Key      string
	ConnNum  chan int
	connPool map[uint64]*ahlConn
}

func (s *AHL) serve(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

	} else if r.Method == http.MethodPost {

	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("."))
	}
}

// Listen Listen
func (s *AHL) Listen() net.Listener {
	mux := http.NewServeMux()
	mux.HandleFunc("/"+s.URI, s.serve)
	//http.HandleFunc("/"+s.URI, s.serve)
	go func() {
		if len(s.Cert) == 0 || len(s.Key) == 0 {
			err := http.ListenAndServe(s.IP+":"+s.Port, mux)
			if err != nil {
				panic("ListenAndServe: " + err.Error())
			}
		} else {
			err := http.ListenAndServeTLS(s.IP+":"+s.Port, s.Cert, s.Key, mux)
			if err != nil {
				panic("ListenAndServe: " + err.Error())
			}
		}

	}()
	return nil
}

// ReadRemote ReadRemote
func (s *AHL) ReadRemote(c net.Conn) (model.ANetAddr, error) {
	if ac, ok := c.(*ahlConn); ok {
		return ac.remoteAddr, nil
	}
	return model.ANetAddr{}, errors.New("wrong conn type")
}

// Accept Accept
func (s *AHL) Accept() (net.Conn, error) {
	return nil, nil
}

// Close Close
func (s *AHL) Close() error {
	return nil
}

// Addr Addr
func (s *AHL) Addr() net.Addr {
	return &net.IPAddr{
		IP:   net.ParseIP(s.IP),
		Zone: "",
	}
}

type connData struct {
	data []byte
	n    int
	err  error
}

type ahlConn struct {
	readChan   chan connData
	writeChan  chan connData
	close      chan int8
	remoteAddr model.ANetAddr
}

func (c *ahlConn) Read(b []byte) (n int, err error) {
	select {
	case data := <-c.readChan:
		copy(b, data.data)
		return data.n, data.err
	case <-time.After(time.Minute):
		return 0, io.EOF
	case <-c.close:
		return 0, io.EOF
	}
}
func (c *ahlConn) Write(b []byte) (n int, err error) {
	data := connData{data: make([]byte, 65536), n: len(b)}
	copy(data.data, b)
	select {
	case c.writeChan <- data:
		return len(b), nil
	case <-time.After(time.Minute):
		return 0, io.EOF
	case <-c.close:
		return 0, io.EOF
	}
}
func (c *ahlConn) Close() error {
	c.close <- 1
	return nil
}
func (c *ahlConn) LocalAddr() net.Addr {
	return nil
}
func (c *ahlConn) RemoteAddr() net.Addr {
	return nil
}
func (c *ahlConn) SetDeadline(t time.Time) error {
	return nil
}
func (c *ahlConn) SetReadDeadline(t time.Time) error {
	return nil
}
func (c *ahlConn) SetWriteDeadline(t time.Time) error {
	return nil
}
