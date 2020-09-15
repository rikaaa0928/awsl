package alistener

import (
	"context"
	"io"

	"github.com/rikaaa0928/awsl/aconn"
)

type RealListener struct {
	a Accepter
	io.Closer
}

func (l *RealListener) RegisterAccepter(mid AcceptMid) {
	l.a = mid(l.a)
}

func (l *RealListener) Accept(ctx context.Context) (context.Context, aconn.AConn, error) {
	return l.a(ctx)
}

func (l *RealListener) Close() error {
	return l.Closer.Close()
}
