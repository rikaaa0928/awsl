package arouter

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"sync"

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
	data, err := conf.GetStrMap("router", "data")
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
	lock := sync.RWMutex{}
	c, err := conf.GetMap("router", "default")
	if err != nil {
		panic(err)
	}
	ruleL := make([]rule, 0)
	rli, ok := c["rules"]
	if !ok {
		rli = make([]interface{}, 0)
	}
	rl := rli.([]interface{})
	for _, ri := range rl {
		r := ri.(map[string]interface{})
		rules := rule{}
		rules.list = datas[r["data"].(string)]
		rules.tag = r["tag"].(string)
		ruleL = append(ruleL, rules)
	}
	lock.Lock()
	routeMap["default"] = route{
		rules: ruleL,
		tag:   c["tag"].(string),
	}
	lock.Unlock()
	return func(ctx context.Context, addr net.Addr) context.Context {
		inTag := ctx.Value(global.CTXInTag).(string)
		lock.RLock()
		router, ok := routeMap[inTag]
		lock.RUnlock()
		if !ok {
			c, err := conf.GetMap("router", inTag)
			if err != nil || c == nil {
				lock.RLock()
				router = routeMap["default"]
				lock.RUnlock()
			} else {
				ruleL := make([]rule, 0)
				rli, ok := c["rules"]
				if !ok {
					rli = make([]interface{}, 0)
				}
				rl := rli.([]interface{})
				for _, ri := range rl {
					r := ri.(map[string]interface{})
					rules := rule{}
					rules.list = datas[r["data"].(string)]
					rules.tag = r["tag"].(string)
					ruleL = append(ruleL, rules)
				}
				lock.Lock()
				routeMap[inTag] = route{
					rules: ruleL,
					tag:   c["tag"].(string),
				}
				router = routeMap[inTag]
				lock.Unlock()
			}
		}
		host, _, err := net.SplitHostPort(addr.String())
		if err != nil {
			return context.WithValue(ctx, global.CTXOutTag, router.tag)
		}
		for _, l := range router.rules {
			for _, r := range l.list {
				if r.Include(host) {
					ctx = context.WithValue(ctx, global.CTXOutTag, l.tag)
					fmt.Println(addr.String(), l.tag)
					return ctx
				}
			}
		}
		ctx = context.WithValue(ctx, global.CTXOutTag, router.tag)
		fmt.Println(addr.String(), router.tag)
		return ctx
	}
}
