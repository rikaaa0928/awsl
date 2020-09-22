package server

import (
	"context"
	"net"
	"strconv"

	"github.com/rikaaa0928/awsl/aconn"
	"github.com/rikaaa0928/awsl/alistener"
)

func NewBaseTcp(listenHost string, listenPort int) BaseTcp {
	return BaseTcp{
		ip:   listenHost,
		port: listenPort,
	}
}

type BaseTcp struct {
	ip   string
	port int
}

func (s BaseTcp) Listen() alistener.AListener {
	l, e := net.Listen("tcp", net.JoinHostPort(s.ip, strconv.Itoa(s.port)))
	if e != nil {
		panic(e)
	}
	return &baseAListerWrapper{l}
}

func (s BaseTcp) Handler() AHandler {
	return DefaultAHandler
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
