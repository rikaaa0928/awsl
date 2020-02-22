package servers

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Evi1/awsl/model"
	"github.com/Evi1/awsl/tools"
)

// NewAHL NewAHL
func NewAHL(listenHost, listenPort, uri, auth, key, cert string, connsSize int) *AHL {
	return &AHL{
		IP:        listenHost,
		Port:      listenPort,
		URI:       uri + "/",
		Auth:      auth,
		Key:       key,
		Cert:      cert,
		connPool:  make(map[uint64]*ahlConn),
		Conns:     make(chan net.Conn, connsSize),
		CloseChan: make(chan int8),
		poolLock:  sync.Mutex{}}
}

// AHL AHL
type AHL struct {
	IP        string
	Port      string
	URI       string
	Auth      string
	Cert      string
	Key       string
	connPool  map[uint64]*ahlConn
	poolLock  sync.Mutex
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
				log.Println(pw.Value, s.Auth)
				w.WriteHeader(http.StatusNonAuthoritativeInfo)
				return
			}
			addrStr := r.Form.Get("addr")
			addr := model.ANetAddr{}
			json.Unmarshal([]byte(addrStr), &addr)
			s.poolLock.Lock()
			_, ok := s.connPool[s.id]
			for ok {
				s.id++
				_, ok = s.connPool[s.id]
			}
			s.connPool[s.id] = newAHlConn(addr, s, s.id)
			s.poolLock.Unlock()
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
			id := string([]rune(r.URL.RequestURI())[len(s.URI)+1:])
			uid, err := strconv.ParseUint(id, 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusNotAcceptable)
				return
			}
			s.poolLock.Lock()
			conn, ok := s.connPool[uid]
			s.poolLock.Unlock()
			if !ok {
				w.WriteHeader(http.StatusNotAcceptable)
				return
			}
			conn.Close()
			//delete(s.connPool, uid)
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
		conn, ok := s.connPool[uid]
		if !ok || conn.closed {
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}
		if len(req) > 0 {
			data := connData{data: make([]byte, len(req)), n: len(req)}
			copy(data.data, req)
			conn.ReadLock.Lock()
			conn.ReadBuf = append(conn.ReadBuf, data)
			conn.ReadLock.Unlock()
		}
		// time.Sleep(10 * time.Millisecond)
		timeOut := time.Now().Add(3 * time.Second)
		for time.Now().Before(timeOut) {
			if len(conn.WriteBuf) == 0 {
				time.Sleep(time.Second)
				continue
			}
			conn.WriteLock.Lock()
			if len(conn.WriteBuf) == 0 {
				conn.WriteLock.Unlock()
				continue
			}
			w.Header().Set("Num", strconv.Itoa(len(conn.WriteBuf)-1))
			w.Header().Set("Content-Type", "application/octet-stream")
			w.WriteHeader(http.StatusOK)
			buf := conn.WriteBuf[0]
			_, err = w.Write(buf.data)
			if err == nil {
				conn.WriteBuf = conn.WriteBuf[1:]
				log.Println(conn.id, len(conn.WriteBuf))
			}
			conn.WriteLock.Unlock()
			return
		}
		w.WriteHeader(http.StatusAccepted)

		/*conn.WriteLock.Lock()
		bufLen := len(conn.WriteBuf)
		if bufLen == 0 {
			conn.WriteLock.Unlock()
			time.Sleep(time.Second)
			conn.WriteLock.Lock()
			bufLen = len(conn.WriteBuf)
		}
		defer conn.WriteLock.Unlock()
		if bufLen == 0 {
			w.WriteHeader(http.StatusAccepted)
			return
		}
		w.Header().Set("Num", strconv.Itoa(bufLen-1))
		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		buf := conn.WriteBuf[0]
		w.Write(buf.data)
		conn.WriteBuf = conn.WriteBuf[1:]*/
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
	return s
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

func newAHlConn(addr model.ANetAddr, l *AHL, id uint64) *ahlConn {
	return &ahlConn{
		ReadBuf:    make([]connData, 0),
		WriteBuf:   make([]connData, 0),
		remoteAddr: addr,
		listenner:  l,
		id:         id}
}

type ahlConn struct {
	id         uint64
	ReadBuf    []connData
	WriteBuf   []connData
	ReadLock   sync.Mutex
	WriteLock  sync.Mutex
	closed     bool
	remoteAddr model.ANetAddr
	listenner  *AHL
}

func (c *ahlConn) Read(b []byte) (n int, err error) {
	return readTimeout(b, &c.ReadBuf, &c.closed, &c.ReadLock, time.Minute)
}

func readTimeout(b []byte, readBuf *[]connData, closed *bool, lock *sync.Mutex, t time.Duration) (n int, err error) {
	timeOut := time.Now().Add(t)
	for time.Now().Before(timeOut) {
		if len(*readBuf) == 0 && !*closed {
			time.Sleep(t / 5)
			continue
		}
		if *closed {
			return 0, io.EOF
		}
		lock.Lock()
		if len(*readBuf) == 0 {
			lock.Unlock()
			continue
		}
		data := (*readBuf)[0]
		if len(b) < data.n {
			copy(b, data.data[:len(b)])
			data.data = data.data[len(b):]
			data.n -= len(b)
			(*readBuf)[0] = data
			lock.Unlock()
			return len(b), nil
		}
		copy(b, data.data)
		*readBuf = (*readBuf)[1:]
		lock.Unlock()
		return data.n, data.err
	}
	return 0, tools.ErrTimeout
}

func (c *ahlConn) Write(b []byte) (n int, err error) {
	if c.closed {
		return 0, io.ErrUnexpectedEOF
	}
	data := connData{data: make([]byte, len(b)), n: len(b)}
	copy(data.data, b)
	c.WriteLock.Lock()
	defer c.WriteLock.Unlock()
	c.WriteBuf = append(c.WriteBuf, data)
	log.Println(c.id, len(c.WriteBuf))
	return len(b), nil
}

func (c *ahlConn) Close() error {
	c.closed = true
	if c.listenner != nil {
		c.listenner.poolLock.Lock()
		defer c.listenner.poolLock.Unlock()
		delete(c.listenner.connPool, c.id)
	}
	c.ReadBuf = make([]connData, 0)
	c.WriteBuf = make([]connData, 0)
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
