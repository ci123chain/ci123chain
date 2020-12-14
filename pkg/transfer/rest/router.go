package rest

import (
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	types2 "github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/transfer/types"
	"github.com/gorilla/mux"
)

var cdc = types2.MakeCodec()
func RegisterTxRoutes(cliCtx context.Context, r *mux.Router)  {
	r.HandleFunc("/tx/query", QueryTxRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/tx/transfers", rest.MiddleHandler(cliCtx, SendRequestHandlerFn, types.DefaultCodespace)).Methods("POST")
	r.HandleFunc("/tx/broadcast", BroadcastTxRequest(cliCtx)).Methods("POST")
	r.HandleFunc("/tx/broadcast_async", BroadcastTxRequestAsync(cliCtx)).Methods("POST")
	r.HandleFunc("/tx/query/block_txs", QueryTxsWithHeight(cliCtx)).Methods("POST")
}
