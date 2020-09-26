package server

import (
	"context"
	"net"
	"net/http"

	"github.com/rikaaa0928/awsl/aconn"
	"golang.org/x/net/websocket"
)

type awslAListerWrapper struct {
	*hbaseAListerWrapper
}

func (l *awslAListerWrapper) handle(conn *websocket.Conn) {
	ctx, cancel := context.WithCancel(context.Background())
	ac := &awslConn{
		Conn:   conn,
		ctx:    ctx,
		cancel: cancel,
	}

	l.cons <- aconn.NewAConn(ac)
	<-ac.ctx.Done()
}

func (l *awslAListerWrapper) h(w http.ResponseWriter, r *http.Request) {
	websocket.Handler(l.handle).ServeHTTP(w, r)
}

type awslConn struct {
	net.Conn
	ctx    context.Context
	cancel context.CancelFunc
}

func (c *awslConn) Close() error {
	err := c.Conn.Close()
	c.cancel()
	return err
}
