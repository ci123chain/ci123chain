package rest

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/staking/types"
	"github.com/ci123chain/ci123chain/pkg/transfer"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterTxRoutes(cliCtx context.Context, r *mux.Router) {
	// Get all validators
	r.HandleFunc("/staking/validator/all", validatorsHandlerFn(cliCtx), ).Methods("POST")
	// Get a single validator info
	r.HandleFunc("/staking/validator/info", validatorHandlerFn(cliCtx), ).Methods("POST")
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
		height := request.FormValue("height")
		prove := request.FormValue("prove")
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
		if !rest.CheckHeightAndProve(writer, height, prove, types.DefaultCodespace) {
			return
		}

		isProve := false
		if prove == "true" {
			isProve = true
		}
		res, _, proof, err := cliCtx.Query("/custom/" + types.ModuleName + "/" + types.QueryDelegatorDelegations, bz, isProve)
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
		value := delegations
		resp := rest.BuildQueryRes(height, isProve, value, proof)
		rest.PostProcessResponseBare(writer, cliCtx, resp)
	}
}

func validatorsHandlerFn(cliCtx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		height := request.FormValue("height")
		prove := request.FormValue("prove")
		cliCtx, ok, err := rest.ParseQueryHeightOrReturnBadRequest(writer, cliCtx, request, "")
		if !ok {
			rest.WriteErrorRes(writer, err)
			return
		}
		if !rest.CheckHeightAndProve(writer, height, prove, types.DefaultCodespace) {
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

		isProve := false
		if prove == "true" {
			isProve = true
		}
		res, _, proof, err := cliCtx.Query("/custom/" + types.ModuleName + "/" + types.QueryValidators, bz, isProve)
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
		value := validators
		resp := rest.BuildQueryRes(height, isProve, value, proof)
		rest.PostProcessResponseBare(writer, cliCtx, resp)
	}
}

// HTTP request handler to query the validator information from a given validator address
func validatorHandlerFn(cliCtx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		height := request.FormValue("height")
		prove := request.FormValue("prove")
		cliCtx, ok, err := rest.ParseQueryHeightOrReturnBadRequest(writer, cliCtx, request, "")
		if !ok {
			rest.WriteErrorRes(writer, err)
			return
		}
		if !rest.CheckHeightAndProve(writer, height, prove, types.DefaultCodespace) {
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

		isProve := false
		if prove == "true" {
			isProve = true
		}
		res, _, proof, err := cliCtx.Query("/custom/" + types.ModuleName + "/" + types.QueryValidator, bz, isProve)
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
		value := validator
		resp := rest.BuildQueryRes(height, isProve, value, proof)
		rest.PostProcessResponseBare(writer, cliCtx, resp)
	}
}

// HTTP request handler to query all unbonding delegations from a validator
func validatorDelegationsHandlerFn(cliCtx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		height := request.FormValue("height")
		prove := request.FormValue("prove")
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
		if !rest.CheckHeightAndProve(writer, height, prove, types.DefaultCodespace) {
			return
		}

		isProve := false
		if prove == "true" {
			isProve = true
		}
		res, _, proof, err := cliCtx.Query("/custom/" + types.ModuleName + "/" + types.QueryValidatorDelegations, bz, isProve)
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
		value := delegations
		resp := rest.BuildQueryRes(height, isProve, value, proof)
		rest.PostProcessResponseBare(writer, cliCtx, resp)
	}
}

// HTTP request handler to query all delegator bonded validators
func delegatorValidatorsHandlerFn(cliCtx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		//vars := mux.Vars(request)
		//delegatorAddress := vars["delegatorAddr"]
		height := request.FormValue("height")
		prove := request.FormValue("prove")
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
		if !rest.CheckHeightAndProve(writer, height, prove, types.DefaultCodespace) {
			return
		}

		isProve := false
		if prove == "true" {
			isProve = true
		}
		res, _, proof, err := cliCtx.Query("/custom/" + types.ModuleName + "/" + types.QueryDelegatorValidators, bz, isProve)
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
		value := validators
		resp := rest.BuildQueryRes(height, isProve, value, proof)
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
		height := request.FormValue("height")
		prove := request.FormValue("prove")
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
		if !rest.CheckHeightAndProve(writer, height, prove, types.DefaultCodespace) {
			return
		}
		isProve := false
		if prove == "true" {
			isProve = true
		}
		res, _, proof, err := cliCtx.Query("/custom/" + types.ModuleName + "/" + types.QueryDelegatorValidator, bz, isProve)
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
		value := validator
		resp := rest.BuildQueryRes(height, isProve, value, proof)
		rest.PostProcessResponseBare(writer, cliCtx, resp)
	}
}


// HTTP request handler to query a delegation
func delegationHandlerFn(cliCtx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		//vars := mux.Vars(request)
		//validatorAddress := vars["validatorAddr"]
		//delegatorAddress := vars["delegatorAddr"]
		height := request.FormValue("height")
		prove := request.FormValue("prove")
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
		if !rest.CheckHeightAndProve(writer, height, prove, types.DefaultCodespace) {
			return
		}
		isProve := false
		if prove == "true" {
			isProve = true
		}
		res, _, proof, err := cliCtx.Query("/custom/" + types.ModuleName + "/" + types.QueryDelegation, bz, isProve)
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
		value := delegation
		resp := rest.BuildQueryRes(height, isProve, value, proof)
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
		height := r.FormValue("height")
		prove := r.FormValue("prove")
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
		if !rest.CheckHeightAndProve(w, height, prove, types.DefaultCodespace) {
			return
		}
		isProve := false
		if prove == "true" {
			isProve = true
		}
		res, _, proof, err := cliCtx.Query("/custom/" + types.ModuleName + "/" + types.QueryRedelegations, bz, isProve)
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
		value := redelegations
		resp := rest.BuildQueryRes(height, isProve, value, proof)
		rest.PostProcessResponseBare(w, cliCtx, resp)
	}
}