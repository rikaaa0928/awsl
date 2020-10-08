package server

import (
	"context"
	"io"
	"log"
	"sync"

	"github.com/rikaaa0928/awsl/aconn"
	"github.com/rikaaa0928/awsl/adialer"
	"github.com/rikaaa0928/awsl/arouter"
	"github.com/rikaaa0928/awsl/consts"
	"github.com/rikaaa0928/awsl/utils"
)

type AHandler func(context.Context, aconn.AConn, arouter.ARouter, adialer.DialerFactory)

var DefaultAHandler AHandler = func(ctx context.Context, sConn aconn.AConn, route arouter.ARouter, getDialer adialer.DialerFactory) {
	superType := ctx.Value(consts.CTXSuperType)
	if superType != nil {
		switch superType.(string) {
		case "udp":
			log.Println("handel udp ", ctx.Value(consts.CTXSuperData))
			return
		default:
		}
		return
	}
	defer sConn.Close()
	ctx = route(ctx, sConn.EndAddr())
	dial := getDialer(ctx)
	if dial == nil {
		log.Println("nil dial")
		return
	}
	ctx, cConn, err := dial(ctx, sConn.EndAddr())
	if err != nil {
		log.Println("dial error: " + err.Error())
		return
	}
	defer cConn.Close()
	w := sync.WaitGroup{}
	w.Add(2)
	go func() {
		buf := utils.GetMem(65536)
		defer utils.PutMem(buf)
		io.CopyBuffer(cConn, sConn, buf)
		w.Done()
	}()
	go func() {
		buf := utils.GetMem(65536)
		defer utils.PutMem(buf)
		io.CopyBuffer(sConn, cConn, buf)
		w.Done()
	}()
	w.Wait()
}
