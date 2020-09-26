package server

import (
	"context"
	"errors"
	"net"
	"net/http"
	"runtime"
	"strconv"

	"github.com/rikaaa0928/awsl/aconn"
	"github.com/rikaaa0928/awsl/alistener"
)

const (
	TypeHTTP = iota
	TypeWebsocket
)

type HTTP struct {
	host string
	port int
	uri  string
	cert string
	key  string
	l    serveListener
	typ  int
}

func NewHTTPServer(typ, host, uri, cert, key string, port int) *HTTP {
	s := &HTTP{
		host: host,
		port: port,
		uri:  uri,
		cert: cert,
		key:  key,
	}
	switch typ {
	case "h2c":
		s.typ = TypeHTTP
		s.l = &h2cAListerWrapper{
			&hbaseAListerWrapper{
				cons: make(chan aconn.AConn, 2*runtime.NumCPU()),
			},
		}
	case "awsl":
		s.typ = TypeWebsocket
		s.l = &awslAListerWrapper{
			&hbaseAListerWrapper{
				cons: make(chan aconn.AConn, 2*runtime.NumCPU()),
			},
		}
	default:
	}
	return s
}

func (s *HTTP) Listen() alistener.AListener {
	mux := http.NewServeMux()
	mux.HandleFunc("/"+s.uri, s.l.h)
	s.l.setSrv(&http.Server{Addr: net.JoinHostPort(s.host, strconv.Itoa(s.port)), Handler: mux})
	go func() {
		if len(s.cert) == 0 || len(s.key) == 0 {
			err := s.l.srv().ListenAndServe()
			if err != nil {
				panic("ListenAndServe: " + err.Error())
			}
		} else {
			err := s.l.srv().ListenAndServeTLS(s.cert, s.key)
			if err != nil {
				panic("ListenAndServe: " + err.Error())
			}
		}
	}()
	return s.l
}

func (s *HTTP) Handler() AHandler {
	return DefaultAHandler
}

type hbaseAListerWrapper struct {
	cons chan aconn.AConn
	s    *http.Server
}

func (l *hbaseAListerWrapper) srv() *http.Server {
	return l.s
}

func (l *hbaseAListerWrapper) setSrv(s *http.Server) {
	l.s = s
}

func (l *hbaseAListerWrapper) h(w http.ResponseWriter, r *http.Request) {

}

func (l *hbaseAListerWrapper) Accept(ctx context.Context) (context.Context, aconn.AConn, error) {
	conn, ok := <-l.cons
	if ok {
		return ctx, conn, nil
	}
	return ctx, nil, errors.New("h2c server closed")
}

func (l *hbaseAListerWrapper) Close() error {
	l.srv().Shutdown(context.Background())
	close(l.cons)
	return nil
}
