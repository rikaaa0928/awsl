package server

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/rikaaa0928/awsl/aconn"
	"github.com/rikaaa0928/awsl/consts"
	"github.com/rikaaa0928/awsl/utils"
	"github.com/rikaaa0928/awsl/utils/rwconn"
)

type h2cAListerWrapper struct {
	*hbaseAListerWrapper
}

func (l *h2cAListerWrapper) h(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "PUT required.", http.StatusBadRequest)
		return
	}

	addr := aconn.AddrInfo{}
	addrCookie, err := r.Cookie("addr")
	if err != nil {
		http.Error(w, "address required.", http.StatusBadRequest)
		return
	}
	addrStr, err := url.QueryUnescape(addrCookie.Value)
	if err != nil {
		http.Error(w, "address required.", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal([]byte(addrStr), &addr)
	if err != nil {
		http.Error(w, "address required.", http.StatusBadRequest)
		return
	}
	auth, err := r.Cookie("auth")
	if err != nil {
		http.Error(w, "wtf?", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	f, ok := w.(http.Flusher)
	if !ok {
		log.Println("not flusher")
		http.Error(w, "flush required.", http.StatusHTTPVersionNotSupported)
		return
	}
	f.Flush()
	wr, ww := io.Pipe()
	rr, rw := io.Pipe()
	defer func() {
		wr.Close()
		ww.Close()
		rr.Close()
		rw.Close()
	}()
	c := &h2cConn{RWConn: rwconn.NewRWConn(ww, rr), auth: auth.Value, addr: addr}
	l.cons <- c
	wait := sync.WaitGroup{}
	wait.Add(2)
	go func() {
		buf := utils.GetMem(65536)
		defer utils.PutMem(buf)
		io.CopyBuffer(rewrite{w}, wr, buf)
		wait.Done()
	}()
	go func() {
		buf := utils.GetMem(65536)
		defer utils.PutMem(buf)
		io.CopyBuffer(rw, r.Body, buf)
		wait.Done()
	}()
	wait.Wait()
}

func (l *h2cAListerWrapper) Accept(ctx context.Context) (context.Context, aconn.AConn, error) {
	conn, ok := <-l.cons
	if ok {
		ctx = context.WithValue(ctx, consts.CTXReceiveAuth, conn.(*h2cConn).auth)
		return ctx, conn, nil
	}
	return ctx, nil, errors.New("h2c server closed")
}

type h2cConn struct {
	*rwconn.RWConn
	auth string
	addr net.Addr
}

func (c *h2cConn) Read(b []byte) (n int, err error) {
	return c.RWConn.Read(b)
}

func (c *h2cConn) Write(b []byte) (n int, err error) {
	return c.RWConn.Write(b)
}

func (c *h2cConn) Close() error {
	return c.RWConn.Close()
}
func (c *h2cConn) LocalAddr() net.Addr {
	return nil
}
func (c *h2cConn) RemoteAddr() net.Addr {
	return nil
}
func (c *h2cConn) SetDeadline(t time.Time) error {
	return nil
}
func (c *h2cConn) SetReadDeadline(t time.Time) error {
	return nil
}
func (c *h2cConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func (c *h2cConn) EndAddr() net.Addr {
	return c.addr
}
func (c *h2cConn) SetEndAddr(addr net.Addr) {
	c.addr = addr
}

type rewrite struct {
	w http.ResponseWriter
}

func (w rewrite) Write(b []byte) (n int, err error) {
	n, err = w.w.Write(b)
	w.w.(http.Flusher).Flush()
	return
}
