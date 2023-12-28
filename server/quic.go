package server

//import (
//	"context"
//	"crypto/tls"
//	"errors"
//	"net"
//	"runtime"
//	"strconv"
//
//	"github.com/lucas-clemente/quic-go"
//	"github.com/rikaaa0928/awsl/aconn"
//	"github.com/rikaaa0928/awsl/alistener"
//)
//
//func NewBaseQUIC(listenHost string, listenPort int, cert, key string) BaseQUIC {
//	return BaseQUIC{
//		ip:   listenHost,
//		port: listenPort,
//		cert: cert,
//		key:  key,
//	}
//}
//
//type BaseQUIC struct {
//	ip   string
//	port int
//	cert string
//	key  string
//}
//
//func (s BaseQUIC) Listen() alistener.AListener {
//	c, err := tls.LoadX509KeyPair(s.cert, s.key)
//	if err != nil {
//		panic(err)
//	}
//	listener, err := quic.ListenAddr(net.JoinHostPort(s.ip, strconv.Itoa(s.port)), &tls.Config{Certificates: []tls.Certificate{c}, NextProtos: []string{"awsl-quic"}}, nil)
//	if err != nil {
//		panic(err)
//	}
//	ctx, cancel := context.WithCancel(context.Background())
//	l := &quicAListerWrapper{c: make(chan aconn.AConn, 2*runtime.NumCPU()), cancel: cancel}
//	go func() {
//		for !l.closed {
//			sess, err := listener.Accept(ctx)
//			if err != nil {
//				continue
//			}
//			go func() {
//				for !l.closed {
//					stream, err := sess.AcceptStream(ctx)
//					if err != nil {
//						break
//					}
//					if !l.closed {
//						l.c <- aconn.NewAConn(streamConn{Stream: stream, l: sess.LocalAddr(), r: sess.RemoteAddr()})
//					}
//				}
//			}()
//		}
//		close(l.c)
//	}()
//	return l
//}
//
//func (s BaseQUIC) Handler() AHandler {
//	return DefaultAHandler
//}
//
//type quicAListerWrapper struct {
//	c      chan aconn.AConn
//	closed bool
//	cancel context.CancelFunc
//}
//
//func (l *quicAListerWrapper) Accept(ctx context.Context) (context.Context, aconn.AConn, error) {
//	select {
//	case c, ok := <-l.c:
//		if !ok {
//			return ctx, nil, errors.New("quicAListerWrapper closed")
//		}
//		return ctx, c, nil
//	}
//}
//func (l *quicAListerWrapper) Close() error {
//	l.closed = true
//	l.cancel()
//	return nil
//}
//
//type streamConn struct {
//	quic.Stream
//	l net.Addr
//	r net.Addr
//}
//
//func (c streamConn) RemoteAddr() net.Addr {
//	return c.r
//}
//
//func (c streamConn) LocalAddr() net.Addr {
//	return c.l
//}
