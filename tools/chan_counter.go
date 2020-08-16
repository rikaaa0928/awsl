package tools

// ChanCounter ChanCounter
type ChanCounter struct {
	sum chan int64
}

// Add Add
func (c *ChanCounter) Add(a int64) {
	s := <-c.sum
	c.sum <- s + a
}

// Get Get
func (c *ChanCounter) Get() int64 {
	s := <-c.sum
	c.sum <- s
	return s
}

// Set Set
func (c *ChanCounter) Set(a int64) int64 {
	s := <-c.sum
	c.sum <- a
	return s
}
