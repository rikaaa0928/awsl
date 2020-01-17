package router

import "github.com/Evi1/awsl/servers"

type ARouter struct{}

func (r ARouter) Route(_ servers.ANetAddr) int {
	return 0
}
