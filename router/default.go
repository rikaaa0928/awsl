package router

import "github.com/Evi1/awsl/servers"

// ARouter ARouter
type ARouter struct{}

// Route Route
func (r ARouter) Route(_ servers.ANetAddr) int {
	return 0
}
