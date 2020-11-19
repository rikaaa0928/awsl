package adialer

import (
	"context"
	"crypto/tls"
	"net"
	"strconv"
	"sync"

	"github.com/lucas-clemente/quic-go"
	"github.com/rikaaa0928/awsl/aconn"
	"github.com/rikaaa0928/awsl/consts"
	"github.com/rikaaa0928/awsl/utils/ctxdatamap"
)

var QUICConf = struct {
	sync.RWMutex
	Inited     bool
	remoteHost string
	remotePort string
	auth       string
	skipVerify bool
	sess       quic.Session
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
		var session quic.Session
		var err error
		var stream quic.Stream
		if QUICConf.sess == nil {
			session, err = quic.DialAddr(net.JoinHostPort(QUICConf.remoteHost, QUICConf.remotePort), tlsConf, nil)
			if err != nil {
				return ctx, nil, err
			}
			QUICConf.sess = session
		} else {
			session = QUICConf.sess
		}
		for i := 0; i < 2; i++ {
			stream, err = session.OpenStreamSync(ctx)
			if err != nil {
				var err2 error
				session, err2 = quic.DialAddr(net.JoinHostPort(QUICConf.remoteHost, QUICConf.remotePort), tlsConf, nil)
				if err2 != nil {
					QUICConf.sess = nil
					return ctx, nil, err
				}
				QUICConf.sess = session
			} else {
				break
			}
		}
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
