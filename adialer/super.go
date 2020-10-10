package adialer

import (
	"context"
	"net"
	"sync"

	"github.com/rikaaa0928/awsl/aconn"
)

var defaultPool = super{
	pools: make(map[string]map[string]aconn.AConn),
}

type super struct {
	sync.RWMutex
	pools map[string]map[string]aconn.AConn
}

func getSuperConn(tag, src, dst string, conf map[string]interface{}) ADialer {
	defaultPool.RLock()
	var conn aconn.AConn
	var ok bool
	var err error
	pool, ok := defaultPool.pools[tag]
	if !ok {
		defaultPool.RUnlock()
		defaultPool.Lock()
		defaultPool.pools[tag] = make(map[string]aconn.AConn)
		defaultPool.Unlock()
		defaultPool.RLock()
	}
	conn, ok = pool[src+"-"+dst]
	if !ok {
		defaultPool.RUnlock()
		defaultPool.Lock()
		switch conf["type"] {
		case "free":

		default:
		}
		defaultPool.pools[tag][src+"-"+dst] = conn
		defaultPool.Unlock()
		defaultPool.RLock()
	}
	defer defaultPool.RUnlock()
	return func(ctx context.Context, _ net.Addr) (context.Context, aconn.AConn, error) {
		return ctx, conn, err
	}
}
