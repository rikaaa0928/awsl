package alistener

import (
	"context"

	"github.com/rikaaa0928/awsl/aconn"
)

type contextString string

const (
	CTXIntag contextString = "intag"
)

type AListener interface {
	Accept(context.Context) (context.Context, aconn.AConn, error)
	Close() error
}

type Accepter func(context.Context) (context.Context, aconn.AConn, error)

type AcceptMid func(Accepter) Accepter
