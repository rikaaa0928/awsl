package router

import "github.com/Evi1/awsl/model"

// Router router
type Router interface {
	Route(src int, addr model.ANetAddr) int
}
