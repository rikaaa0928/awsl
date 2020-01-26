package servers

import (
	"errors"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/Evi1/awsl/config"
	"github.com/Evi1/awsl/model"
)

// NewHTTP NewHTTP
func NewHTTP(listenHost, listenPort string, connsSize int) *HTTPServer {
	a := &HTTPServer{
		IP:   listenHost,
		Port: listenPort,
		Listener: &AWSListener{
			Conns: make(chan net.Conn, connsSize),
			IP:    listenHost,
			Port:  listenPort,
		},
		ConnNum: make(chan int, 1),
	}
	a.ConnNum <- 0
	return a
}

// HTTPServer HTTPServer
type HTTPServer struct {
	IP       string
	Port     string
	Listener *AWSListener
	ConnNum  chan int
	Max      int
}

// Listen server
func (s *HTTPServer) Listen() net.Listener {
	server := &http.Server{
		Addr: s.IP + ":" + s.Port,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodConnect {
				w.WriteHeader(http.StatusOK)
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

				con := &httpConn{Conn: clientConn, CloseChan: make(chan int8), addr: addr}
				s.Listener.Conns <- con

				if config.Debug {
					num := <-s.ConnNum
					num++
					if num > s.Max {
						s.Max = num
						log.Println("max conn: " + strconv.Itoa(num))
					}
					s.ConnNum <- num
				}

				<-con.CloseChan

				if config.Debug {
					num := <-s.ConnNum
					num--
					if num == 0 {
						log.Println("Connection clear")
					}
					s.ConnNum <- num
				}

			} else {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}),
	}
	go server.ListenAndServe()
	return s.Listener
}

// ReadRemote server
func (s *HTTPServer) ReadRemote(c net.Conn) (model.ANetAddr, error) {
	conn, ok := c.(*httpConn)
	if !ok {
		return model.ANetAddr{}, errors.New("conn not httpConn")
	}
	return conn.addr, nil
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
