package bitmap

import (
	"fmt"
	"testing"
)

func TestBit(t *testing.T) {
	b := NewBitMap(8)
	b.Set(3)
	t.Log(fmt.Sprintf("%08b", b.m))
	b.Set(7)
	t.Log(fmt.Sprintf("%08b", b.m))
	b.Set(6)
	t.Log(fmt.Sprintf("%08b", b.m))
	b.Del(6)
	t.Log(fmt.Sprintf("%08b", b.m))
	b.Del(5)
	t.Log(fmt.Sprintf("%08b", b.m))
	b.Del(7)
	t.Log(fmt.Sprintf("%08b", b.m))
	b.Set(8)
	t.Log(fmt.Sprintf("%08b", b.m))
	b.Del(3)
	t.Log(fmt.Sprintf("%08b", b.m))
	b.Set(15)
	t.Log(fmt.Sprintf("%08b", b.m))
	b.Set(16)
	t.Log(fmt.Sprintf("%08b", b.m))
	t.Log(uint64(0xffffffffffffffff))
	t.Log(fmt.Sprintf("%b", uint64(0xffffffffffffffff)))
	t.Log(len(fmt.Sprintf("%b", uint64(0xffffffffffffffff))))
	t.Log(uint32(0xffffffff)/64 + 1)
	b2 := NewBitMap(0xffffffff)
	t.Log(len(b2.m))
	//t.Log(len(fmt.Sprintf("%b", b2.m)))
}
