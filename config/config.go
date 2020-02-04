package config

import "flag"

import "io/ioutil"

import "encoding/json"

import "github.com/Evi1/awsl/model"

import "path/filepath"

import "sync"

var conf *model.Object
var lock sync.Mutex

// Debug debug
var Debug bool

// UDP useudp
var UDP bool

func initConf() {
	configFile := flag.String("c", "/etc/awsl/config.json", "path to config file")
	debug := flag.Bool("d", false, "debug")
	udp := flag.Bool("u", false, "udp")
	flag.Parse()
	Debug = *debug
	UDP = *udp
	confBytes, err := ioutil.ReadFile(filepath.FromSlash(*configFile))
	if err != nil {
		panic(err)
	}
	conf = &model.Object{}
	err = json.Unmarshal(confBytes, conf)
	if err != nil {
		panic(err)
	}
	if conf.BufSize == 0 {
		conf.BufSize = 32
	}
	if Debug {
		conf.NoVerify = true
	}
}

// GetConf GetConf
func GetConf() model.Object {
	if conf == nil {
		lock.Lock()
		if conf == nil {
			initConf()
		}
		lock.Unlock()
	}
	return *conf
}
