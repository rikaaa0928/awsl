package router

import (
	"encoding/json"
	"testing"

	"github.com/rikaaa0928/awsl/model"
)

func Test_R(t *testing.T) {
	conf := model.Object{}
	confStr := `{
		"outs": [
			{
				"tag":"out1"
			},
			{
				"tag":"out2"
			}
		],
		"ins": [
			{
				"tag":"in1"
			},
			{
				"tag":"in2"
			}
		],
		"data": {
			"cn":{
				"name":"/home/kiririn/src/awsl/test/ip.txt",
				"type":0
			}
		},
		"routerules":[
			{
				"intags":["in1"],
				"outtags":["out2"],
				"datatags":["cn"]
			},
			{
				"intags":["in2","in1"],
				"outtags":["out1","out2"]
			}
		]
	}`
	err := json.Unmarshal([]byte(confStr), &conf)
	if err != nil {
		t.Error(err)
	}
	//t.Log(confStr,conf)
	r := NewDefaultRouter(conf)
	//t.Log(r.Resolver.Resolve("huya.com"))
	t.Log(r.Resolver.Resolve("api.steampowered.com"))
	//t.Log(r.Route(0, model.ANetAddr{Host: "huya.com", Typ: 1}))
	//t.Log(r.Route(0, model.ANetAddr{Host: "api.steampowered.com", Typ: 1}))
	t.Log(r.Route(0, model.ANetAddr{Host: "api.steampowered.com", Typ: 1}))
}
