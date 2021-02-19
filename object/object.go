package object

import (
	"context"
	"log"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/rikaaa0928/awsl/aconn"
	"github.com/rikaaa0928/awsl/adialer"
	"github.com/rikaaa0928/awsl/alistener"
	"github.com/rikaaa0928/awsl/arouter"
	"github.com/rikaaa0928/awsl/config"
	"github.com/rikaaa0928/awsl/global"
	"github.com/rikaaa0928/awsl/server"
)

type Object func(context.Context, *sync.WaitGroup, string, config.Configs)

var DefaultObject Object = func(ctx context.Context, wg *sync.WaitGroup, tag string, c config.Configs) {
	//closed := false
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
	alistener.DefaultAcceptMids(ctx, l, typ, tag, conf)
	//go func(closed *bool) {
	//	select {
	//	case <-ctx.Done():
	//		*closed = true
	//		l.Close()
	//	}
	//}(&closed)
	ctx = context.WithValue(ctx, global.CTXInTag, tag)
	ctx = context.WithValue(ctx, global.CTXInType, typ)
	for /*!closed*/ {
		select {
		case <-ctx.Done():
			l.Close()
			break
		default:

		}
		ctx, ac, err := l.Accept(ctx)
		if err != nil {
			log.Println("accept error: ", err)
			continue
		}
		go func() {
			if global.Tracing {
				tracer := otel.Tracer("awsl")
				var span trace.Span
				ctx, span = tracer.Start(ctx, "object_go_routine")
				defer span.End()
			}
			rc := aconn.CreateRealConn(ac)
			rc.RegisterCloser(aconn.NewMetricsMid(ctx, tag, typ, rc.EndAddr().String()).MetricsClose)
			outsConf, err := c.GetMap("outs")
			if err != nil {
				log.Println("c.GetMap('outs'), err: ", err)
				return
			}
			handle(ctx, rc, arouter.NewRouter(c), adialer.NewFactory(outsConf))
		}()
	}
	wg.Done()
}
