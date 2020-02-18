package rest

import (
	"github.com/gorilla/mux"
	"github.com/tanhuiya/ci123chain/pkg/abci/types/rest"
	"github.com/tanhuiya/ci123chain/pkg/client/context"
	"github.com/tanhuiya/ci123chain/pkg/staking/types"
	"github.com/tanhuiya/ci123chain/pkg/transfer"
	"net/http"
)

func RegisterTxRoutes(cliCtx context.Context, r *mux.Router) {
	r.HandleFunc("/staking/getDelegation", QueryDelegationRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/staking/getAllDelegation", QueryAllDelegationRequestHandlerFn(cliCtx)).Methods("POST")
}

func QueryDelegationRequestHandlerFn(cliCtx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		//
		vars := mux.Vars(request)
		validatorAddress := vars["validatorAddress"]
		delegatorAddress := vars["delegatorAddress"]
		cliCtx, ok, err := rest.ParseQueryHeightOrReturnBadRequest(writer, cliCtx, request, "")
		if !ok {
			rest.WriteErrorRes(writer, err)
			return
		}

		res, _, err := cliCtx.Query("/custom/" + types.ModuleName + "/delegation/" + validatorAddress + "/" + delegatorAddress, nil)
		if err != nil {
			rest.WriteErrorRes(writer, err)
			return
		}
		if len(res) < 1 {
			rest.WriteErrorRes(writer, transfer.ErrQueryTx(types.DefaultCodespace, "query response length less than 1"))
			return
		}
		rest.PostProcessResponseBare(writer, cliCtx, res)
	}
}

func QueryAllDelegationRequestHandlerFn(cliCtx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		//
		vars := mux.Vars(request)
		validatorAddress := vars["validatorAddress"]
		cliCtx, ok, err := rest.ParseQueryHeightOrReturnBadRequest(writer, cliCtx, request, "")
		if !ok {
			rest.WriteErrorRes(writer, err)
			return
		}

		res, _, err := cliCtx.Query("/custom/" + types.ModuleName + "/allDelegation/" + validatorAddress, nil)
		if err != nil {
			rest.WriteErrorRes(writer, err)
			return
		}
		if len(res) < 1 {
			rest.WriteErrorRes(writer, transfer.ErrQueryTx(types.DefaultCodespace, "query response length less than 1"))
			return
		}
		rest.PostProcessResponseBare(writer, cliCtx, res)
	}
}