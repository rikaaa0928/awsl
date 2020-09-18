package adialer

import "context"

var DefaultFactory = func(_ context.Context) ADialer {
	return FreeTCP
}
