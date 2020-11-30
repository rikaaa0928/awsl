package server

import (
	"context"
	"io"
	"log"
	"sync"

	"github.com/rikaaa0928/awsl/aconn"
	"github.com/rikaaa0928/awsl/adialer"
	"github.com/rikaaa0928/awsl/arouter"
	"github.com/rikaaa0928/awsl/utils"
)

type AHandler func(context.Context, aconn.AConn, arouter.ARouter, adialer.DialerFactory)

var DefaultAHandler AHandler = func(ctx context.Context, sConn aconn.AConn, route arouter.ARouter, getDialer adialer.DialerFactory) {
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
	rcConn.RegisterCloser(aconn.NewMetricsMidForOut(ctx, rcConn.EndAddr().String()).MetricsClose)
	defer rcConn.Close()
	w := sync.WaitGroup{}
	w.Add(2)
	// debug := strings.Contains(sConn.EndAddr().String(), "steam")
	debug := false
	go func() {
		buf := utils.GetMem(65536)
		defer utils.PutMem(buf)
		if debug {
			n, err := io.CopyBuffer(rcConn, sConn, buf)
			log.Println("io.CopyBuffer(cConn, sConn, buf)", sConn.EndAddr().String(), n, err)
		} else {
			io.CopyBuffer(rcConn, sConn, buf)
		}
		w.Done()
	}()
	go func() {
		buf := utils.GetMem(65536)
		defer utils.PutMem(buf)
		if debug {
			n, err := io.CopyBuffer(sConn, rcConn, buf)
			log.Println("io.CopyBuffer(sConn, cConn, buf)", sConn.EndAddr().String(), n, err)
		} else {
			io.CopyBuffer(sConn, rcConn, buf)
		}
		w.Done()
	}()
	w.Wait()
}
