package main

import (
	"context"
	"flag"
	"log"
	"runtime"
	"sync"

	"cloud.google.com/go/profiler"

	"github.com/rikaaa0928/awsl/config"
	"github.com/rikaaa0928/awsl/global"
	"github.com/rikaaa0928/awsl/object"
	"github.com/rikaaa0928/awsl/utils/metrics"
)

func main() {
	cFile := flag.String("c", "/etc/awsl/config.json", "path to config file")
	flag.Parse()
	runtime.GOMAXPROCS(int(float64(runtime.NumCPU()) * 1.4))
	ctx, cancel := context.WithCancel(context.Background())
	conf := config.NewJsonConfig()
	err := conf.Open(*cFile)
	if err != nil {
		panic(err)
	}

	gcp, err := conf.GetBool("gcp")
	if err == nil {
		global.GCP = gcp
	}
	if global.GCP {
		cfg := profiler.Config{
			Service:        "awsl",
			ServiceVersion: "1.0.0",
			// ProjectID must be set if not running on GCP.
			// ProjectID: "my-project",

			// For OpenCensus users:
			// To see Profiler agent spans in APM backend,
			// set EnableOCTelemetry to true
			EnableOCTelemetry: true,
		}

		// Profiler initialization, best done as early as possible.
		if err := profiler.Start(cfg); err != nil {
			// TODO: Handle error.
			log.Println(err)
		}
	}

	ins, err := conf.GetMap("ins")
	if err != nil {
		panic(err)
	}
	go metrics.StartMetrics(conf)
	timeOut, err := conf.GetInt("timeout")
	if err == nil {
		log.Println("timeout : ", timeOut)
		global.TimeOut = timeOut
	}
	wg := &sync.WaitGroup{}
	for k := range ins {
		wg.Add(1)
		go object.DefaultObject(ctx, wg, k, conf)
	}
	wg.Wait()
	cancel()
}
