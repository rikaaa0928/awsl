package router

import "github.com/rikaaa0928/awsl/model"

// Router router
type Router interface {
	Route(src int, addr model.ANetAddr) []int
	GetCache(src int) string
}
