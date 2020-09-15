package server

import (
	"context"
	"net"

	"github.com/rikaaa0928/awsl/aconn"
	"github.com/rikaaa0928/awsl/alistener"
)

func NewBaseTcp(listenHost, listenPort, tag string, id int) BaseTcp {
	return BaseTcp{
		ip:   listenHost,
		port: listenPort,
		tag:  tag,
	}
}

type BaseTcp struct {
	ip   string
	port string
	tag  string
}

// Listen server
func (s BaseTcp) Listen() alistener.AListener {
	l, e := net.Listen("tcp", net.JoinHostPort(s.ip, s.port))
	if e != nil {
		panic(e)
	}
	return &baseAListernWrapepr{l}
}

type baseAListernWrapepr struct {
	net.Listener
}

func (l *baseAListernWrapepr) Accept(ctx context.Context) (context.Context, aconn.AConn, error) {
	conn, err := l.Listener.Accept()
	if err != nil {
		return ctx, nil, err
	}
	return ctx, aconn.NewAConn(conn), nil
}
func (l *baseAListernWrapepr) Close() error {
	return l.Listener.Close()
}
