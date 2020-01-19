package test

import (
	"github.com/Evi1/awsl/clients"
	"github.com/Evi1/awsl/object"
	"github.com/Evi1/awsl/router"
	"github.com/Evi1/awsl/servers"
	"testing"
)

func TestDefault(t *testing.T) {
	s := []servers.Server{servers.Socke5Server{IP: "0.0.0.0", Port: "48888"}, servers.Socke5Server{IP: "0.0.0.0", Port: "58888"}}
	c := []clients.Client{clients.DirectOut{}}
	m := make([]chan object.DefaultRemoteMsg, len(c))
	for i := range m {
		m[i] = make(chan object.DefaultRemoteMsg, 10)
	}
	o := object.DefaultObject{S: s,
		C:     c,
		Msg:   m,
		Close: make(chan int8),
		R:     router.ARouter{}}
	o.Run()
}
