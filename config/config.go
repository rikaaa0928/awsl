package config

import "flag"

import "io/ioutil"

import "encoding/json"

import "github.com/Evi1/awsl/model"

import "path/filepath"

// Conf conf
var Conf model.Object

// Debug debug
var Debug bool

func init() {
	configFile := flag.String("c", "/etc/awsl/config.json", "path to config file")
	debug := flag.Bool("d", false, "debug")
	flag.Parse()
	Debug = *debug
	confBytes, err := ioutil.ReadFile(filepath.FromSlash(*configFile))
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(confBytes, &Conf)
	if err != nil {
		panic(err)
	}
	if Conf.BufSize == 0 {
		Conf.BufSize = 32
	}
	if Debug {
		Conf.NoVerify = true
	}
}
