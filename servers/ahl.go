package servers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Evi1/awsl/model"
	"github.com/Evi1/awsl/tools"
)

// AHL AHL
type AHL struct {
	IP        string
	Port      string
	URI       string
	Auth      string
	Cert      string
	Key       string
	connPool  map[uint64]*ahlConn
	id        uint64
	Conns     chan net.Conn
	CloseChan chan int8
}

func (s *AHL) serve(w http.ResponseWriter, r *http.Request) {
	if strings.ToUpper(r.Method) == http.MethodGet {
		r.ParseForm()
		switch r.Form.Get("action") {
		case "connect":
			pw, err := r.Cookie("pw")
			if err != nil {
				w.WriteHeader(http.StatusNotAcceptable)
				return
			}
			if pw.Value != s.Auth {
				w.WriteHeader(http.StatusNonAuthoritativeInfo)
				return
			}
			addrStr := r.Form.Get("addr")
			addr := model.ANetAddr{}
			json.Unmarshal([]byte(addrStr), &addr)
			_, ok := s.connPool[s.id]
			for ok {
				s.id++
				_, ok = s.connPool[s.id]
			}
			s.connPool[s.id] = newAHlConn(addr)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(strconv.FormatUint(s.id, 10)))
			s.Conns <- s.connPool[s.id]
		case "close":
			pw, err := r.Cookie("pw")
			if err != nil {
				w.WriteHeader(http.StatusNotAcceptable)
				return
			}
			if pw.Value != s.Auth {
				w.WriteHeader(http.StatusNonAuthoritativeInfo)
				return
			}
			id, err := r.Cookie("id")
			if err != nil {
				w.WriteHeader(http.StatusNotAcceptable)
				return
			}
			uid, err := strconv.ParseUint(id.Value, 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusNotAcceptable)
				return
			}
			conn, ok := s.connPool[uid]
			if !ok {
				w.WriteHeader(http.StatusNotAcceptable)
				return
			}
			conn.Close()
			delete(s.connPool, uid)
		default:
			w.WriteHeader(http.StatusNotAcceptable)
		}
	} else if strings.ToUpper(r.Method) == http.MethodPost {
		pw, err := r.Cookie("pw")
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}
		if pw.Value != s.Auth {
			w.WriteHeader(http.StatusNonAuthoritativeInfo)
			return
		}
		id := string([]rune(r.URL.RequestURI())[len(s.URI)+1:])
		uid, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}
		req, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}
		fmt.Println(string(req))
		conn, ok := s.connPool[uid]
		if !ok {
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}
		buf := tools.MemPool.Get(65536)
		defer tools.MemPool.Put(buf)

		n, err := writeToChan(conn.ReadChan, conn.close, req)
		fmt.Println("write req", n, string(req))
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte(err.Error()))
			conn.Close()
			return
		}

		n, err = readFromChan(conn.WriteChan, conn.close, buf)
		fmt.Println("readFromChan WriteChan", n, string(buf[:n]))
		if n != 0 {
			w.WriteHeader(http.StatusOK)
			w.Write(buf[:n])
		}
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte(err.Error()))
			conn.Close()
		}

	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(r.Method))
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
	select {
	case conn := <-s.Conns:
		return conn, nil
	case <-s.CloseChan:
	}
	return nil, errors.New("closed")
}

// Close Close
func (s *AHL) Close() error {
	s.CloseChan <- 1
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

func newAHlConn(addr model.ANetAddr) *ahlConn {
	return &ahlConn{ReadChan: make(chan connData, 1),
		WriteChan:  make(chan connData, 1),
		close:      make(chan int8),
		remoteAddr: addr}
}

type ahlConn struct {
	ReadChan   chan connData
	WriteChan  chan connData
	close      chan int8
	remoteAddr model.ANetAddr
}

func (c *ahlConn) Read(b []byte) (n int, err error) {
	return readFromChan(c.ReadChan, c.close, b)
	/*select {
	case data := <-c.ReadChan:
		copy(b, data.data)
		return data.n, data.err
	case <-time.After(time.Minute):
		return 0, io.EOF
	case <-c.close:
		return 0, io.EOF
	}*/
}
func (c *ahlConn) Write(b []byte) (n int, err error) {
	return writeToChan(c.WriteChan, c.close, b)
	/*data := connData{data: make([]byte, len(b)), n: len(b)}
	copy(data.data, b)
	select {
	case c.WriteChan <- data:
		return len(b), nil
	case <-time.After(time.Minute):
		return 0, io.EOF
	case <-c.close:
		return 0, io.EOF
	}*/
}

func readFromChan(c chan connData, close chan int8, b []byte) (n int, err error) {
	select {
	case data := <-c:
		copy(b, data.data)
		return data.n, data.err
	case <-time.After(time.Minute):
		return 0, io.EOF
	case <-close:
		return 0, io.EOF
	}
}

func writeToChan(c chan connData, close chan int8, b []byte) (n int, err error) {
	data := connData{data: make([]byte, len(b)), n: len(b)}
	copy(data.data, b)
	select {
	case c <- data:
		return len(b), nil
	case <-time.After(time.Minute):
		return 0, io.EOF
	case <-close:
		return 0, io.EOF
	}
}

func (c *ahlConn) Close() error {
	defer func() {
		recover()
	}()
	close(c.close)
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
