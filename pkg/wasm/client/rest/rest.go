package rest

import (
	"github.com/gorilla/mux"
	"github.com/tanhuiya/ci123chain/pkg/client/context"
)

func RegisterRoutes(cliCtx context.Context, r *mux.Router) {
	registerTxRoutes(cliCtx, r)
	registerQueryRoutes(cliCtx, r)
}