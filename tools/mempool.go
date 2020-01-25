package tools

import "sync"

// MemPool MemPool
var MemPool *MPool

func init() {
	MemPool = NewMPool()
}

// MPool mem pool
type MPool struct {
	Pool map[int]chan []byte
	l    sync.Mutex
}

// NewMPool new
func NewMPool() *MPool {
	m := MPool{Pool: make(map[int]chan []byte), l: sync.Mutex{}}
	return &m
}

// Get Get
func (m *MPool) Get(size int) []byte {
	c, ok := m.Pool[size]
	if !ok {
		m.l.Lock()
		c, ok = m.Pool[size]
		if !ok {
			m.Pool[size] = make(chan []byte, 32)
			c = m.Pool[size]
		}
		m.l.Unlock()
	}
	var r []byte
	select {
	case r = <-c:
		return r
	default:
		return make([]byte, size)
	}
}

// Put Put
func (m *MPool) Put(bytes []byte) {
	size := cap(bytes)
	c, ok := m.Pool[size]
	if !ok {
		return
	}
	select {
	case c <- bytes:
	default:
		return
	}

}
