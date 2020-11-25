package rest

import (
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/vm/wasmtypes"
	"github.com/gorilla/mux"
)

func RegisterRoutes(cliCtx context.Context, r *mux.Router) {
	registerTxRoutes(cliCtx, r)
	registerQueryRoutes(cliCtx, r)
}

func registerTxRoutes(cliCtx context.Context, r *mux.Router)  {
	r.HandleFunc("/vm/contract/upload", rest.MiddleHandler(cliCtx, uploadContractHandler, types.DefaultCodespace)).Methods("POST")
	r.HandleFunc("/vm/contract/init", rest.MiddleHandler(cliCtx, instantiateContractHandler, types.DefaultCodespace)).Methods("POST")
	r.HandleFunc("/vm/contract/execute", rest.MiddleHandler(cliCtx, executeContractHandler, types.DefaultCodespace)).Methods("POST")
	r.HandleFunc("/vm/contract/migrate", rest.MiddleHandler(cliCtx, migrateContractHandler, types.DefaultCodespace)).Methods("POST")

}

func registerQueryRoutes(cliCtx context.Context, r *mux.Router) {
	//r.HandleFunc("/wasm/codeSearch/list", listCodesHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/vm/contract/meta", queryCodeHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/vm/contract/list", listContractsByCodeHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/vm/contract/info", queryContractHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/vm/contract/query", queryContractStateAllHandlerFn(cliCtx)).Methods("POST")
}
