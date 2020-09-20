package adialer

import "context"

var TestFactory = func(_ context.Context, _ ...[]byte) ADialer {
	return FreeTCP
}
