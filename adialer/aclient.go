package adialer

import (
	"context"
	"net"

	"github.com/rikaaa0928/awsl/aconn"
)

type ADialer func(ctx context.Context, addr net.Addr) (aconn.AConn, error)

type DialerFactory func(context.Context, ...[]byte) ADialer
