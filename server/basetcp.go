package server

import (
	"context"
	"io"
	"log"
	"net"
	"sync"

	"github.com/rikaaa0928/awsl/aconn"
	"github.com/rikaaa0928/awsl/adialer"
	"github.com/rikaaa0928/awsl/alistener"
	"github.com/rikaaa0928/awsl/arouter"
	"github.com/rikaaa0928/awsl/utils"
)

func NewBaseTcp(listenHost, listenPort string) BaseTcp {
	return BaseTcp{
		ip:   listenHost,
		port: listenPort,
	}
}

type BaseTcp struct {
	ip   string
	port string
}

// Listen server
func (s BaseTcp) Listen() alistener.AListener {
	l, e := net.Listen("tcp", net.JoinHostPort(s.ip, s.port))
	if e != nil {
		panic(e)
	}
	return &baseAListerWrapper{l}
}

// Listen server
func (s BaseTcp) Handler() AHandler {
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

type baseAListerWrapper struct {
	net.Listener
}

func (l *baseAListerWrapper) Accept(ctx context.Context) (context.Context, aconn.AConn, error) {
	conn, err := l.Listener.Accept()
	if err != nil {
		return ctx, nil, err
	}
	return ctx, aconn.NewAConn(conn), nil
}
func (l *baseAListerWrapper) Close() error {
	return l.Listener.Close()
}
