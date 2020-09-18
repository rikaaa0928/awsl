package arouter

import (
	"context"
	"net"
)

type contextString string

const (
	CTXRoute contextString = "route"
)

type ARouter func(context.Context, net.Addr) context.Context

var NopRouter = func(ctx context.Context, _ net.Addr) context.Context {
	return context.WithValue(ctx, CTXRoute, "")
}
