package router

import (
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/Evi1/awsl/model"
	"github.com/Evi1/awsl/tools/dns"
	inlist "github.com/Evi1/awsl/tools/inList"
)

// NewDefaultRouter NewDefaultRouter
func NewDefaultRouter(conf model.Object) ARouter {
	inMap := make(map[int]string)
	for i, v := range conf.Ins {
		if len(v.Tag) == 0 {
			continue
		}
		inMap[i] = v.Tag
	}
	outMap := make(map[string]int)
	for i, v := range conf.Outs {
		if len(v.Tag) == 0 {
			continue
		}
		outMap[v.Tag] = i
	}
	if conf.Data == nil {
		return ARouter{RuleSet: nil, RulesForIn: nil, InMap: inMap, OutMap: outMap,
			Resolver: dns.DoH{URL: "https://cloudflare-dns.com/dns-query"},
			Cache:    make(map[string]int)}
	}
	ruleSet := make(map[string]inlist.InList)
	for k, v := range conf.Data {
		ruleStr, err := ioutil.ReadFile(v.Name)
		if err != nil {
			panic(err)
		}
		ruleList := strings.Split(strings.Replace(string(ruleStr), "\r\n", "\n", -1), "\n")
		if v.Type == 1 {
			continue
		}
		ruleSet[k] = inlist.NewIPList(ruleList)
	}
	ruleForIn := make(map[string][]routeRule)
	for _, v := range conf.RouteRules {
		rr := routeRule{}
		rr.OutTag = v.OutTag
		for _, vv := range v.DataTags {
			rr.RuleTag = vv
			for _, vvv := range v.InTags {
				_, ok := ruleForIn[vvv]
				if !ok {
					ruleForIn[vvv] = make([]routeRule, 0)
				}
				ruleForIn[vvv] = append(ruleForIn[vvv], rr)
			}
		}
	}
	return ARouter{RuleSet: ruleSet, RulesForIn: ruleForIn, InMap: inMap, OutMap: outMap,
		Resolver: dns.DoH{URL: "https://cloudflare-dns.com/dns-query"},
		Cache:    make(map[string]int)}
}

type routeRule struct {
	RuleTag string
	OutTag  string
}

// ARouter ARouter
type ARouter struct {
	RuleSet    map[string]inlist.InList
	RulesForIn map[string][]routeRule
	InMap      map[int]string
	OutMap     map[string]int
	Resolver   dns.DNS
	Cache      map[string]int
}

// Route Route
func (r ARouter) Route(src int, addr model.ANetAddr) int {
	if r.RuleSet == nil || r.RulesForIn == nil {
		return 0
	}
	// cache
	result, ok := r.Cache[strconv.Itoa(src)+addr.Host]
	if ok {
		return result
	}
	// resolve
	inTag, ok := r.InMap[src]
	if !ok {
		return 0
	}
	rules, ok := r.RulesForIn[inTag]
	if !ok {
		return 0
	}
	for _, v := range rules {
		ruleList, ok := r.RuleSet[v.RuleTag]
		if !ok {
			return 0
		}
		host := addr.Host
		if addr.Typ == model.RAWADDR {
			result, _ := r.Resolver.Resolve(addr.Host)
			if len(result.V4)+len(result.V6) == 0 {
				return 0
			}
			host = result.V4
			if len(host) == 0 {
				host = result.V6
			}
		}
		if ruleList.Include(host) {
			outID, ok := r.OutMap[v.OutTag]
			if !ok {
				return 0
			}
			r.Cache[strconv.Itoa(src)+addr.Host] = outID
			return outID
		}
	}
	r.Cache[strconv.Itoa(src)+addr.Host] = 0
	return 0
}
