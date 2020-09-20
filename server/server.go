package server

import (
	"fmt"

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
	default:
	}
	return nil, fmt.Errorf("error type: %v", typ)
}
