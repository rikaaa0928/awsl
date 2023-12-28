package server

import (
	"fmt"
	"net/http"

	"github.com/rikaaa0928/awsl/alistener"
)

type AServer interface {
	Listen() alistener.AListener
	Handler() AHandler
}

func NewServer(typ string, conf map[string]interface{}) (AServer, error) {
	switch typ {
	case "socks", "socks5", "socks4", "tcp":
		return NewBaseTcp(conf["host"].(string), int(conf["port"].(float64))), nil
	case "awsl", "http", "pprof":
		uri, ok := conf["uri"]
		if !ok {
			uri = ""
		}
		cert, ok := conf["cert"]
		if !ok {
			cert = ""
		}
		key, ok := conf["key"]
		if !ok {
			key = ""
		}
		return NewHTTPServer(typ, conf["host"].(string), uri.(string), cert.(string), key.(string),
			int(conf["port"].(float64))), nil
	//case "quic":
	//	cert, ok := conf["cert"]
	//	if !ok {
	//		cert = ""
	//	}
	//	key, ok := conf["key"]
	//	if !ok {
	//		key = ""
	//	}
	//	return NewBaseQUIC(conf["host"].(string), int(conf["port"].(float64)), cert.(string), key.(string)), nil
	default:
	}
	return nil, fmt.Errorf("error type: %v", typ)
}

type serveListener interface {
	alistener.AListener
	h(w http.ResponseWriter, r *http.Request)
	setSrv(*http.Server)
	srv() *http.Server
	handler() AHandler
}
