package main

import (
	"log"

	"github.com/rikaaa0928/awsl/config"
	"github.com/rikaaa0928/awsl/object"
)

func main() {
	conf := config.NewJsonConfig()
	err := conf.Open("./test/conf.json")
	if err != nil {
		panic(err)
	}
	ins, err := conf.GetMap("ins")
	log.Println(len(ins))
	for k := range ins {
		object.DefaultObject(k, conf)
	}
}
