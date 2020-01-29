package test

import (
	"net"
	"testing"

	"github.com/Evi1/awsl/tools"
)

func TestSize(t *testing.T) {
	l, _ := net.Listen("tcp", ":1996")
	c, _ := l.Accept()
	buf := tools.MemPool.Get(65536)

	t.Log(len(buf), cap(buf), buf)
	n, _ := c.Read(buf)
	t.Log(n, len(buf), cap(buf), buf[:n])
	t.Log(n, len(buf), cap(buf), buf[:n])
	t.Log(len(buf), cap(buf), buf)
	c.Write([]byte("\x05\x00"))

	t.Log(len(buf), cap(buf), buf)
	n, _ = c.Read(buf)
	t.Log(n, len(buf), cap(buf), buf[:n])

	tools.MemPool.Put(buf)
	buf2 := tools.MemPool.Get(65536)

	t.Log(len(buf2), cap(buf2), buf2)
}
