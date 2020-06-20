package main

import (
	"github.com/Evi1/awsl/config"
	"github.com/Evi1/awsl/manage"
	"github.com/Evi1/awsl/object"
)

func main() {
	o := object.NewObject(*config.GetConf())
	go manage.Manage(o)
	o.Run()
}
