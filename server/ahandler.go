package server

import (
	"context"

	"github.com/rikaaa0928/awsl/aconn"
	"github.com/rikaaa0928/awsl/adialer"
	"github.com/rikaaa0928/awsl/arouter"
)

type AHandler func(context.Context, aconn.AConn, arouter.ARouter, adialer.DialerFactory)
