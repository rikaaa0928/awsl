package router

import "github.com/Evi1/awsl/model"

// ARouter ARouter
type ARouter struct{}

// Route Route
func (r ARouter) Route(_ model.ANetAddr) int {
	return 0
}
