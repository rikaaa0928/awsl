package tools

// NewCounter NewCounter
func NewCounter(t string) (c Counter) {
	if t == "chan" {
		c = &ChanCounter{sum: make(chan int64, 1)}
		c.(*ChanCounter).sum <- 0
	} else {
		c = &AtomicCounter{}
	}
	return
}

// Counter Counter Interface
type Counter interface {
	Set(a int64) int64
	Get() int64
	Add(a int64)
}
