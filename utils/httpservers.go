package utils

import "net/http"

var hServers map[int]*http.ServeMux

func init() {
	hServers = make(map[int]*http.ServeMux)
}
