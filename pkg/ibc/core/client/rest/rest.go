package rest

import (
	"github.com/ci123chain/ci123chain/pkg/client/context"
	ibcchannel "github.com/ci123chain/ci123chain/pkg/ibc/core/channel/client/rest"
	ibcclient "github.com/ci123chain/ci123chain/pkg/ibc/core/clients/client/rest"
	ibcconnection "github.com/ci123chain/ci123chain/pkg/ibc/core/connection/client/rest"
	"github.com/gorilla/mux"
)

func RegisterRoutes(cliCtx context.Context, r *mux.Router) {
	ibcchannel.RegisterQueryRoutes(cliCtx, r)
	ibcconnection.RegisterQueryRoutes(cliCtx, r)
	ibcclient.RegisterQueryRoutes(cliCtx, r)
}