package main

import "github.com/Evi1/awsl/object"

import "github.com/Evi1/awsl/servers"

import "github.com/Evi1/awsl/clients"

func main() {
	o := object.DefaultObject{S: servers.Socke5Server{IP: "0.0.0.0", Port: "48888"},
		C:     clients.DirectOut{},
		Msg:   make(chan object.DefaultRemoteMsg, 10),
		Close: make(chan int8)}
	o.Run()
}
