package adialer

import (
	"context"
	"encoding/json"
	"errors"
	"net"

	"github.com/rikaaa0928/awsl/aconn"
	"github.com/rikaaa0928/awsl/consts"
)

func NewAuthDataMid(next ADialer) ADialer {
	return func(ctx context.Context, addr net.Addr) (context.Context, aconn.AConn, error) {
		ctx, conn, err := next(ctx, addr)
		if err != nil {
			return ctx, nil, err
		}
		auth := ctx.Value(consts.CTXSendAuth)
		if auth == nil {
			return ctx, nil, errors.New("auth data mid: nil auth")
		}
		data := ctx.Value(consts.CTXSendData)
		var dataMap map[string]interface{}
		if data == nil {
			dataMap = make(map[string]interface{})
		} else {
			err = json.Unmarshal([]byte(data.(string)), &dataMap)
			if err != nil {
				return ctx, nil, err
			}
		}
		dataMap["auth"] = auth.(string)
		dataStr, err := json.Marshal(dataMap)
		if err != nil {
			return ctx, nil, err
		}
		ctx = context.WithValue(ctx, consts.CTXSendData, dataStr)
		return ctx, conn, nil
	}
}

func NewAddrDataMid(next ADialer) ADialer {
	return func(ctx context.Context, addr net.Addr) (context.Context, aconn.AConn, error) {
		ctx, conn, err := next(ctx, addr)
		if err != nil {
			return ctx, nil, err
		}
		data := ctx.Value(consts.CTXSendData)
		var dataMap map[string]interface{}
		if data == nil {
			dataMap = make(map[string]interface{})
		} else {
			err = json.Unmarshal([]byte(data.(string)), &dataMap)
			if err != nil {
				return ctx, nil, err
			}
		}
		ai, ok := addr.(aconn.AddrInfo)
		if !ok {
			(&ai).Parse(addr.Network(), addr.String())
		}
		dataMap["addr"], err = json.Marshal(ai)
		dataStr, err := json.Marshal(dataMap)
		if err != nil {
			return ctx, nil, err
		}
		ctx = context.WithValue(ctx, consts.CTXSendData, dataStr)
		return ctx, conn, nil
	}
}
