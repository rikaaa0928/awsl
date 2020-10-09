package adialer

import (
	"context"
	"net"
	"sync"

	"github.com/rikaaa0928/awsl/aconn"
)

var defaultPool = super{
	pool: make(map[string]aconn.AConn),
}

type super struct {
	sync.RWMutex
	pool map[string]aconn.AConn
}

func getSuperConn(srcDst string, conf map[string]interface{}) ADialer {
	defaultPool.Lock()
	defer defaultPool.Unlock()
	var conn aconn.AConn
	var ok bool
	var err error
	conn, ok = defaultPool.pool[srcDst]
	if !ok {
		defaultPool.pool[srcDst] = conn
	}
	return func(ctx context.Context, _ net.Addr) (context.Context, aconn.AConn, error) {
		return ctx, conn, err
	}
}
