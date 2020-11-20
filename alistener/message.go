package alistener

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/rikaaa0928/awsl/aconn"
	"github.com/rikaaa0928/awsl/consts"
	"github.com/rikaaa0928/awsl/utils"
	"github.com/rikaaa0928/awsl/utils/ctxdatamap"
)

func NewMessageMid(_ context.Context, inTag string, conf map[string]interface{}) AcceptMid {
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
			lenBytes := utils.GetMem(4)
			defer utils.PutMem(lenBytes)
			_, err = io.ReadFull(conn, lenBytes)
			if err != nil {
				conn.Close()
				return ctx, nil, err
			}
			length := binary.BigEndian.Uint32(lenBytes)
			buf := make([]byte, length)
			//buf := utils.GetMem(int(length))
			//defer utils.PutMem(buf)
			n, err := io.ReadFull(conn, buf)
			if err != nil {
				conn.Close()
				return ctx, nil, err
			}
			data := buf[:n]
			ctx = ctxdatamap.Parse(ctx, data)
			//fmt.Println(length, len(data), string(data), string(ctxdatamap.Bytes(ctx)))
			rAuth = ctxdatamap.Get(ctx, consts.TransferAuth)
			if rAuth == nil {
				conn.Close()
				return ctx, nil, errors.New("no auth in map. map:" + fmt.Sprintf("%+v", string(ctxdatamap.Bytes(ctx))) + "\nread message data:" + string(data))
			}
			if auth.(string) != rAuth.(string) {
				conn.Close()
				return ctx, nil, errors.New("auth failed")
			}

			addrIn := ctxdatamap.Get(ctx, consts.TransferAddr)
			if addrIn == nil {
				conn.Close()
				return ctx, nil, errors.New("no addr in map:" + fmt.Sprintf("%+v", string(ctxdatamap.Bytes(ctx))) + "\nread message data:" + string(data))
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
