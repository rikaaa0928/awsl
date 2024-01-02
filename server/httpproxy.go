package server

import (
	"context"
	"github.com/rikaaa0928/awsl/utils/safer"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/rikaaa0928/awsl/aconn"
	"github.com/rikaaa0928/awsl/adialer"
	"github.com/rikaaa0928/awsl/arouter"
	"github.com/rikaaa0928/awsl/utils"
)

type hpAListerWrapper struct {
	*hbaseAListerWrapper
}

func (l *hpAListerWrapper) h(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodConnect {
		w.WriteHeader(http.StatusOK)

		sl := strings.Split(r.Host, ":")
		if len(sl) != 2 {
			log.Println("addr error : " + r.Host)
			http.Error(w, "addr error : "+r.Host, http.StatusBadRequest)
			return
		}
		port, err := strconv.Atoi(sl[1])
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		addr := aconn.NewAddr(sl[0], port, "tcp")

		hijacker, ok := w.(http.Hijacker)
		if !ok {
			log.Println("Hijacking not supported")
			http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
			return
		}
		clientConn, _, err := hijacker.Hijack()
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		ctx, cancel := context.WithCancel(context.Background())
		con := &waitCloseConn{
			Conn:   clientConn,
			ctx:    ctx,
			cancel: cancel,
		}
		ac := aconn.NewAConn(con)
		ac.SetEndAddr(addr)
		l.cons <- ac

		<-con.ctx.Done()
	} else {
		log.Println("http proxy Method not CONNECT.", r.Method, r.Host)
		rHost := ""
		rPort := 80
		var err error
		// host port
		if strings.Contains(r.Host, ":") {
			sl := strings.Split(r.Host, ":")
			if len(sl) != 2 {
				log.Println("addr error : " + r.Host)
				http.Error(w, "addr error : "+r.Host, http.StatusBadRequest)
				return
			}
			rHost = sl[0]
			rPort, err = strconv.Atoi(sl[1])
			if err != nil {
				log.Println(err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		} else {
			rHost = r.Host
		}

		ctx, cancel := context.WithCancel(context.Background())
		addr := aconn.NewAddr(rHost, rPort, "tcp")
		hc := &HTTPGetConn{W: w, R: r, End: addr, ctx: ctx, cancel: cancel}

		l.cons <- hc

		<-hc.ctx.Done()
	}
}

func (l *hpAListerWrapper) handler() AHandler {
	return func(ctx context.Context, sConn aconn.AConn, route arouter.ARouter, getDialer adialer.DialerFactory) {
		hc, ok := sConn.(*aconn.RealConn).AConn.(*HTTPGetConn)
		if ok {
			log.Println("handle http get")
			defer sConn.Close()
			ctx = route(ctx, sConn.EndAddr())
			dial := getDialer(ctx)
			if dial == nil {
				log.Println("nil dial")
				return
			}
			_, cConn, err := dial(ctx, sConn.EndAddr())
			if err != nil {
				log.Println("dial error: " + err.Error())
				return
			}
			rcConn := aconn.CreateRealConn(cConn)
			rcConn.RegisterReader(safer.IOSaferFactory(rcConn.Magic(), true))
			rcConn.RegisterWriter(safer.IOSaferFactory(rcConn.Magic(), false))
			rcConn.RegisterCloser(aconn.NewMetricsMid(ctx, "", "", rcConn.EndAddr().String()).MetricsClose)
			defer rcConn.Close()

			trans := http.Transport{DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return rcConn, nil
			}}
			resp, err := trans.RoundTrip(hc.R)
			if err != nil {
				http.Error(hc.W, err.Error(), http.StatusServiceUnavailable)
				log.Println("object http roudtrip error", sConn.EndAddr(), err)
				return
			}
			defer resp.Body.Close()
			utils.CopyHeader(hc.W.Header(), resp.Header)
			hc.W.WriteHeader(resp.StatusCode)
			buf := utils.GetMem(65536)
			defer utils.PutMem(buf)
			io.CopyBuffer(hc.W, resp.Body, buf)
			// n, err := io.CopyBuffer(hc.W, resp.Body, buf)
			//log.Println("http", n, err, cConn.EndAddr())
			hc.Close()
		} else {
			DefaultAHandler(ctx, sConn, route, getDialer)
		}
	}
}

var _ aconn.AConn = &HTTPGetConn{}

// HTTPGetConn HTTPGetConn
type HTTPGetConn struct {
	W   http.ResponseWriter
	R   *http.Request
	End net.Addr
	net.Conn
	ctx    context.Context
	cancel context.CancelFunc
}

func (c *HTTPGetConn) Magic() *uint32 {
	return nil
}

func (c *HTTPGetConn) SetMagic(u uint32) {
}

func (c *HTTPGetConn) Close() error {
	var err error
	if c.Conn != nil {
		err = c.Conn.Close()
	}
	c.cancel()
	return err
}

func (c *HTTPGetConn) EndAddr() net.Addr {
	return c.End
}

func (c *HTTPGetConn) SetEndAddr(addr net.Addr) {
	c.End = addr
}
