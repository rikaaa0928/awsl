package manage

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Evi1/awsl/config"
	om "github.com/Evi1/awsl/object/manage"
)

var connectionNum = "cnum/"
var serverSide = "server/"
var clientSide = "client/"

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
		cnmBytes, err := json.Marshal(cnm)
		if err != nil {
			log.Panic(err)
		}
		res := string(cnmBytes) + " , "
		for _, v := range cnm {
			sum += v.Get()
		}
		res += strconv.FormatInt(sum, 10)
		w.Write([]byte(res))
	}
}

// Manage manage
func Manage() {
	http.ListenAndServe(":"+strconv.Itoa(config.Manage), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uri := r.RequestURI
		uri = strings.TrimPrefix(uri, "/")
		if strings.HasPrefix(uri, connectionNum) {
			uri = strings.TrimPrefix(uri, connectionNum)
			if strings.HasPrefix(uri, serverSide) {
				connectionNums(w, uri, true)
				return
			} else if strings.HasPrefix(uri, clientSide) {
				connectionNums(w, uri, false)
				return
			}
		}
		w.Write([]byte(uri))
	}))
}
