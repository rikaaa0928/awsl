package alistener

import (
	"context"

	"github.com/rikaaa0928/awsl/consts"
)

func DefaultAcceptMids(ctx context.Context, l AcceptMidor, ty, tag string, conf map[string]interface{}) {
	switch ty {
	case "socks", "socks5", "socks4":
		l.RegisterAcceptor(NewSocksAcceptMid(ctx, tag, conf))
	case "awsl", "tcp", "h2c":
		l.RegisterAcceptor(NewMessageMid(ctx, tag, conf))
	default:
	}
	ctx = context.WithValue(ctx, consts.CTXInTag, tag)
}
