package alistener

import (
	"context"
	"io"

	"github.com/rikaaa0928/awsl/aconn"
)

func NewRealListener(l AListener) *RealListener {
	return &RealListener{Closer: l, a: l.Accept}
}

type RealListener struct {
	a Acceptor
	io.Closer
}

func (l *RealListener) RegisterAcceptor(mid AcceptMid) {
	l.a = mid(l.a)
}

func (l *RealListener) Accept(ctx context.Context) (context.Context, aconn.AConn, error) {
	return l.a(ctx)
}

func (l *RealListener) Close() error {
	return l.Closer.Close()
}
