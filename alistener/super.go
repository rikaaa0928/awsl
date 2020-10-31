package alistener

import (
	"context"
	"errors"

	"github.com/rikaaa0928/awsl/aconn"
	"github.com/rikaaa0928/awsl/consts"
	"github.com/rikaaa0928/awsl/utils/superlib"
)

var ErrorSupper = errors.New("supper connection")

func NewSupperMid(ctx context.Context, inTag string, conf map[string]interface{}) AcceptMid {
	return func(next Acceptor) Acceptor {
		return func(ctx context.Context) (context.Context, aconn.AConn, error) {
			ctx, conn, err := next(ctx)
			if err != nil {
				return ctx, conn, nil
			}
			if ctx.Value(consts.TransferSupper) == nil {
				return ctx, conn, nil
			}
			go func() {
				superlib.NewID()
			}()
			return ctx, nil, ErrorSupper
		}
	}
}
