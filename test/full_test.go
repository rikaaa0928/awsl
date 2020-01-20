package test

import (
	"github.com/Evi1/awsl/clients"
	"github.com/Evi1/awsl/object"
	"github.com/Evi1/awsl/servers"
	"testing"
)

func TestFull(t *testing.T) {
	o1 := object.NewDefault([]clients.Client{clients.NewAWSL("127.0.0.1", "1928", "wss")}, []servers.Server{servers.NewSocks5("127.0.0.1", "48888")})
	s := servers.NewAWSL("127.0.0.1", "1928", "wss", GetTestPath()+"/server.key", GetTestPath()+"/server.crt")
	o2 := object.NewDefault([]clients.Client{clients.DirectOut{}}, []servers.Server{s})
	go o2.Run()
	o1.Run()
}
