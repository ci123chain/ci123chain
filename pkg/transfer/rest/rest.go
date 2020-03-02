package rest

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tanhuiya/ci123chain/pkg/abci/types/rest"
	"github.com/tanhuiya/ci123chain/pkg/client/context"
	"github.com/tanhuiya/ci123chain/pkg/transfer/rest/utils"
	"github.com/tanhuiya/ci123chain/pkg/transfer/types"
	"net/http"
)

func RegisterTxRoutes(cliCtx context.Context, r *mux.Router)  {
	r.HandleFunc("/tx", QueryTxRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/tx/sign_transfer", SignTxRequestHandler(cliCtx)).Methods("POST")
	r.HandleFunc("/tx/transfers", SendRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/tx/broadcast", BroadcastTxRequest(cliCtx)).Methods("POST")
	r.HandleFunc("/tx/broadcast_async", BroadcastTxRequestAsync(cliCtx)).Methods("POST")
}

type TxRequestParams struct {
	Hash    string    `json:"hash"`
	Height  string    `json:"height"`
}

/*type QueryTxRequestParams struct{
	//
	Data    TxRequestParams `json:"data"`
}*/

func QueryTxRequestHandlerFn(cliCtx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		//vars := mux.Vars(request)
		hashHexStr := request.FormValue("hash")

		/*
		var params QueryTxRequestParams
		b, readErr := ioutil.ReadAll(request.Body)
		readErr = json.Unmarshal(b, &params)
		if readErr != nil {
			//
		}
		*/

		cliCtx, ok, err := rest.ParseQueryHeightOrReturnBadRequest(writer, cliCtx, request, "")
		if !ok {
			rest.WriteErrorRes(writer, err)
			return
		}

		resp, err := utils.QueryTx(cliCtx, hashHexStr)
		if err != nil {
			rest.WriteErrorRes(writer, err)
			return
		}
		if resp.Empty() {
			rest.WriteErrorRes(writer, types.ErrQueryTx(types.DefaultCodespace,fmt.Sprintf("no transfer found with hash %s", hashHexStr)))
		}
		rest.PostProcessResponseBare(writer, cliCtx, resp)
	}
}