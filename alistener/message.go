package alistener

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/rikaaa0928/awsl/aconn"
	"github.com/rikaaa0928/awsl/consts"
	"github.com/rikaaa0928/awsl/utils"
	"github.com/rikaaa0928/awsl/utils/ctxdatamap"
)

func NewMessageMid(ctx context.Context, inTag string, conf map[string]interface{}) AcceptMid {
	return func(next Acceptor) Acceptor {
		return func(ctx context.Context) (context.Context, aconn.AConn, error) {
			ctx, conn, err := next(ctx)
			if err != nil {
				return ctx, conn, err
			}
			ctx = context.WithValue(ctx, consts.CTXInTag, inTag)
			auth, ok := conf["auth"]
			if !ok {
				conn.Close()
				return ctx, nil, errors.New("no auth in conf. map:" + fmt.Sprintf("%+v", conf))
			}
			rAuth := ctx.Value(consts.CTXReceiveAuth)
			if conn.EndAddr() != nil && rAuth != nil {
				if auth.(string) != rAuth.(string) {
					conn.Close()
					return ctx, nil, errors.New("auth failed")
				} else {
					return ctx, conn, nil
				}
			}
			if err != nil {
				return ctx, nil, err
			}
			buf := utils.GetMem(65536)
			defer utils.PutMem(buf)
			n, err := conn.Read(buf)
			data := buf[:n]
			if err != nil {
				conn.Close()
				return ctx, nil, err
			}
			ctx = ctxdatamap.Parse(ctx, data)
			// dataMap := map[string]interface{}{}
			// err = json.Unmarshal(data, &dataMap)
			// if err != nil {
			// 	conn.Close()
			// 	return ctx, nil, err
			// }
			// rAuth, ok = dataMap["auth"]
			// if !ok {
			// 	conn.Close()
			// 	return ctx, nil, errors.New("no auth in map. map:" + fmt.Sprintf("%+v", dataMap))
			// }
			rAuth = ctxdatamap.Get(ctx, consts.TransferAuth)
			if rAuth == nil {
				conn.Close()
				return ctx, nil, errors.New("no auth in map. map:" + fmt.Sprintf("%+v", string(ctxdatamap.Bytes(ctx))))
			}
			if auth.(string) != rAuth.(string) {
				conn.Close()
				return ctx, nil, errors.New("auth failed")
			}

			//addrStr, ok := dataMap["addr"]
			addrIn := ctxdatamap.Get(ctx, consts.TransferAddr)
			if addrIn == nil {
				conn.Close()
				return ctx, nil, errors.New("no addr in map:" + fmt.Sprintf("%+v", string(ctxdatamap.Bytes(ctx))))
			}
			addr := aconn.AddrInfo{}
			err = json.Unmarshal([]byte(addrIn.(string)), &addr)
			if err != nil {
				conn.Close()
				return ctx, nil, err
			}
			conn.SetEndAddr(addr)
			return ctx, conn, nil
		}
	}
}
