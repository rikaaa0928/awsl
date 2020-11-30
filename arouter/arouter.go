package arouter

import (
	"context"
	"net"

	"github.com/rikaaa0928/awsl/global"
)

type ARouter func(context.Context, net.Addr) context.Context

var NopRouter = func(ctx context.Context, _ net.Addr) context.Context {
	return context.WithValue(ctx, global.CTXOutTag, "default")
}
