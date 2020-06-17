package tools

import "context"

// NewCloseWait NewCloseWait
func NewCloseWait(ctx context.Context) *CloseWait {
	nctx, cancel := context.WithCancel(ctx)
	return &CloseWait{Ctx: nctx, cancel: cancel}
}

// CloseWait CloseWait
type CloseWait struct {
	Ctx    context.Context
	cancel context.CancelFunc
}

// Close Close
func (c *CloseWait) Close() {
	c.cancel()
}

// WaitClose Wait for Close
func (c *CloseWait) WaitClose() <-chan struct{} {
	return c.Ctx.Done()
}
