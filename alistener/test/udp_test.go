package test_test

import (
	"context"
	"encoding/hex"
	"log"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/rikaaa0928/awsl/alistener"
	"github.com/rikaaa0928/awsl/server"
	"github.com/txthinking/socks5"
)

func TestUDP(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(ctx context.Context) {
		l := alistener.NewRealListener(server.NewBaseTcp("127.0.0.1", 4888).Listen())
		l.RegisterAcceptor(alistener.NewSocksAcceptMid(ctx, map[string]interface{}{"host": "127.0.0.1", "port": 4888.0}))
		go func() {
			<-ctx.Done()
			l.Close()
		}()
	main:
		for {
			select {
			case <-ctx.Done():
				break main
			default:

			}
			_, _, err := l.Accept(ctx)
			if err != nil && err != alistener.ErrUDP {
				t.Log(err)
				continue
			}
			t.Log(1)
		}
		wg.Done()
	}(ctx)
	time.Sleep(time.Second)

	c, err := socks5.NewClient("127.0.0.1:4888", "", "", 0, 60)
	if err != nil {
		panic(err)
	}
	conn, err := c.Dial("udp", "8.8.8.8:53")
	if err != nil {
		panic(err)
	}
	b, err := hex.DecodeString("0001010000010000000000000a74787468696e6b696e6703636f6d0000010001")
	if err != nil {
		panic(err)
	}
	t.Log(b)
	if _, err := conn.Write(b); err != nil {
		panic(err)
	}
	b2 := make([]byte, 2048)
	n, err := conn.Read(b2)
	if err != nil {
		panic(err)
	}
	b2 = b2[:n]
	b2 = b2[len(b2)-4:]
	log.Println("result", net.IPv4(b2[0], b2[1], b2[2], b2[3]))
	if _, err := conn.Write(b); err != nil {
		panic(err)
	}
	b2 = make([]byte, 2048)
	n, err = conn.Read(b2)
	if err != nil {
		panic(err)
	}
	b2 = b2[:n]
	b2 = b2[len(b2)-4:]
	log.Println("result", net.IPv4(b2[0], b2[1], b2[2], b2[3]))
	cancel()
	wg.Wait()
}
