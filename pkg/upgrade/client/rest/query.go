package rest

import (
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/upgrade/types"
	"github.com/gorilla/mux"
	"net/http"
)

// RegisterQueryRoutes - Central function to define routes that get registered by the main application
func RegisterQueryRoutes(cliCtx context.Context, r *mux.Router) {
	r.HandleFunc("/upgrade/current", upgradeCurrentQueryHandleFn(cliCtx, r))
	r.HandleFunc("/upgrade/applied", upgradeAppliedQueryHandleFn(cliCtx, r))

}

func upgradeCurrentQueryHandleFn(cliCtx context.Context, r *mux.Router) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		//vars := mux.Vars(request)
		//delegatorAddress := vars["delegatorAddr"]
		height := request.FormValue("height")
		prove := request.FormValue("prove")

		cliCtx, ok, err := rest.ParseQueryHeightOrReturnBadRequest(writer, cliCtx, request, "")
		if !ok || err != nil {
			rest.WriteErrorRes(writer, "get clictx failed")
			return
		}
		if !rest.CheckHeightAndProve(writer, height, prove, types.DefaultCodespace) {
			return
		}

		isProve := false
		if prove == "true" {
			isProve = true
		}
		res, _, proof, err := cliCtx.Query("/custom/" + types.ModuleName + "/" + types.QueryCurrent, nil, isProve)
		if err != nil {
			rest.WriteErrorRes(writer, err.Error())
			return
		}
		if len(res) < 1 {
			rest.WriteErrorRes(writer, fmt.Sprintf("unexpected res: %v", res))
			return
		}
		var plan types.Plan
		cliCtx.Cdc.MustUnmarshalJSON(res, &plan)
		value := plan
		resp := rest.BuildQueryRes(height, isProve, value, proof)
		rest.PostProcessResponseBare(writer, cliCtx, resp)
	}
}



func upgradeAppliedQueryHandleFn(cliCtx context.Context, r *mux.Router) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		//vars := mux.Vars(request)
		//delegatorAddress := vars["delegatorAddr"]
		height := request.FormValue("height")
		name := request.FormValue("name")
		prove := request.FormValue("prove")

		cliCtx, ok, err := rest.ParseQueryHeightOrReturnBadRequest(writer, cliCtx, request, "")
		if !ok || err != nil {
			rest.WriteErrorRes(writer, "get clictx failed")
			return
		}
		if !rest.CheckHeightAndProve(writer, height, prove, types.DefaultCodespace) {
			return
		}

		isProve := false
		if prove == "true" {
			isProve = true
		}
		param := types.QueryAppliedParams{Name: name}
		b := cliCtx.Cdc.MustMarshalJSON(param)
		res, _, proof, err := cliCtx.Query("/custom/" + types.ModuleName + "/" + types.QueryCurrent, b, isProve)
		if err != nil {
			rest.WriteErrorRes(writer, err.Error())
			return
		}
		if len(res) < 1 {
			rest.WriteErrorRes(writer, fmt.Sprintf("unexpected res: %v", res))
			return
		}
		var plan types.Plan
		cliCtx.Cdc.MustUnmarshalJSON(res, &plan)
		value := plan
		resp := rest.BuildQueryRes(height, isProve, value, proof)
		rest.PostProcessResponseBare(writer, cliCtx, resp)
	}
}

