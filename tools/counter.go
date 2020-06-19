package tools

// NewCounter NewCounter
func NewCounter() *Counter {
	c := &Counter{sum: make(chan int64, 1)}
	c.sum <- 0
	return c
}

// Counter Counter
type Counter struct {
	sum chan int64
}

// Add Add
func (c *Counter) Add(a int64) {
	s := <-c.sum
	c.sum <- s + a
}

// Get Get
func (c *Counter) Get() int64 {
	s := <-c.sum
	c.sum <- s
	return s
}

// GetSet GetSet
func (c *Counter) GetSet(a int64) int64 {
	s := <-c.sum
	c.sum <- a
	return s
}
