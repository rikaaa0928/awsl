package router

import "github.com/Evi1/awsl/servers"

// Router router
type Router interface {
	Route(addr servers.ANetAddr) int
}
