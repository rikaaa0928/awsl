package adialer

import (
	"context"
	"crypto/tls"
	"github.com/lucas-clemente/quic-go"
	"github.com/rikaaa0928/awsl/aconn"
	"github.com/rikaaa0928/awsl/consts"
	"github.com/rikaaa0928/awsl/utils/ctxdatamap"
	"net"
	"strconv"
	"sync"
)

var QUICConf = struct {
	sync.RWMutex
	Inited     bool
	remoteHost string
	remotePort string
	auth       string
	skipVerify bool
}{}

func NewQUIC(conf map[string]interface{}) ADialer {
	QUICConf.RLock()
	if !QUICConf.Inited {
		QUICConf.RUnlock()
		QUICConf.Lock()
		if !QUICConf.Inited {
			if skip, ok := conf["skipVerify"].(bool); ok && skip {
				QUICConf.skipVerify = true
			}
			QUICConf.remoteHost = conf["host"].(string)
			QUICConf.remotePort = strconv.Itoa(int(conf["port"].(float64)))
			QUICConf.auth = conf["auth"].(string)
			QUICConf.Inited = true
		}
		QUICConf.Unlock()
	} else {
		QUICConf.RUnlock()
	}
	return func(ctx context.Context, addr net.Addr) (context.Context, aconn.AConn, error) {
		tlsConf := &tls.Config{
			InsecureSkipVerify: QUICConf.skipVerify,
			NextProtos:         []string{"awsl-quic"},
		}
		session, err := quic.DialAddr(net.JoinHostPort(QUICConf.remoteHost, QUICConf.remotePort), tlsConf, nil)
		if err != nil {
			return ctx, nil, err
		}

		stream, err := session.OpenStreamSync(context.Background())
		if err != nil {
			return ctx, nil, err
		}
		conn := aconn.NewAConn(streamConn{Stream: stream, l: session.LocalAddr(), r: session.RemoteAddr()})
		conn.SetEndAddr(addr)
		ctx = ctxdatamap.Set(ctx, consts.TransferAuth, QUICConf.auth)
		return ctx, conn, nil
	}
}

type streamConn struct {
	quic.Stream
	l net.Addr
	r net.Addr
}

func (c streamConn) RemoteAddr() net.Addr {
	return c.r
}

func (c streamConn) LocalAddr() net.Addr {
	return c.l
}
