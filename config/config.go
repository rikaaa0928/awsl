package config

import "flag"

import "io/ioutil"

import "encoding/json"

import "github.com/Evi1/awsl/model"

// Conf conf
var Conf model.Object

func init() {
	configFile := flag.String("c", "config.json", "path to config file")
	flag.Parse()
	confBytes, err := ioutil.ReadFile(*configFile)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(confBytes, &Conf)
	if err != nil {
		panic(err)
	}
}
