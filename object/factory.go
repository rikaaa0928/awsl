package object

import (
	"context"

	"github.com/Evi1/awsl/clients"
	"github.com/Evi1/awsl/model"
	"github.com/Evi1/awsl/router"
	"github.com/Evi1/awsl/servers"
	"github.com/Evi1/awsl/tools"
)

// NewObject NewObject
func NewObject(conf model.Object) Object {
	closeWait := tools.NewCloseWait(context.Background())
	//config.MainContext = closeWait.Ctx
	return NewDefault(clients.NewClients(conf.Outs), servers.NewServers(closeWait.Ctx, conf.Ins), router.NewDefaultRouter(conf), closeWait)
}
