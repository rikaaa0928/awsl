package server

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"net/http/pprof"
	"runtime"
	"strconv"

	"github.com/rikaaa0928/awsl/aconn"
	"github.com/rikaaa0928/awsl/alistener"
)

type HTTP struct {
	host   string
	port   int
	uri    string
	cert   string
	key    string
	l      serveListener
	routed bool
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
	case "awsl":
		s.l = &awslAListerWrapper{
			&hbaseAListerWrapper{
				cons: make(chan aconn.AConn, 2*runtime.NumCPU()),
			},
		}
	case "http":
		s.l = &hpAListerWrapper{
			&hbaseAListerWrapper{
				cons: make(chan aconn.AConn, 2*runtime.NumCPU()),
			},
		}
		s.routed = true
	case "pprof":
		mux := http.NewServeMux()
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
		s.l = &pprofAListerWrapper{
			hbaseAListerWrapper: &hbaseAListerWrapper{
				//cons: make(chan aconn.AConn, 2*runtime.NumCPU()),
			},
			mux: mux,
			c:   make(chan struct{}),
		}
		s.routed = true
	default:
	}
	return s
}

func (s *HTTP) Listen() alistener.AListener {
	log.Println("base http listen: " + net.JoinHostPort(s.host, strconv.Itoa(s.port)) + "/" + s.uri)
	if s.routed {
		s.l.setSrv(&http.Server{Addr: net.JoinHostPort(s.host, strconv.Itoa(s.port)), Handler: http.HandlerFunc(s.l.h)})
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
	} else {
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
	}
	return s.l
}

func (s *HTTP) Handler() AHandler {
	return s.l.handler()
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

func (l *hbaseAListerWrapper) handler() AHandler {
	return DefaultAHandler
}
