package tools

// NewCounter NewCounter
func NewCounter() *Counter {
	c := &Counter{sum: make(chan uint64, 1)}
	c.sum <- 0
	return c
}

// Counter Counter
type Counter struct {
	sum chan uint64
}

// Add Add
func (c *Counter) Add(a uint64) {
	s := <-c.sum
	c.sum <- s + a
}

// Get Get
func (c *Counter) Get() uint64 {
	s := <-c.sum
	c.sum <- s
	return s
}

// GetSet GetSet
func (c *Counter) GetSet(a uint64) uint64 {
	s := <-c.sum
	c.sum <- a
	return s
}
