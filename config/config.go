package config

import (
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/Evi1/awsl/model"
)

var conf *model.Object
var lock sync.Mutex

// Debug debug
var Debug bool

// UDP useudp
var UDP bool

// ConfigFile cf
var ConfigFile string

// RouteCache RouteCache
var RouteCache bool

// InitConf InitConf
func InitConf() {
	configFile := flag.String("c", "/etc/awsl/config.json", "path to config file")
	debug := flag.Bool("d", false, "debug")
	udp := flag.Bool("u", false, "udp")
	nostd := flag.Bool("nostd", false, "nostd")
	routeCache := flag.Bool("rc", false, "route cache")
	logFile := flag.String("l", "", "log file location")
	flag.Parse()
	Debug = *debug
	UDP = *udp
	RouteCache = *routeCache
	confBytes, err := ioutil.ReadFile(filepath.FromSlash(*configFile))
	if err != nil {
		panic(err)
	}
	conf = &model.Object{}
	err = json.Unmarshal(confBytes, conf)
	if err != nil {
		panic(err)
	}
	ConfigFile = *configFile
	if conf.BufSize == 0 {
		conf.BufSize = 32
	}
	if Debug {
		conf.NoVerify = true
	}
	// log
	var writerList []io.Writer
	if !*nostd {
		writerList = append(writerList, os.Stdout)
	}
	if logFile != nil && len(*logFile) != 0 {
		writer, err := os.OpenFile(*logFile, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}
		writerList = append(writerList, writer)
	}
	if len(writerList) == 0 {
		log.SetOutput(ioutil.Discard)
	} else {
		log.SetOutput(io.MultiWriter(writerList...))
	}

}

// GetConf GetConf
func GetConf() *model.Object {
	if conf == nil {
		lock.Lock()
		if conf == nil {
			InitConf()
		}
		lock.Unlock()
	}
	return conf
}
