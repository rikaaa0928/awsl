package arouter

import (
	"context"
	"net"

	"github.com/rikaaa0928/awsl/consts"
)

type ARouter func(context.Context, net.Addr) context.Context

var NopRouter = func(ctx context.Context, _ net.Addr) context.Context {
	return context.WithValue(ctx, consts.CTXRoute, "default")
}
