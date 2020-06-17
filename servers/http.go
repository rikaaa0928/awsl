package servers

import (
	"errors"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/Evi1/awsl/model"
)

// NewHTTP NewHTTP
func NewHTTP(listenHost, listenPort string, connsSize int) *HTTPServer {
	a := &HTTPServer{
		IP:    listenHost,
		Port:  listenPort,
		Conns: make(chan net.Conn, connsSize),
	}
	return a
}

// HTTPServer HTTPServer
type HTTPServer struct {
	IP    string
	Port  string
	Conns chan net.Conn
	Max   int
}

// Listen server
func (s *HTTPServer) Listen() net.Listener {
	server := &http.Server{
		Addr: s.IP + ":" + s.Port,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodConnect {
				w.WriteHeader(http.StatusOK)

				// host port
				sl := strings.Split(r.Host, ":")
				if len(sl) != 2 {
					log.Println("addr error : " + r.Host)
					http.Error(w, "addr error : "+r.Host, http.StatusBadRequest)
					return
				}
				port, err := strconv.Atoi(sl[1])
				if err != nil {
					log.Println(err.Error())
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				addr := model.ANetAddr{Host: sl[0], Port: port}
				if net.ParseIP(addr.Host) == nil {
					addr.Typ = model.RAWADDR
				} else if net.ParseIP(addr.Host).To4() != nil {
					addr.Typ = model.IPV4ADDR
				} else if net.ParseIP(addr.Host).To16() != nil {
					addr.Typ = model.IPV6ADDR
				}

				hijacker, ok := w.(http.Hijacker)
				if !ok {
					log.Println("Hijacking not supported")
					http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
					return
				}
				clientConn, _, err := hijacker.Hijack()
				if err != nil {
					log.Println(err.Error())
					http.Error(w, err.Error(), http.StatusServiceUnavailable)
					return
				}

				con := &httpConn{Conn: clientConn, CloseChan: make(chan int8), addr: addr}
				s.Conns <- con
				<-con.CloseChan
			} else {
				rHost := ""
				rPort := 80
				var err error
				// host port
				if strings.Contains(r.Host, ":") {
					sl := strings.Split(r.Host, ":")
					if len(sl) != 2 {
						log.Println("addr error : " + r.Host)
						http.Error(w, "addr error : "+r.Host, http.StatusBadRequest)
						return
					}
					rHost = sl[0]
					rPort, err = strconv.Atoi(sl[1])
					if err != nil {
						log.Println(err.Error())
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}
				} else {
					rHost = r.Host
				}

				addr := model.ANetAddr{Host: rHost, Port: rPort}

				if net.ParseIP(addr.Host) == nil {
					addr.Typ = model.RAWADDR
				} else if net.ParseIP(addr.Host).To4() != nil {
					addr.Typ = model.IPV4ADDR
				} else if net.ParseIP(addr.Host).To16() != nil {
					addr.Typ = model.IPV6ADDR
				}

				conn := &HTTPGetConn{W: w, R: r, addr: addr, CloseChan: make(chan int8)}
				s.Conns <- conn
				<-conn.CloseChan
			}
		}),
	}
	go server.ListenAndServe()
	return s
}

// ReadRemote server
func (s *HTTPServer) ReadRemote(c net.Conn) (model.ANetAddr, error) {
	conn, ok := c.(model.AWSLConn)
	if !ok {
		return model.ANetAddr{}, errors.New("conn not httpConn")
	}
	return conn.GetAddr(), nil
}

// Accept Accept
func (s *HTTPServer) Accept() (net.Conn, error) {
	c := <-s.Conns
	return c, nil
}

// Close Close
func (s *HTTPServer) Close() error {
	return nil
}

// Addr Addr
func (s *HTTPServer) Addr() net.Addr {
	return &net.IPAddr{
		IP:   net.ParseIP(s.IP),
		Zone: "",
	}
}

type httpConn struct {
	net.Conn
	addr      model.ANetAddr
	CloseChan chan int8
}

func (c *httpConn) Close() error {
	err := c.Conn.Close()
	c.CloseChan <- 1
	return err
}

func (c *httpConn) GetAddr() model.ANetAddr {
	return c.addr
}

// HTTPGetConn HTTPGetConn
type HTTPGetConn struct {
	W    http.ResponseWriter
	R    *http.Request
	addr model.ANetAddr
	net.Conn
	CloseChan chan int8
}

// GetAddr GetAddr
func (c *HTTPGetConn) GetAddr() model.ANetAddr {
	return c.addr
}

// Close Close
func (c *HTTPGetConn) Close() error {
	defer func() {
		recover()
	}()
	close(c.CloseChan)
	return nil
}
