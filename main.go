package main

import (
	"context"
	"flag"
	"log"
	"runtime"
	"sync"

	"github.com/rikaaa0928/awsl/config"
	"github.com/rikaaa0928/awsl/object"
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
	ins, err := conf.GetMap("ins")
	log.Println(len(ins))
	wg := &sync.WaitGroup{}
	for k := range ins {
		wg.Add(1)
		go object.DefaultObject(ctx, wg, k, conf)
	}
	wg.Wait()
	cancel()
}
