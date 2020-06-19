package main

import (
	"github.com/Evi1/awsl/config"
	"github.com/Evi1/awsl/manage"
	"github.com/Evi1/awsl/object"
)

func main() {
	go manage.Manage()
	o := object.NewObject(*config.GetConf())
	o.Run()
}
