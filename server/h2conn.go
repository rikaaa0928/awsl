package server

import (
	"context"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"runtime"
	"strconv"
	"sync"

	"github.com/rikaaa0928/awsl/aconn"
	"github.com/rikaaa0928/awsl/adialer"
	"github.com/rikaaa0928/awsl/alistener"
	"github.com/rikaaa0928/awsl/arouter"
	"github.com/rikaaa0928/awsl/utils"
)

type H2Conn struct {
	host string
	port int
	uri  string
	cert string
	key  string
	srv  http.Server
}

func (s *H2Conn) Listen() alistener.AListener {
	mux := http.NewServeMux()
	ctx, cancel := context.WithCancel(context.Background())
	l := &h2cAListerWrapper{
		cons:   make(chan aconn.AConn, 2*runtime.NumGoroutine()),
		ctx:    ctx,
		cancel: cancel,
	}
	mux.HandleFunc("/"+s.uri+"/", l.handel)
	s.srv = http.Server{Addr: net.JoinHostPort(s.host, strconv.Itoa(s.port)), Handler: mux}
	return l
}
func (s *H2Conn) Handler() AHandler {
	return func(ctx context.Context, sConn aconn.AConn, route arouter.ARouter,
		getDialer adialer.DialerFactory) {
		defer sConn.Close()
		ctx = route(ctx, sConn.EndAddr())
		dial := getDialer(ctx)
		cConn, err := dial(ctx, sConn.EndAddr())
		if err != nil {
			log.Println(err)
			return
		}
		defer cConn.Close()
		w := sync.WaitGroup{}
		w.Add(2)
		go func() {
			buf := utils.GetMem(65536)
			defer utils.PutMem(buf)
			io.CopyBuffer(cConn, sConn, buf)
			w.Done()
		}()
		go func() {
			buf := utils.GetMem(65536)
			defer utils.PutMem(buf)
			io.CopyBuffer(sConn, cConn, buf)
			w.Done()
		}()
		w.Wait()
	}
}

type h2cAListerWrapper struct {
	cons   chan aconn.AConn
	ctx    context.Context
	cancel context.CancelFunc
	srv    http.Server
}

func (l *h2cAListerWrapper) handel(w http.ResponseWriter, req *http.Request) {

}

func (l *h2cAListerWrapper) Accept(ctx context.Context) (context.Context, aconn.AConn, error) {
	select {
	case conn := <-l.cons:
		return ctx, conn, nil
	case <-ctx.Done():
		l.cancel()
		return ctx, nil, errors.New("h2c server closed")
	case <-l.ctx.Done():
	}
	return ctx, nil, errors.New("h2c server closed")
}

func (l *h2cAListerWrapper) Close() error {
	l.cancel()
	return nil
}
