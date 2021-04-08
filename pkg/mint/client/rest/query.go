package rest


import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/mint/types"
	"github.com/gorilla/mux"
	"net/http"
)


func RegisterQueryRoutes(cliCtx context.Context, r *mux.Router) {
	r.HandleFunc("/mint/parameters", queryParamsHandleFn(cliCtx)).Methods("POST")
	r.HandleFunc("/mint/inflation", queryInflationHandleFn(cliCtx)).Methods("POST")
	r.HandleFunc("/mint/annual_provisions", queryAnnualProvisionsHandlerFn(cliCtx)).Methods("POST")
}


func queryParamsHandleFn(cliCtx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("/custom/%s/%s", types.QuerierRoute, types.QueryParameters)
		res, _, _, err := cliCtx.Query(route, nil, false)

		if err != nil {
			rest.WriteErrorRes(w, err.Error())
			return
		}
		var params types.Params
		cliCtx.Cdc.MustUnmarshalJSON(res, &params)
		rest.PostProcessResponseBare(w, cliCtx, params)
	}
}



func queryInflationHandleFn(cliCtx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("/custom/%s/%s", types.QuerierRoute, types.QueryInflation)
		res, _, _, err := cliCtx.Query(route, nil, false)

		if err != nil {
			rest.WriteErrorRes(w, err.Error())
			return
		}
		var inflation sdk.Dec
		cliCtx.Cdc.MustUnmarshalJSON(res, &inflation)
		rest.PostProcessResponseBare(w, cliCtx, inflation)
	}
}




func queryAnnualProvisionsHandlerFn(cliCtx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("/custom/%s/%s", types.QuerierRoute, types.QueryAnnualProvisions)
		res, _, _, err := cliCtx.Query(route, nil, false)

		if err != nil {
			rest.WriteErrorRes(w, err.Error())
			return
		}
		var annualPrivisions sdk.Dec
		cliCtx.Cdc.MustUnmarshalJSON(res, &annualPrivisions)
		rest.PostProcessResponseBare(w, cliCtx, annualPrivisions)
	}
}
