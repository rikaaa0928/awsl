package alistener

import (
	"context"
)

func DefaultAcceptMids(ctx context.Context, l AcceptMidor, typ, tag string, conf map[string]interface{}) {
	switch typ {
	case "socks", "socks5", "socks4":
		l.RegisterAcceptor(NewSocksAcceptMid(ctx, tag, conf))
	case "awsl", "tcp":
		l.RegisterAcceptor(NewMessageMid(ctx, tag, conf))
	default:
	}
}
