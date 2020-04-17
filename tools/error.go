package tools

// ErrTimeout ErrTimeout
var ErrTimeout error = &TimeoutError{}

// TimeoutError TimeoutError
type TimeoutError struct{}

// Implement the net.Error interface.
// Error Error
func (e *TimeoutError) Error() string { return "i/o timeout" }

// Timeout Timeout
func (e *TimeoutError) Timeout() bool { return true }

// Temporary Temporary
func (e *TimeoutError) Temporary() bool { return true }
