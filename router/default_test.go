package router

import (
	"encoding/json"
	"testing"

	"github.com/Evi1/awsl/model"
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
				"DataTags":["cn"]
			}
		]
	}`
	err := json.Unmarshal([]byte(confStr), &conf)
	if err != nil {
		t.Error(err)
	}
	//t.Log(confStr,conf)
	r := NewDefaultRouter(conf)
	t.Log(r.Resolver.Resolve("live.bilibili.com"))
	t.Log(r.Resolver.Resolve("pan.baidu.com"))
	t.Log(r.Route(0, model.ANetAddr{Host: "live.bilibili.com", Typ: 1}))
	t.Log(r.Route(0, model.ANetAddr{Host: "pan.baidu.com", Typ: 1}))
}
