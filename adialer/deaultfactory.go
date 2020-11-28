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
		outTag := ctx.Value(global.CTXRoute)
		if outTag == nil {
			return nil
		}
		tag, ok := outTag.(string)
		if !ok {
			return nil
		}
		tagConf := conf[tag].(map[string]interface{})
		var d ADialer
		//superTyp := ctx.Value(global.CTXSuperType)
		//if superTyp != nil {
		//	superData := ctx.Value(global.CTXSuperData).(string)
		//	var udpMsg superlib.UDPMSG
		//	err := json.Unmarshal([]byte(superData), &udpMsg)
		//	if err != nil {
		//		log.Println(err)
		//		return d
		//	}
		//	inTag := ctx.Value(global.CTXInTag)
		//	if inTag == nil {
		//		log.Println("nil inTag")
		//		return d
		//	}
		//	d = getSuperConn(tag, inTag.(string)+":"+superlib.GetID(ctx), udpMsg.SrcStr, udpMsg.DstStr, tagConf)
		//	return d
		//}
		typ := tagConf["type"].(string)
		switch typ {
		case "free":
			d = Free
		case "awsl":
			d = NewAWSL(tagConf)
		case "quic":
			d = NewQUIC(tagConf)
		default:
		}
		d = DefaultDialMids(d, typ)
		return d
	}
}
