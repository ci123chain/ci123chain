package rpc

import (
	"github.com/gorilla/mux"
	"github.com/ci123chain/ci123chain/pkg/client/context"
)

func RegisterRoutes(cliCtx context.Context, r *mux.Router)  {

	r.HandleFunc("/node_info", NodeInfoRequestHandlerFn(cliCtx)).Methods("GET")

}

