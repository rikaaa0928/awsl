package test

import (
	"context"
	"testing"
	"time"
)

func Test_Context(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	ctx2, _ := context.WithCancel(ctx)
	go func() {
	FOR:
		for {
			select {
			case <-ctx.Done():
				t.Log("done")
				break FOR
			default:
				t.Log("cancel1")
				cancel()
			}
		}
	}()
	go func() {
	FOR2:
		for {
			select {
			case <-ctx2.Done():
				t.Log("2done")
				break FOR2
			default:
				t.Log("cancel2")
				cancel()
			}
		}
	}()
	time.Sleep(time.Second)
	go cancel()
	go cancel()
	time.Sleep(time.Second)
	t.Log("main done")
}
