package adialer

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"net"

	"github.com/rikaaa0928/awsl/aconn"
	"github.com/rikaaa0928/awsl/global"
	"github.com/rikaaa0928/awsl/utils"
	"github.com/rikaaa0928/awsl/utils/ctxdatamap"
)

func NewAddrDataMid(next ADialer) ADialer {
	return func(ctx context.Context, addr net.Addr) (context.Context, aconn.AConn, error) {
		ctx, conn, err := next(ctx, addr)
		if err != nil {
			return ctx, nil, err
		}
		// data := ctx.Value(global.CTXSendData)
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
		ctx = ctxdatamap.Set(ctx, global.TransferAddr, string(addrBytes))
		// dataMap["addr"] = string(addrBytes)
		// dataBytes, err := json.Marshal(dataMap)
		// if err != nil {
		// 	conn.Close()
		// 	return ctx, nil, err
		// }
		// ctx = context.WithValue(ctx, global.CTXSendData, string(dataBytes))
		return ctx, conn, nil
	}
}

func NewSendDataMid(next ADialer) ADialer {
	return func(ctx context.Context, addr net.Addr) (context.Context, aconn.AConn, error) {
		ctx, conn, err := next(ctx, addr)
		if err != nil {
			return ctx, nil, err
		}
		//data := ctx.Value(global.CTXSendData)
		data := ctxdatamap.Bytes(ctx)
		length := len(data)
		if length == 0 {
			return ctx, conn, nil
		}
		lenBytes := utils.GetMem(4)
		defer utils.PutMem(lenBytes)
		binary.BigEndian.PutUint32(lenBytes, uint32(length))
		_, err = conn.Write(append(lenBytes, data...))
		if err != nil {
			conn.Close()
			return ctx, nil, err
		}
		//time.Sleep(500 * time.Millisecond)
		//fmt.Println("client write data done ", ctx.Value(ctxdatamap.CTXMapData))
		return ctx, conn, nil
	}
}
