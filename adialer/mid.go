package adialer

import (
	"context"
	"encoding/json"
	"errors"
	"net"

	"github.com/rikaaa0928/awsl/aconn"
	"github.com/rikaaa0928/awsl/consts"
	"github.com/rikaaa0928/awsl/utils/ctxdatamap"
)

func NewAuthDataMid(next ADialer) ADialer {
	return func(ctx context.Context, addr net.Addr) (context.Context, aconn.AConn, error) {
		ctx, conn, err := next(ctx, addr)
		if err != nil {
			return ctx, nil, err
		}
		auth := ctx.Value(consts.CTXSendAuth)
		if auth == nil {
			conn.Close()
			return ctx, nil, errors.New("auth data mid: nil auth")
		}
		// data := ctx.Value(consts.CTXSendData)
		// var dataMap map[string]interface{}
		// if data == nil {
		// 	dataMap = make(map[string]interface{})
		// } else {
		// 	err = json.Unmarshal([]byte(data.(string)), &dataMap)
		// 	if err != nil {
		// 		conn.Close()
		// 		return ctx, nil, err
		// 	}
		// }
		// dataMap["auth"] = auth.(string)
		// dataBytes, err := json.Marshal(dataMap)
		// if err != nil {
		// 	conn.Close()
		// 	return ctx, nil, err
		// }
		// ctx = context.WithValue(ctx, consts.CTXSendData, string(dataBytes))
		ctx = ctxdatamap.Set(ctx, consts.TransferAuth, auth)
		return ctx, conn, nil
	}
}

func NewAddrDataMid(next ADialer) ADialer {
	return func(ctx context.Context, addr net.Addr) (context.Context, aconn.AConn, error) {
		ctx, conn, err := next(ctx, addr)
		if err != nil {
			return ctx, nil, err
		}
		// data := ctx.Value(consts.CTXSendData)
		// var dataMap map[string]interface{}
		// if data == nil {
		// 	dataMap = make(map[string]interface{})
		// } else {
		// 	err = json.Unmarshal([]byte(data.(string)), &dataMap)
		// 	if err != nil {
		// 		conn.Close()
		// 		return ctx, nil, err
		// 	}
		// }
		ai, ok := addr.(aconn.AddrInfo)
		if !ok {
			(&ai).Parse(addr.Network(), addr.String())
		}
		addrBytes, err := json.Marshal(ai)
		if err != nil {
			conn.Close()
			return ctx, nil, err
		}
		ctx = ctxdatamap.Set(ctx, consts.TransferAddr, string(addrBytes))
		// dataMap["addr"] = string(addrBytes)
		// dataBytes, err := json.Marshal(dataMap)
		// if err != nil {
		// 	conn.Close()
		// 	return ctx, nil, err
		// }
		// ctx = context.WithValue(ctx, consts.CTXSendData, string(dataBytes))
		return ctx, conn, nil
	}
}

func NewSendDataMid(next ADialer) ADialer {
	return func(ctx context.Context, addr net.Addr) (context.Context, aconn.AConn, error) {
		ctx, conn, err := next(ctx, addr)
		if err != nil {
			return ctx, nil, err
		}
		//data := ctx.Value(consts.CTXSendData)
		data := ctxdatamap.Bytes(ctx)
		if len(data) == 0 {
			return ctx, conn, nil
		}
		_, err = conn.Write(data)
		if err != nil {
			conn.Close()
			return ctx, nil, err
		}
		return ctx, conn, nil
	}
}
