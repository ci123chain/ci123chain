package rest

import (
	"CI123Chain/pkg/abci/types/rest"
	"CI123Chain/pkg/client/context"
	"CI123Chain/pkg/client/helper"
	"github.com/gorilla/mux"
	"net/http"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.Context, r *mux.Router) {
	//r.HandleFunc("/bank/accounts/{address}/transfers", SendRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/bank/balances/{address}", QueryBalancesRequestHandlerFn(cliCtx)).Methods("GET")
}

func QueryBalancesRequestHandlerFn(cliCtx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(request)
		addr := vars["address"]

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, request)
		if !ok {
			return
		}
		addrBytes, err := helper.ParseAddrs(addr)
		if len(addrBytes) < 1 || err != nil {
			return
		}
		//params := types.NewQueryBalanceParams(addr)
		res, err := cliCtx.GetBalanceByAddress(addrBytes[0])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		rest.PostProcessResponseBare(w, cliCtx, res)
	}
}