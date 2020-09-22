package adialer

import (
	"context"

	"github.com/rikaaa0928/awsl/consts"
)

var TestFactory = func(_ context.Context, _ ...[]byte) ADialer {
	return FreeTCP
}

func NewFactory(conf map[string]interface{}) DialerFactory {
	return func(ctx context.Context, _ ...[]byte) ADialer {
		outTag := ctx.Value(consts.CTXRoute)
		if outTag == nil {
			return nil
		}
		tag, ok := outTag.(string)
		if !ok {
			return nil
		}
		tagConf := conf[tag].(map[string]interface{})
		switch tagConf["type"].(string) {
		case "h2c":
			return NewH2C(tagConf)
		case "free":
			return FreeTCP
		default:
		}
		return nil
	}
}
