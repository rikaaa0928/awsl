package config

import "flag"

import "io/ioutil"

import "encoding/json"

import "github.com/Evi1/awsl/model"

// Conf conf
var Conf model.Object

// Debug debug
var Debug bool

func init() {
	configFile := flag.String("c", "/etc/awsl/config.json", "path to config file")
	debug := flag.Bool("d", false, "debug")
	flag.Parse()
	Debug = *debug
	confBytes, err := ioutil.ReadFile(*configFile)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(confBytes, &Conf)
	if err != nil {
		panic(err)
	}
}
