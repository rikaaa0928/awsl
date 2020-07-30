package main

import (
	"github.com/rikaaa0928/awsl/config"
	"github.com/rikaaa0928/awsl/manage"
	"github.com/rikaaa0928/awsl/object"
)

func main() {
	o := object.NewObject(*config.GetConf())
	go manage.Manage(o)
	o.Run()
}
