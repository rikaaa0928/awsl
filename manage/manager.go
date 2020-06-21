package manage

import (
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
var serverSide = "server/"
var clientSide = "client/"

var obj object.Object

func connectionNums(w http.ResponseWriter, uri string, isServer bool) {
	cnm := om.ServerConnectionNumber
	if !isServer {
		cnm = om.ClientConnectionNumber
	}
	end := serverSide
	if !isServer {
		end = clientSide
	}
	uri = strings.TrimPrefix(uri, end)
	if len(uri) == 0 {
		var sum int64
		res := ""
		for k, v := range cnm {
			n := v.Get()
			sum += v.Get()
			res += strconv.Itoa(k) + " : " + strconv.FormatInt(n, 10) + " , "
		}
		res += strconv.FormatInt(sum, 10)
		w.Write([]byte(res))
	}
}

func routerCaches(w http.ResponseWriter, uri string) {
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

// Manage manage
func Manage(o object.Object) {
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
				connectionNums(w, uri, true)
				return
			} else if strings.HasPrefix(uri, clientSide) {
				connectionNums(w, uri, false)
				return
			}
		} else if strings.HasPrefix(uri, routerCache) {
			routerCaches(w, uri)
			return
		}
		w.Write([]byte(uri))
	}))
}
