package config_test

import (
	"encoding/json"
	"testing"

	"github.com/rikaaa0928/awsl/config"
)

func TestJson(t *testing.T) {
	var c config.Configs
	c = config.NewJsonConfig()

	err := c.Open("../test/conf.json")
	if err != nil {
		t.Fatal(err)
	}

	typ, err := c.GetString("ins", "socks5", "type")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(typ)

	switch typ {
	case "socks":
		t.Log(typ)
	default:
		t.Fatal(typ)
	}

	s, err := c.Get("ins")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(s.(json.RawMessage)))
	s, err = c.Get("ins", "socks5")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(s.(json.RawMessage)))
	s, err = c.Get("ins", "socks5", "type")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(s.(json.RawMessage)))

	port, err := c.GetString("ins", "socks5", "port")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(port)

	server0, err := c.GetMap("ins", "socks5")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(server0)
	t.Log(int(server0["port"].(float64)))
}
