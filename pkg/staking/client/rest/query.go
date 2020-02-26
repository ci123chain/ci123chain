package rest

import (
	"github.com/gorilla/mux"
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/abci/types/rest"
	"github.com/tanhuiya/ci123chain/pkg/client/context"
	"github.com/tanhuiya/ci123chain/pkg/staking/types"
	"github.com/tanhuiya/ci123chain/pkg/transfer"
	"net/http"
)

func RegisterTxRoutes(cliCtx context.Context, r *mux.Router) {
	// Get all validators
	r.HandleFunc("/staking/validators", validatorsHandlerFn(cliCtx), ).Methods("POST")
	// Get a single validator info
	r.HandleFunc("/staking/validator", validatorHandlerFn(cliCtx), ).Methods("POST")
	// Get all delegations to a validator
	r.HandleFunc("/staking/validator/delegations", validatorDelegationsHandlerFn(cliCtx), ).Methods("POST")
	// Query all validators that a delegator is bonded to
	r.HandleFunc("/staking/delegator/validators", delegatorValidatorsHandlerFn(cliCtx), ).Methods("POST")

	// Query a validator that a delegator is bonded to
	r.HandleFunc("/staking/delegator/validator", delegatorValidatorHandlerFn(cliCtx), ).Methods("POST")

	// Query a delegation between a delegator and a validator
	r.HandleFunc("/staking/delegator/delegation", delegationHandlerFn(cliCtx), ).Methods("POST")
	// Get all delegations from a delegator
	r.HandleFunc("/staking/delegator/delegations", delegatorDelegationsHandlerFn(cliCtx), ).Methods("POST")
	// Query redelegations (filters in query params)
	r.HandleFunc("/staking/redelegations", redelegationsHandlerFn(cliCtx), ).Methods("POST")

}

func delegatorDelegationsHandlerFn(cliCtx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		//vars := mux.Vars(request)
		//delegatorAddress := vars["delegatorAddr"]
		delegatorAddress := request.FormValue("delegatorAddr")
		delegatorAddr := sdk.HexToAddress(delegatorAddress)
		params := types.NewQueryDelegatorParams(delegatorAddr)
		bz, Err := cliCtx.Cdc.MarshalJSON(params)
		if Err != nil {
			rest.WriteErrorRes(writer, sdk.ErrInternal("marshal failed"))
			return
		}
		cliCtx, ok, err := rest.ParseQueryHeightOrReturnBadRequest(writer, cliCtx, request, "")
		if !ok {
			rest.WriteErrorRes(writer, err)
			return
		}

		res, _, err := cliCtx.Query("/custom/" + types.ModuleName + "/delegatorDelegations/", bz)
		if err != nil {
			rest.WriteErrorRes(writer, err)
			return
		}
		if len(res) < 1 {
			rest.WriteErrorRes(writer, transfer.ErrQueryTx(types.DefaultCodespace, "query response length less than 1"))
			return
		}
		var delegations types.DelegationResponses
		cliCtx.Cdc.MustUnmarshalJSON(res, &delegations)
		resp := delegations
		rest.PostProcessResponseBare(writer, cliCtx, resp)
	}
}

func validatorsHandlerFn(cliCtx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		cliCtx, ok, err := rest.ParseQueryHeightOrReturnBadRequest(writer, cliCtx, request, "")
		if !ok {
			rest.WriteErrorRes(writer, err)
			return
		}

		_, page, limit, Err := rest.ParseHTTPArgsWithLimit(request, 0)
		if err != nil {
			rest.WriteErrorRes(writer, err)
			return
		}

		status := request.FormValue("status")
		if status == "" {
			status = sdk.BondStatusBonded
		}

		params := types.NewQueryValidatorsParams(page, limit, status)
		bz, Err := cliCtx.Cdc.MarshalJSON(params)
		if Err != nil {
			rest.WriteErrorRes(writer, sdk.ErrInternal("marshal failed"))
			return
		}

		res, _, err := cliCtx.Query("/custom/" + types.ModuleName + "/validators/", bz)
		if err != nil {
			rest.WriteErrorRes(writer, err)
			return
		}
		if len(res) < 1 {
			rest.WriteErrorRes(writer, transfer.ErrQueryTx(types.DefaultCodespace, "query response length less than 1"))
			return
		}
		var validators []types.Validator
		cliCtx.Cdc.MustUnmarshalJSON(res, &validators)
		/*Err = types.StakingCodec.UnmarshalJSON(res, &validators)
		if Err != nil {
			//rest.WriteErrorRes(writer, transfer.ErrQueryTx(types.DefaultCodespace, "unmarshal failed"))
			//return
		}*/
		resp := validators
		rest.PostProcessResponseBare(writer, cliCtx, resp)
	}
}

// HTTP request handler to query the validator information from a given validator address
func validatorHandlerFn(cliCtx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		cliCtx, ok, err := rest.ParseQueryHeightOrReturnBadRequest(writer, cliCtx, request, "")
		if !ok {
			rest.WriteErrorRes(writer, err)
			return
		}
		//vars := mux.Vars(request)
		//validatorAddress := vars["validatorAddr"]
		validatorAddress := request.FormValue("validatorAddr")
		validatorAddr := sdk.HexToAddress(validatorAddress)
		params := types.NewQueryValidatorParams(validatorAddr)
		bz, Err := cliCtx.Cdc.MarshalJSON(params)
		if Err != nil {
			rest.WriteErrorRes(writer, sdk.ErrInternal("marshal failed"))
			return
		}

		res, _, err := cliCtx.Query("/custom/" + types.ModuleName + "/validator/", bz)
		if err != nil {
			rest.WriteErrorRes(writer, err)
			return
		}
		if len(res) < 1 {
			rest.WriteErrorRes(writer, transfer.ErrQueryTx(types.DefaultCodespace, "query response length less than 1"))
			return
		}
		var validator types.Validator
		cliCtx.Cdc.MustUnmarshalJSON(res, &validator)
		resp := validator
		rest.PostProcessResponseBare(writer, cliCtx, resp)
	}
}

// HTTP request handler to query all unbonding delegations from a validator
func validatorDelegationsHandlerFn(cliCtx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		//vars := mux.Vars(request)
		validatorAddress := request.FormValue("validatorAddr")
		validatorAddr := sdk.HexToAddress(validatorAddress)
		params := types.NewQueryValidatorParams(validatorAddr)
		bz, Err := cliCtx.Cdc.MarshalJSON(params)
		if Err != nil {
			rest.WriteErrorRes(writer, sdk.ErrInternal("marshal failed"))
			return
		}

		cliCtx, ok, err := rest.ParseQueryHeightOrReturnBadRequest(writer, cliCtx, request, "")
		if !ok {
			rest.WriteErrorRes(writer, err)
			return
		}

		res, _, err := cliCtx.Query("/custom/" + types.ModuleName + "/validatorDelegations/", bz)
		if err != nil {
			rest.WriteErrorRes(writer, err)
			return
		}
		if len(res) < 1 {
			rest.WriteErrorRes(writer, transfer.ErrQueryTx(types.DefaultCodespace, "query response length less than 1"))
			return
		}
		var delegations types.DelegationResponses
		cliCtx.Cdc.MustUnmarshalJSON(res, &delegations)
		resp := delegations
		rest.PostProcessResponseBare(writer, cliCtx, resp)
	}
}

// HTTP request handler to query all delegator bonded validators
func delegatorValidatorsHandlerFn(cliCtx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		//vars := mux.Vars(request)
		//delegatorAddress := vars["delegatorAddr"]
		delegatorAddress := request.FormValue("delegatorAddr")
		delegatorAddr := sdk.HexToAddress(delegatorAddress)
		params := types.NewQueryDelegatorParams(delegatorAddr)
		bz, Err := cliCtx.Cdc.MarshalJSON(params)
		if Err != nil {
			rest.WriteErrorRes(writer, sdk.ErrInternal("marshal failed"))
			return
		}

		cliCtx, ok, err := rest.ParseQueryHeightOrReturnBadRequest(writer, cliCtx, request, "")
		if !ok {
			rest.WriteErrorRes(writer, err)
			return
		}

		res, _, err := cliCtx.Query("/custom/" + types.ModuleName + "/delegatorValidators/", bz)
		if err != nil {
			rest.WriteErrorRes(writer, err)
			return
		}
		if len(res) < 1 {
			rest.WriteErrorRes(writer, transfer.ErrQueryTx(types.DefaultCodespace, "query response length less than 1"))
			return
		}
		var validators []types.Validator
		cliCtx.Cdc.MustUnmarshalJSON(res, &validators)
		resp := validators
		rest.PostProcessResponseBare(writer, cliCtx, resp)
	}
}

// HTTP request handler to get information from a currently bonded validator
func delegatorValidatorHandlerFn(cliCtx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		//vars := mux.Vars(request)
		//validatorAddress := vars["validatorAddr"]
		//delegatorAddress := vars["delegatorAddr"]
		validatorAddress := request.FormValue("validatorAddr")
		delegatorAddress := request.FormValue("delegatorAddr")
		validatorAddr := sdk.HexToAddress(validatorAddress)
		delegatorAddr := sdk.HexToAddress(delegatorAddress)
		params := types.NewQueryBondsParams(delegatorAddr, validatorAddr)
		bz, Err := cliCtx.Cdc.MarshalJSON(params)
		if Err != nil {
			rest.WriteErrorRes(writer, sdk.ErrInternal("marshal failed"))
			return
		}

		cliCtx, ok, err := rest.ParseQueryHeightOrReturnBadRequest(writer, cliCtx, request, "")
		if !ok {
			rest.WriteErrorRes(writer, err)
			return
		}

		res, _, err := cliCtx.Query("/custom/" + types.ModuleName + "/delegatorValidator/", bz)
		if err != nil {
			rest.WriteErrorRes(writer, err)
			return
		}
		if len(res) < 1 {
			rest.WriteErrorRes(writer, transfer.ErrQueryTx(types.DefaultCodespace, "query response length less than 1"))
			return
		}
		var validator types.Validator
		cliCtx.Cdc.MustUnmarshalJSON(res, &validator)
		resp := validator
		rest.PostProcessResponseBare(writer, cliCtx, resp)
	}
}


// HTTP request handler to query a delegation
func delegationHandlerFn(cliCtx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		//vars := mux.Vars(request)
		//validatorAddress := vars["validatorAddr"]
		//delegatorAddress := vars["delegatorAddr"]
		validatorAddress := request.FormValue("validatorAddr")
		delegatorAddress := request.FormValue("delegatorAddr")
		validatorAddr := sdk.HexToAddress(validatorAddress)
		delegatorAddr := sdk.HexToAddress(delegatorAddress)
		params := types.NewQueryBondsParams(delegatorAddr, validatorAddr)
		bz, Err := cliCtx.Cdc.MarshalJSON(params)
		if Err != nil {
			rest.WriteErrorRes(writer, sdk.ErrInternal("marshal failed"))
			return
		}

		cliCtx, ok, err := rest.ParseQueryHeightOrReturnBadRequest(writer, cliCtx, request, "")
		if !ok {
			rest.WriteErrorRes(writer, err)
			return
		}

		res, _, err := cliCtx.Query("/custom/" + types.ModuleName + "/delegation/", bz)
		if err != nil {
			rest.WriteErrorRes(writer, err)
			return
		}
		if len(res) < 1 {
			rest.WriteErrorRes(writer, transfer.ErrQueryTx(types.DefaultCodespace, "query response length less than 1"))
			return
		}
		var delegation types.DelegationResponse
		cliCtx.Cdc.MustUnmarshalJSON(res, &delegation)
		resp := delegation
		rest.PostProcessResponseBare(writer, cliCtx, resp)
	}
}

// HTTP request handler to query redelegations
func redelegationsHandlerFn(cliCtx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var params types.QueryRedelegationParams

		cliCtx, ok, err := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r, "")
		if !ok {
			rest.WriteErrorRes(w, err)
			return
		}
		//vars := mux.Vars(r)
		//delegatorAddress := vars["delegatorAddr"]
		//validatorSrc := vars["validatorSrcAddr"]
		//validatorDst := vars["validatorDsrAddr"]
		delegatorAddress := r.FormValue("delegatorAddr")
		validatorSrc := r.FormValue("validatorSrcAddr")
		validatorDst := r.FormValue("validatorDsrAddr")

		delegatorAddr := sdk.HexToAddress(delegatorAddress)
		validatorSrcAddr := sdk.HexToAddress(validatorSrc)
		validatorDstAddr := sdk.HexToAddress(validatorDst)

		params.DelegatorAddr = delegatorAddr
		params.SrcValidatorAddr = validatorSrcAddr
		params.DstValidatorAddr = validatorDstAddr
		bz, Err := cliCtx.Cdc.MarshalJSON(params)
		if Err != nil {
			rest.WriteErrorRes(w, sdk.ErrInternal("marshal failed"))
			return
		}

		res, _, err := cliCtx.Query("/custom/" + types.ModuleName + "/redelegations/", bz)
		if err != nil {
			rest.WriteErrorRes(w, err)
			return
		}
		if len(res) < 1 {
			rest.WriteErrorRes(w, transfer.ErrQueryTx(types.DefaultCodespace, "query response length less than 1"))
			return
		}
		var redelegations types.RedelegationResponses
		cliCtx.Cdc.MustUnmarshalJSON(res, &redelegations)
		resp := redelegations
		rest.PostProcessResponseBare(w, cliCtx, resp)
	}
}