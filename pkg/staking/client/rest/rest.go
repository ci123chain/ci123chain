package rest

import (
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/gorilla/mux"
)

func RegisterRoutes(cliCtx context.Context, r *mux.Router) {
	RegisterRestNodeRoutes(cliCtx, r)
	RegisterQueryRoutes(cliCtx, r)
	RegisterRestTxRoutes(cliCtx, r)
}