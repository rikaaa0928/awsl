package server

import (
	"context"
	"crypto/tls"
	"net"
	"strconv"

	"github.com/lucas-clemente/quic-go"
	"github.com/rikaaa0928/awsl/aconn"
	"github.com/rikaaa0928/awsl/alistener"
)

func NewBaseQUIC(listenHost string, listenPort int, cert, key string) BaseQUIC {
	return BaseQUIC{
		ip:   listenHost,
		port: listenPort,
		cert: cert,
		key:  key,
	}
}

type BaseQUIC struct {
	ip   string
	port int
	cert string
	key  string
}

func (s BaseQUIC) Listen() alistener.AListener {
	c, err := tls.LoadX509KeyPair(s.cert, s.key)
	if err != nil {
		panic(err)
	}
	listener, err := quic.ListenAddr(net.JoinHostPort(s.ip, strconv.Itoa(s.port)), &tls.Config{Certificates: []tls.Certificate{c}, NextProtos: []string{"awsl-quic"}}, nil)
	if err != nil {
		panic(err)
	}
	return &quicAListerWrapper{listener}
}

func (s BaseQUIC) Handler() AHandler {
	return DefaultAHandler
}

type quicAListerWrapper struct {
	quic.Listener
}

func (l *quicAListerWrapper) Accept(ctx context.Context) (context.Context, aconn.AConn, error) {
	conn, err := l.Listener.Accept(ctx)
	if err != nil {
		return ctx, nil, err
	}
	return ctx, aconn.NewAConn(conn), nil
}
func (l *quicAListerWrapper) Close() error {
	return l.Listener.Close()
}
