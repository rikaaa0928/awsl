package main

import (
	"context"
	"log"
	"os"

	"github.com/rikaaa0928/awsl/aconn"
	"github.com/rikaaa0928/awsl/adialer"
	"github.com/rikaaa0928/awsl/alistener"
	"github.com/rikaaa0928/awsl/arouter"
	"github.com/rikaaa0928/awsl/server"
)

func main() {
	log.SetOutput(os.Stderr)
	s := server.NewBaseTcp("127.0.0.1", "12345")
	handle := s.Handler()
	l := alistener.NewRealListener(s.Listen())
	l.RegisterAcceptor(alistener.NewSocksAcceptMid("socks"))
	for {
		ctx, c, err := l.Accept(context.Background())
		if err != nil {
			log.Println(err)
		}
		go func() {
			rc := aconn.CreateRealConn(c)
			handle(ctx, rc, arouter.NopRouter, adialer.DefaultFactory)
		}()
	}
}
