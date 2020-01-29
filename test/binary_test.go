package test

import (
	"bytes"
	"encoding/binary"
	"net"
	"testing"
)

func TestBinary(t *testing.T) {
	buf := new(bytes.Buffer)
	var num uint16 = 65535
	err := binary.Write(buf, binary.LittleEndian, num)
	if err != nil {
		t.Error("binary.Write failed:", err)
	}
	t.Logf("% x", len(buf.Bytes()))

	b := []byte{ 0x0, 0x01}
	tb := bytes.NewBuffer(b)
	var y int16
	binary.Read(tb, binary.BigEndian, &y)
	t.Log(y)
}

func TestIP(t *testing.T) {
	a := []byte{0, 0, 0, 0, 0, 0}
	ip := net.ParseIP("127.0.0.1").To4()
	t.Logf("% x,%d", []byte(ip), len([]byte(ip)))
	copy(a[1:4], []byte(ip))
	t.Logf("% x", a)
}

type a struct {
	a int
}

func (o a) fa() int {
	return o.a
}

type b struct {
	a
	b int
}

type ia interface {
	fa() int
}

func TestStruct(t *testing.T) {
	var o ia
	o = b{}
	t.Log(o.(b))
}
