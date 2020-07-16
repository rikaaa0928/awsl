package manage

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Evi1/awsl/config"
	"github.com/Evi1/awsl/object"
	om "github.com/Evi1/awsl/object/manage"
)

var base = "/manage/"

var connectionNum = "cnum/"
var routerCache = "routercache/"
var serverFlow = "serverflow/"

var serverSide = "server/"
var clientSide = "client/"
var history = "history/"

var obj object.Object

type sConnNum struct {
	Tag string
	Num int64
}

func handleConnectionNums(w http.ResponseWriter, uri string, isServer bool) {
	cnm := om.ServerConnectionNumber
	if !isServer {
		cnm = om.ClientConnectionNumber
	}
	end := serverSide
	if !isServer {
		end = clientSide
	}
	uri = strings.TrimPrefix(uri, end)
	res := make(map[int]sConnNum, 0)
	if len(uri) == 0 {
		var sum int64
		//res := ""
		for k, v := range cnm {
			n := v.Get()
			sum += n
			//res += strconv.Itoa(k) + "-" + v.Tag + " : " + strconv.FormatInt(n, 10) + " , "
			res[k] = sConnNum{Tag: v.Tag, Num: n}
		}
		//res += "sum : " + strconv.FormatInt(sum, 10)
		res[-1] = sConnNum{Tag: "sum", Num: sum}
		resBytes, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			w.Write([]byte(err.Error() + "\n"))
		}
		w.Write(resBytes)
		return
	}
	id, err := strconv.Atoi(uri)
	if err != nil {
		w.Write([]byte(err.Error() + "\n"))
		return
	}
	num, ok := cnm[id]
	if !ok {
		w.Write([]byte("no id : " + strconv.Itoa(id)))
		return
	}
	w.Write([]byte(strconv.FormatInt(num.Counter.Get(), 10)))
}

func handleRouterCache(w http.ResponseWriter, uri string) {
	uri = strings.TrimPrefix(uri, routerCache)
	if len(uri) == 0 {
		if o, ok := obj.(*object.DefaultObject); ok {
			res := o.R.GetCache(-1)
			w.Write([]byte(res))
			return
		}
	}
	src, err := strconv.Atoi(uri)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	if o, ok := obj.(*object.DefaultObject); ok {
		res := o.R.GetCache(src)
		w.Write([]byte(res))
	}
}

func handleServerFlow(w http.ResponseWriter, uri string) {
	uri = strings.TrimPrefix(uri, serverFlow)
	if len(uri) == 0 {
		res := om.ServerFlowManager.GetRoot()
		w.Write([]byte(res))
		return
	}
	src, err := strconv.Atoi(uri)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	res := om.ServerFlowManager.GetID(src)
	w.Write([]byte(res))
}

// Manage manage
func Manage(o object.Object) {
	if config.Manage <= 0 {
		return
	}
	obj = o
	http.ListenAndServe(":"+strconv.Itoa(config.Manage), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uri := r.RequestURI
		if !strings.HasPrefix(uri, base) {
			w.Write([]byte(uri))
			return
		}
		uri = strings.TrimPrefix(uri, base)
		if strings.HasPrefix(uri, connectionNum) {
			uri = strings.TrimPrefix(uri, connectionNum)
			if strings.HasPrefix(uri, serverSide) {
				handleConnectionNums(w, uri, true)
				return
			} else if strings.HasPrefix(uri, clientSide) {
				handleConnectionNums(w, uri, false)
				return
			}
		} else if strings.HasPrefix(uri, routerCache) {
			handleRouterCache(w, uri)
			return
		} else if strings.HasPrefix(uri, serverFlow) {
			handleServerFlow(w, uri)
			return
		}
		w.Write([]byte(uri))
	}))
}
