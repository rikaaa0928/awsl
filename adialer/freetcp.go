package adialer

import (
	"context"
	"net"

	"github.com/rikaaa0928/awsl/aconn"
)

var FreeTCP = func(ctx context.Context, addr net.Addr) (context.Context, aconn.AConn, error) {
	c, err := net.Dial("tcp", addr.String())
	ac := aconn.NewAConn(c)
	ac.SetEndAddr(addr)
	return ctx, ac, err
}
