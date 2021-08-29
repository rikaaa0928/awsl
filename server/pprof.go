package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/rikaaa0928/awsl/aconn"
)

type pprofAListerWrapper struct {
	*hbaseAListerWrapper
	c   chan struct{}
	mux *http.ServeMux
}

func (l *pprofAListerWrapper) h(w http.ResponseWriter, r *http.Request) {
	l.mux.ServeHTTP(w, r)
}

func (l *pprofAListerWrapper) Accept(ctx context.Context) (context.Context, aconn.AConn, error) {
	<-l.c
	return ctx, nil, errors.New("pprof listen closed")
}

func (l *pprofAListerWrapper) Close() error {
	close(l.c)
	return nil
}
