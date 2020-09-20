package httpservers

import (
	"net/http"
	"sync"
)

var hServers map[int]*http.ServeMux
var hLock sync.Mutex
var started bool

func init() {
	hLock = sync.Mutex{}
	hServers = make(map[int]*http.ServeMux)
}
