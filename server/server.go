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
	case "h2c", "awsl", "http":
		return NewHTTPServer(typ, conf["host"].(string), conf["uri"].(string), conf["cert"].(string), conf["key"].(string), int(conf["port"].(float64))), nil
	default:
	}
	return nil, fmt.Errorf("error type: %v", typ)
}

type serveListener interface {
	alistener.AListener
	h(w http.ResponseWriter, r *http.Request)
	setSrv(*http.Server)
	srv() *http.Server
}
