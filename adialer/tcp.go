package adialer

import (
	"context"
	"fmt"
	"github.com/rikaaa0928/awsl/aconn"
	"github.com/rikaaa0928/awsl/global"
	"github.com/rikaaa0928/awsl/utils/ctxdatamap"
	"net"
	"strconv"
)

func NewTCP(conf map[string]interface{}) ADialer {
	remoteHost := conf["host"].(string)
	remotePort := strconv.Itoa(int(conf["port"].(float64)))
	auth := conf["auth"].(string)
	return func(ctx context.Context, addr net.Addr) (context.Context, aconn.AConn, error) {
		ctx = context.WithValue(ctx, global.CTXOutType, "tcp")
		ctx = ctxdatamap.Set(ctx, global.TransferAuth, auth)
		c, err := net.Dial("tcp", fmt.Sprintf("%s:%s", remoteHost, remotePort))
		ac := aconn.NewAConn(c)
		ac.SetEndAddr(addr)
		return ctx, ac, err
	}
}
