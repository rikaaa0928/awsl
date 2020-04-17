package router

import (
	"io/ioutil"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Evi1/awsl/model"
	"github.com/Evi1/awsl/tools/dns"
	inlist "github.com/Evi1/awsl/tools/inList"
)

// NewDefaultRouter NewDefaultRouter
func NewDefaultRouter(conf model.Object) *ARouter {

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
	r := &ARouter{RuleSet: nil, RulesForIn: nil, InMap: inMap, OutMap: outMap,
		Resolver: dns.DoH{URL: "https://cloudflare-dns.com/dns-query"},
		Cache:    make(map[string]routeCache),
		CLock:    sync.Mutex{}}
	if conf.Data == nil {
		return r
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
		rr.OutTags = v.OutTags
		if len(rr.OutTags) == 0 {
			continue
		}
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
	r.RuleSet = ruleSet
	r.RulesForIn = ruleForIn
	return r
}

type routeRule struct {
	RuleTag string
	OutTags []string
}

type routeCache struct {
	data     []int
	lastTime time.Time
	timeOut  time.Duration
}

// TempRoute TempRoute
type TempRoute interface {
	TempRoute(src int, addr model.ANetAddr, outID int)
}

// ARouter ARouter
type ARouter struct {
	RuleSet    map[string]inlist.InList
	RulesForIn map[string][]routeRule
	InMap      map[int]string
	OutMap     map[string]int
	Resolver   dns.DNS
	Cache      map[string]routeCache
	CLock      sync.Mutex
	Cleaning   bool
	LastClean  time.Time
}

// Route Route
func (r *ARouter) Route(src int, addr model.ANetAddr) []int {
	if r.RuleSet == nil || r.RulesForIn == nil {
		return []int{0}
	}
	// cache
	if !r.Cleaning && time.Now().After(r.LastClean.Add(6*time.Hour)) {
		go func() {
			r.CLock.Lock()
			defer r.CLock.Unlock()
			r.Cleaning = true
			defer func() {
				r.Cleaning = false
			}()
			if time.Now().Before(r.LastClean.Add(6 * time.Hour)) {
				return
			}
			for k, v := range r.Cache {
				if v.lastTime.Add(6 * time.Hour).Before(time.Now()) {
					delete(r.Cache, k)
				}
			}
			r.LastClean = time.Now()
		}()
	}
	r.CLock.Lock()
	result, ok := r.Cache[strconv.Itoa(src)+"-"+addr.Host]
	if ok {
		if result.lastTime.Add(result.timeOut).After(time.Now()) {
			r.CLock.Unlock()
			return result.data
		}
		delete(r.Cache, strconv.Itoa(src)+addr.Host)
	}
	r.CLock.Unlock()
	// resolve
	inTag, ok := r.InMap[src]
	if !ok {
		return []int{0}
	}
	rules, ok := r.RulesForIn[inTag]
	if !ok || len(rules) == 0 {
		return []int{0}
	}
	for _, v := range rules {
		ruleList, ok := r.RuleSet[v.RuleTag]
		if !ok || len(v.OutTags) == 0 {
			return []int{0}
		}
		host := addr.Host
		if addr.Typ == model.RAWADDR {
			result, _ := r.Resolver.Resolve(addr.Host)
			if len(result.V4)+len(result.V6) == 0 {
				return []int{0}
			}
			host = result.V4
			if len(host) == 0 {
				host = result.V6
			}
		}
		if ruleList.Include(host) {
			outIDs := make([]int, 0, len(v.OutTags))
			for _, outTag := range v.OutTags {
				outID, ok := r.OutMap[outTag]
				if !ok {
					return []int{0}
				}
				outIDs = append(outIDs, outID)
			}
			/*outID, ok := r.OutMap[v.OutTag]
			if !ok {
				return 0
			}*/
			r.CLock.Lock()
			defer r.CLock.Unlock()
			r.Cache[strconv.Itoa(src)+"-"+addr.Host] = routeCache{data: outIDs, lastTime: time.Now(), timeOut: 6 * time.Hour}
			return outIDs
		}
	}
	r.CLock.Lock()
	defer r.CLock.Unlock()
	r.Cache[strconv.Itoa(src)+"-"+addr.Host] = routeCache{data: []int{0}, lastTime: time.Now(), timeOut: 6 * time.Hour}
	return []int{0}
}

// TempRoute TempRoute
func (r *ARouter) TempRoute(src int, addr model.ANetAddr, outID int) {
	r.CLock.Lock()
	defer r.CLock.Unlock()
	r.Cache[strconv.Itoa(src)+"-"+addr.Host] = routeCache{data: []int{outID}, lastTime: time.Now(), timeOut: 10 * time.Minute}
}
