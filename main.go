package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"github.com/rikaaa0928/awsl/utils/metrics"
	"log"
	"runtime"
	"sync"

	"github.com/rikaaa0928/awsl/config"
	"github.com/rikaaa0928/awsl/global"
	"github.com/rikaaa0928/awsl/object"
)

//go:embed version
var version string

func main() {
	checkVer := flag.Bool("v", false, "version")
	cFile := flag.String("c", "/etc/awsl/config.json", "path to config file")
	flag.Parse()
	if *checkVer {
		fmt.Println(version)
		return
	}
	runtime.GOMAXPROCS(int(float64(runtime.NumCPU()) * 1.4))
	ctx, cancel := context.WithCancel(context.Background())
	conf := config.NewJsonConfig()
	err := conf.Open(*cFile)
	if err != nil {
		panic(err)
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
