package main

import (
	"context"
	"log"
	"sync"

	"github.com/rikaaa0928/awsl/config"
	"github.com/rikaaa0928/awsl/object"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	conf := config.NewJsonConfig()
	err := conf.Open("./test/conf.json")
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
