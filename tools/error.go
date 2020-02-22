package tools

var ErrTimeout error = &TimeoutError{}

type TimeoutError struct{}

// Implement the net.Error interface.
func (e *TimeoutError) Error() string   { return "i/o timeout" }
func (e *TimeoutError) Timeout() bool   { return true }
func (e *TimeoutError) Temporary() bool { return true }
