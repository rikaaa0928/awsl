package test

import (
	"context"
	"log"
	"net"
	"runtime"
	"sync"
	"testing"

	"github.com/rikaaa0928/awsl/config"
	"github.com/rikaaa0928/awsl/object"
)

func TestClient(t *testing.T) {
	runtime.GOMAXPROCS(int(float64(runtime.NumCPU()) * 1.4))
	ctx, cancel := context.WithCancel(context.Background())
	conf := config.NewJsonConfig()
	err := conf.Open("./client.json")
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

func TestSplitHP(t *testing.T) {
	t.Log(net.SplitHostPort("[::]:123"))
}
