package object

import (
	"context"

	"github.com/rikaaa0928/awsl/clients"
	"github.com/rikaaa0928/awsl/model"
	"github.com/rikaaa0928/awsl/router"
	"github.com/rikaaa0928/awsl/servers"
	"github.com/rikaaa0928/awsl/tools"
)

// NewObject NewObject
func NewObject(conf model.Object) Object {
	closeWait := tools.NewCloseWait(context.Background())
	//config.MainContext = closeWait.Ctx
	return NewDefault(clients.NewClients(conf.Outs), servers.NewServers(closeWait.Ctx, conf.Ins), router.NewDefaultRouter(conf), closeWait)
}
