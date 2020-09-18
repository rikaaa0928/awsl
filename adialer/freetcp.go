package adialer

import (
	"context"
	"net"

	"github.com/rikaaa0928/awsl/aconn"
)

var FreeTCP = func(_ context.Context, addr net.Addr) (aconn.AConn, error) {
	c, err := net.Dial("tcp", addr.String())
	ac := aconn.NewAConn(c)
	ac.SetEndAddr(addr)
	return ac, err
}
