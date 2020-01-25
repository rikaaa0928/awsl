package main

import (
	"github.com/Evi1/awsl/config"
	"github.com/Evi1/awsl/object"
)

func main() {
	o := object.NewObject(config.Conf)
	o.Run()
}
