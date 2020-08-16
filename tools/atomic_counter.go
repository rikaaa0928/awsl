package tools

import "sync/atomic"

// AtomicCounter AtomicCounter
type AtomicCounter struct {
	sum int64
}

// Add Add
func (c *AtomicCounter) Add(a int64) {
	atomic.AddInt64(&c.sum, a)
}

// Get Get
func (c *AtomicCounter) Get() int64 {
	s := atomic.LoadInt64(&c.sum)
	return s
}

// Set Set
func (c *AtomicCounter) Set(a int64) int64 {
	s := atomic.LoadInt64(&c.sum)
	atomic.StoreInt64(&c.sum, a)
	return s
}
