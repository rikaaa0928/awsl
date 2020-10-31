package adialer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/rikaaa0928/awsl/consts"
	"github.com/rikaaa0928/awsl/utils/superlib"
)

var TestFactory = func(_ context.Context, _ ...[]byte) ADialer {
	return Free
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
			var udpMsg superlib.UDPMSG
			err := json.Unmarshal([]byte(superData), &udpMsg)
			if err != nil {
				log.Println(err)
				return d
			}
			inTag := ctx.Value(consts.CTXInTag)
			if inTag == nil {
				log.Println("nil inTag")
				return d
			}
			d = getSuperConn(tag, inTag.(string)+":"+superlib.GetID(ctx), udpMsg.SrcStr, udpMsg.DstStr, tagConf)
			return d
		}
		typ := tagConf["type"].(string)
		switch typ {
		case "h2c":
			d = NewH2C(tagConf)
		case "free":
			d = Free
		case "awsl":
			d = NewAWSL(tagConf)
		default:
		}
		d = DefaultDialMids(d, typ)
		return d
	}
}
