package rest

import (
	"github.com/gorilla/mux"
	"github.com/tanhuiya/ci123chain/pkg/abci/types/rest"
	"github.com/tanhuiya/ci123chain/pkg/account/types"
	"github.com/tanhuiya/ci123chain/pkg/client"
	"github.com/tanhuiya/ci123chain/pkg/client/context"
	"github.com/tanhuiya/ci123chain/pkg/client/helper"
	"github.com/tanhuiya/ci123chain/pkg/transfer"
	"net/http"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.Context, r *mux.Router) {

	r.HandleFunc("/bank/balances/{address}", QueryBalancesRequestHandlerFn(cliCtx)).Methods("GET")
}

type BalanceData struct {
	Balance uint64 `json:"balance"`
}
func QueryBalancesRequestHandlerFn(cliCtx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(request)
		addr := vars["address"]

		cliCtx, ok, err := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, request)
		if !ok {
			rest.WriteErrorRes(w, err)
			return
		}
		addrBytes, err2 := helper.ParseAddrs(addr)
		if len(addrBytes) < 1 || err2 != nil {
			rest.WriteErrorRes(w, client.ErrParseAddr(types.DefaultCodespace, err2))
			return
		}
		//params := types.NewQueryBalanceParams(addr)
		res, err2 := cliCtx.GetBalanceByAddress(addrBytes[0])
		if err2 != nil {
			rest.WriteErrorRes(w, transfer.ErrQueryTx(types.DefaultCodespace, err2.Error()))
			return
		}
		resp := BalanceData{Balance:res}
		rest.PostProcessResponseBare(w, cliCtx, resp)
	}
}