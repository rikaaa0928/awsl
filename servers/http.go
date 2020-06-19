package servers

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/Evi1/awsl/config"
	"github.com/Evi1/awsl/model"
	"github.com/Evi1/awsl/servers/manage"
	"github.com/Evi1/awsl/tools"
)

// NewHTTP NewHTTP
func NewHTTP(ctx context.Context, conf model.In, id int) *HTTP {
	a := &HTTP{
		IP:        conf.Host,
		Port:      conf.Port,
		Conns:     make(chan net.Conn, conf.HTTP.Chan),
		closeWait: tools.NewCloseWait(ctx),
		id:        id,
		tag:       conf.Tag,
	}
	return a
}

// HTTP HTTP
type HTTP struct {
	IP        string
	Port      string
	Conns     chan net.Conn
	Max       int
	closeWait *tools.CloseWait
	Srv       http.Server
	id        int
	tag       string
}

// Listen server
func (s *HTTP) Listen() net.Listener {
	s.Srv = http.Server{
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

				//con := &httpConn{Conn: clientConn, CloseChan: make(chan int8), addr: addr}
				con := &httpConn{Conn: clientConn, closeWait: tools.NewCloseWait(s.closeWait.Ctx), addr: addr}
				s.Conns <- con
				if config.Manage > 0 {
					manage.NewConnectionCount(s.IDTag())
				}
				con.closeWait.WaitClose()
				if config.Manage > 0 {
					manage.ConnectionCloseCount(s.id)
				}
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

				conn := &HTTPGetConn{W: w, R: r, addr: addr, closeWait: tools.NewCloseWait(s.closeWait.Ctx)}
				s.Conns <- conn
				if config.Manage > 0 {
					manage.NewConnectionCount(s.IDTag())
				}
				conn.closeWait.WaitClose()
				if config.Manage > 0 {
					manage.ConnectionCloseCount(s.id)
				}
			}
		}),
	}
	go s.Srv.ListenAndServe()
	return s
}

// ReadRemote server
func (s *HTTP) ReadRemote(c net.Conn) (model.ANetAddr, error) {
	conn, ok := c.(model.AWSLConn)
	if !ok {
		return model.ANetAddr{}, errors.New("conn not httpConn")
	}
	return conn.GetAddr(), nil
}

// IDTag id and tag
func (s *HTTP) IDTag() (int, string) {
	return s.id, s.tag
}

// Accept Accept
func (s *HTTP) Accept() (net.Conn, error) {
	select {
	case conn := <-s.Conns:
		return conn, nil
	case <-s.closeWait.WaitClose():
	}
	return nil, errors.New("http server closed")
}

// Close Close
func (s *HTTP) Close() error {
	s.closeWait.Close()
	return s.Srv.Shutdown(context.Background())
}

// Addr Addr
func (s *HTTP) Addr() net.Addr {
	return &net.IPAddr{
		IP:   net.ParseIP(s.IP),
		Zone: "",
	}
}

type httpConn struct {
	net.Conn
	addr model.ANetAddr
	//CloseChan chan int8
	closeWait *tools.CloseWait
}

func (c *httpConn) Close() error {
	err := c.Conn.Close()
	c.closeWait.Close()
	//c.CloseChan <- 1
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
	//CloseChan chan int8
	closeWait *tools.CloseWait
}

// GetAddr GetAddr
func (c *HTTPGetConn) GetAddr() model.ANetAddr {
	return c.addr
}

// Close Close
func (c *HTTPGetConn) Close() error {
	/*defer func() {
		recover()
	}()
	close(c.CloseChan)*/
	c.closeWait.Close()
	return nil
}
