package adialer

import (
	"context"
	"encoding/json"
	"log"

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
		var d ADialer
		superTyp := ctx.Value(consts.CTXSuperType)
		if superTyp != nil {
			superData := ctx.Value(consts.CTXSuperData).(string)
			var udpMsg consts.UDPMSG
			err := json.Unmarshal([]byte(superData), &udpMsg)
			if err != nil {
				log.Println(err)
				return d
			}
			d = getSuperConn(tag, udpMsg.SrcStr, udpMsg.DstStr, tagConf)
			return d
		}
		typ := tagConf["type"].(string)
		switch typ {
		case "h2c":
			d = NewH2C(tagConf)
		case "free":
			d = FreeTCP
		case "awsl":
			d = NewAWSL(tagConf)
		default:
		}
		d = DefaultDialMids(d, typ)
		return d
	}
}
