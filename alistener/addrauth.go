package alistener

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/rikaaa0928/awsl/aconn"
	"github.com/rikaaa0928/awsl/consts"
	"github.com/rikaaa0928/awsl/utils"
)

func NewAddrAuthMid(conf map[string]interface{}) AcceptMid {
	return func(next Acceptor) Acceptor {
		return func(ctx context.Context) (context.Context, aconn.AConn, error) {
			ctx, conn, err := next(ctx)
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
			dataMap := map[string]interface{}{}
			err = json.Unmarshal(data, &dataMap)
			if err != nil {
				conn.Close()
				return ctx, nil, err
			}
			rAuth, ok = dataMap["auth"]
			if !ok {
				conn.Close()
				return ctx, nil, errors.New("no auth in map. map:" + fmt.Sprintf("%+v", dataMap))
			}
			if auth.(string) != rAuth.(string) {
				conn.Close()
				return ctx, nil, errors.New("auth failed")
			}

			addrStr, ok := dataMap["addr"]
			if !ok {
				conn.Close()
				return ctx, nil, errors.New("no addr in map:" + fmt.Sprintf("%+v", dataMap))
			}
			addr := aconn.AddrInfo{}
			err = json.Unmarshal([]byte(addrStr.(string)), &addr)
			if err != nil {
				conn.Close()
				return ctx, nil, err
			}
			conn.SetEndAddr(addr)
			return ctx, conn, nil
		}
	}
}
