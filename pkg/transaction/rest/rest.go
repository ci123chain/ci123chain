package rest

import (
	"github.com/tanhuiya/ci123chain/pkg/abci/types/rest"
	"github.com/tanhuiya/ci123chain/pkg/client/context"
	"github.com/tanhuiya/ci123chain/pkg/transaction/rest/utils"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

func RegisterTxRoutes(cliCtx context.Context, r *mux.Router)  {
	r.HandleFunc("/tx/{hash}", QueryTxRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/tx/sign", SignTxRequestHandler(cliCtx)).Methods("POST")
	r.HandleFunc("/tx/broadcast", BraodcastTxRequest(cliCtx)).Methods("POST")
}

func QueryTxRequestHandlerFn(cliCtx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		hashHexStr := vars["hash"]

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(writer, cliCtx, request)
		if !ok {
			return
		}

		output, err := utils.QueryTx(cliCtx, hashHexStr)
		if err != nil {
			if strings.Contains(err.Error(), hashHexStr) {
				rest.WriteErrorResponse(writer, http.StatusNotFound, err.Error())
				return
			}
			rest.WriteErrorResponse(writer, http.StatusInternalServerError, err.Error())
			return
		}
		if output.Empty() {
			rest.WriteErrorResponse(writer, http.StatusNotFound, fmt.Sprintf("no transaction found with hash %s", hashHexStr))
		}
		rest.PostProcessResponseBare(writer, cliCtx, output)
	}
}