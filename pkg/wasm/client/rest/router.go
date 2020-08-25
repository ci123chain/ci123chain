package rest

import (
	"github.com/gorilla/mux"

	"github.com/ci123chain/ci123chain/pkg/client/context"
)

func RegisterRoutes(cliCtx context.Context, r *mux.Router) {
	registerTxRoutes(cliCtx, r)
	registerQueryRoutes(cliCtx, r)
}

func registerTxRoutes(cliCtx context.Context, r *mux.Router) {
	r.HandleFunc("/wasm/contract/init", instantiateContractHandler(cliCtx)).Methods("POST")
	r.HandleFunc("/wasm/contract/execute", executeContractHandler(cliCtx)).Methods("POST")
	r.HandleFunc("/wasm/contract/migrate", migrateContractHandler(cliCtx)).Methods("POST")
}


func registerQueryRoutes(cliCtx context.Context, r *mux.Router) {
	//r.HandleFunc("/wasm/codeSearch/list", listCodesHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/wasm/contract/meta", queryCodeHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/wasm/contract/list", listContractsByCodeHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/wasm/contract/info", queryContractHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/wasm/contract/query", queryContractStateAllHandlerFn(cliCtx)).Methods("POST")
}
