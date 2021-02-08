package arouter

import (
	"context"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"

	"github.com/rikaaa0928/awsl/config"
	"github.com/rikaaa0928/awsl/global"
	"github.com/rikaaa0928/awsl/utils/rulelist"
)

type ARouter func(context.Context, net.Addr) context.Context

var NopRouter = func(ctx context.Context, _ net.Addr) context.Context {
	return context.WithValue(ctx, global.CTXOutTag, "default")
}

type rule struct {
	list []rulelist.Rule
	tag  string
}
type route struct {
	rules []rule
	tag   string
}

func NewRouter(conf config.Configs) ARouter {
	datas := make(map[string][]rulelist.Rule)
	data, err := conf.GetStrMap("data")
	if err != nil {
		data = make(map[string]string)
	}
	for k, v := range data {
		f, err := os.Open(v)
		if err != nil {
			log.Println(err)
			continue
		}
		fData, err := ioutil.ReadAll(f)
		if err != nil {
			log.Println(err)
			continue
		}
		fDataStr := strings.Replace(string(fData), "\r\n", "\n", -1)
		l := strings.Split(fDataStr, "\n")
		datas[k] = make([]rulelist.Rule, 0, len(l))
		for _, line := range l {
			datas[k] = append(datas[k], rulelist.New(line))
		}
		f.Close()
	}
	routeMap := make(map[string]route)
	c, err := conf.GetMap("router", "default")
	if err != nil {
		panic(err)
	}
	rules := rule{}
	r := c["rules"].(map[string]interface{})
	rules.list = datas[r["data"].(string)]
	rules.tag = r["tag"].(string)
	routeMap["default"] = route{
		rules: nil,
		tag:   c["tag"].(string),
	}
	return func(ctx context.Context, addr net.Addr) context.Context {
		inTag := ctx.Value(global.CTXInTag).(string)
		router, ok := routeMap[inTag]
		if !ok {
			c, err := conf.GetMap("router", inTag)
			if err != nil {

			}
		}
	}
}
