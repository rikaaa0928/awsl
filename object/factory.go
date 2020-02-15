package object

import (
	"github.com/Evi1/awsl/clients"
	"github.com/Evi1/awsl/model"
	"github.com/Evi1/awsl/servers"
)

// NewObject NewObject
func NewObject(conf model.Object) Object {
	return NewDefault(clients.NewClients(conf.Outs), servers.NewServers(conf.Ins))
}
