package server

import (
	"context"
	"io"
	"log"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/rikaaa0928/awsl/aconn"
	"github.com/rikaaa0928/awsl/adialer"
	"github.com/rikaaa0928/awsl/arouter"
	"github.com/rikaaa0928/awsl/global"
	"github.com/rikaaa0928/awsl/utils"
)

type AHandler func(context.Context, aconn.AConn, arouter.ARouter, adialer.DialerFactory)

func inStringSlice(s string, slice []string) bool {
	for _, v := range slice {
		if s == v {
			return true
		}
	}
	return false
}

var DefaultAHandler AHandler = func(ctx context.Context, sConn aconn.AConn, route arouter.ARouter, getDialer adialer.DialerFactory) {
	var tracer trace.Tracer
	if global.Tracing && !inStringSlice(ctx.Value(global.CTXInTag).(string), global.TraceBypassTags) {
		tracer = otel.Tracer("awsl")
		var span trace.Span
		ctx, span = tracer.Start(ctx, "default_handler")
		defer span.End()
	}
	defer sConn.Close()
	ctx = route(ctx, sConn.EndAddr())
	dial := getDialer(ctx)
	if dial == nil {
		log.Println("nil dial")
		return
	}
	var cConn aconn.AConn
	var err error
	ctx, cConn, err = dial(ctx, sConn.EndAddr())
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
	go func(ctx context.Context) {
		if tracer != nil {
			_, span := tracer.Start(ctx, "go_io.CopyBuffer_c_s")
			defer span.End()
		}
		defer sConn.Close()
		defer rcConn.Close()
		buf := utils.GetMem(65536)
		defer utils.PutMem(buf)
		if debug {
			n, err := io.CopyBuffer(rcConn, sConn, buf)
			log.Println("io.CopyBuffer(cConn, sConn, buf)", sConn.EndAddr().String(), n, err)
		} else {
			io.CopyBuffer(rcConn, sConn, buf)
		}
		w.Done()
	}(ctx)
	go func(ctx context.Context) {
		if tracer != nil {
			_, span := tracer.Start(ctx, "go_io.CopyBuffer_s_c")
			defer span.End()
		}
		defer sConn.Close()
		defer rcConn.Close()
		buf := utils.GetMem(65536)
		defer utils.PutMem(buf)
		if debug {
			n, err := io.CopyBuffer(sConn, rcConn, buf)
			log.Println("io.CopyBuffer(sConn, cConn, buf)", sConn.EndAddr().String(), n, err)
		} else {
			io.CopyBuffer(sConn, rcConn, buf)
		}
		w.Done()
	}(ctx)
	w.Wait()
}
