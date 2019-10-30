package rpc

import (
	"gitlab.oneitfarm.com/blockchain/ci123chain/pkg/client/context"
	"github.com/gorilla/mux"
)

func RegisterRoutes(cliCtx context.Context, r *mux.Router)  {
	r.HandleFunc("/node_info", NodeInfoRequestHandlerFn(cliCtx)).Methods("GET")

}