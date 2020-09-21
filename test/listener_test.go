package test

import (
	"net"
	"testing"
)

func TestClose(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:12345")
	if err != nil {
		t.Fatal(err)
	}
	l.Close()
	t.Log(l.Accept())
}
