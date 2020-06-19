package manage

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Evi1/awsl/config"
	sm "github.com/Evi1/awsl/servers/manage"
)

var connectionNum = "cnum/"
var serverSide = "server/"

// Manage manage
func Manage() {
	http.ListenAndServe(":"+strconv.Itoa(config.Manage), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uri := r.RequestURI
		uri = strings.TrimLeft(uri, "/")
		if strings.HasPrefix(uri, connectionNum) {
			uri = strings.TrimLeft(uri, connectionNum)
			if strings.HasPrefix(uri, serverSide) {
				uri = strings.TrimLeft(uri, serverSide)
				if len(uri) == 0 {
					var sum int64
					for _, v := range sm.ServerConnectionNumber {
						sum += v.Get()
					}
					w.Write([]byte(strconv.FormatInt(sum, 10)))
				}
			}
		}
	}))
}
