package tools

import (
	"sync"
	"testing"
)

func BenchmarkMem(b *testing.B) {
	a := []byte{'a', 'b'}
	m := sync.WaitGroup{}
	for i := 0; i < 1000*1000; i++ {
		m.Add(1)
		go func() {
			c := make([]byte, 65535)
			copy(c, a)
			m.Done()
		}()
	}
	m.Wait()
}

func BenchmarkMemPool(b *testing.B) {
	a := []byte{'a', 'b'}
	p := NewMPool()
	m := sync.WaitGroup{}
	for i := 0; i < 1000*1000; i++ {
		m.Add(1)
		go func() {
			c := p.Get(655335)
			copy(c, a)
			p.Put(c)
			m.Done()
		}()
	}
	m.Wait()
}

func TestSlice(t *testing.T) {
	s := make([]int, 10)
	s[0] = 1
	t.Logf("%p,%+v\n", s, s)
	c := func(l []int) []int {
		return l[:1]
	}(s)
	t.Logf("%p,%+v\n", c, c)
}
