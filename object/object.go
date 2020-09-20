package object

import (
	"context"
	"log"

	"github.com/rikaaa0928/awsl/aconn"
	"github.com/rikaaa0928/awsl/adialer"
	"github.com/rikaaa0928/awsl/alistener"
	"github.com/rikaaa0928/awsl/arouter"
	"github.com/rikaaa0928/awsl/config"
	"github.com/rikaaa0928/awsl/server"
)

type Object func(string, config.Configs)

var DefaultObject Object = func(tag string, c config.Configs) {
	typ, err := c.GetString("ins", tag, "type")
	if err != nil {
		panic(err)
	}
	conf, err := c.GetMap("ins", tag)
	if err != nil {
		panic(err)
	}
	s, err := server.NewServer(typ, conf)
	if err != nil {
		panic(err)
	}
	//s := server.NewBaseTcp(host, int(port))
	handle := s.Handler()
	l := alistener.NewRealListener(s.Listen())
	alistener.DefaultAcceptMids(l, typ, tag)
	//l.RegisterAcceptor(alistener.NewSocksAcceptMid("socks"))
	for {
		ctx, c, err := l.Accept(context.Background())
		if err != nil {
			log.Println(err)
		}
		go func() {
			rc := aconn.CreateRealConn(c)
			handle(ctx, rc, arouter.NopRouter, adialer.TestFactory)
		}()
	}
}
