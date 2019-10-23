package rpc

import (
	"CI123Chain/pkg/client/context"
	"github.com/gorilla/mux"
)

func RegisterRoutes(cliCtx context.Context, r *mux.Router)  {
	r.HandleFunc("/node_info", NodeInfoRequestHandlerFn(cliCtx)).Methods("GET")

}