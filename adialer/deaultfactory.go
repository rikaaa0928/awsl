package adialer

import (
	"context"

	"github.com/rikaaa0928/awsl/global"
)

var TestFactory = func(_ context.Context, _ ...[]byte) ADialer {
	return Free
}

func NewFactory(conf map[string]interface{}) DialerFactory {
	return func(ctx context.Context, _ ...[]byte) ADialer {
		outTag := ctx.Value(global.CTXOutTag)
		if outTag == nil {
			return nil
		}
		tag, ok := outTag.(string)
		if !ok {
			return nil
		}
		tagConf := conf[tag].(map[string]interface{})
		var d ADialer
		typ := tagConf["type"].(string)
		switch typ {
		case "free":
			d = Free
		case "awsl":
			d = NewAWSL(outTag.(string), tagConf)
		case "quic":
			d = NewQUIC(tagConf)
		default:
		}
		d = DefaultDialMids(d, typ)
		return d
	}
}
