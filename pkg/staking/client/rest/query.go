package rest

import (
	"encoding/hex"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/staking/types"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

func RegisterQueryRoutes(cliCtx context.Context, r *mux.Router) {
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

	r.HandleFunc("/staking/operator", operatorAddressSetQueryHandleFn(cliCtx), ).Methods("POST")

}

func delegatorDelegationsHandlerFn(cliCtx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		//vars := mux.Vars(request)
		//delegatorAddress := vars["delegatorAddr"]
		delegatorAddress := request.FormValue("delegator_address")
		height := request.FormValue("height")
		prove := request.FormValue("prove")
		delegatorAddr := sdk.HexToAddress(delegatorAddress)
		params := types.NewQueryDelegatorParams(delegatorAddr)
		bz, Err := cliCtx.Cdc.MarshalJSON(params)
		if Err != nil {
			rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("cdc marshal failed: %v", Err.Error())).Error())
			return
		}
		cliCtx, ok, err := rest.ParseQueryHeightOrReturnBadRequest(writer, cliCtx, request, "")
		if !ok || err != nil {
			rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInternal, "get clictx failed").Error())
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
			rest.WriteErrorRes(writer, err.Error())
			return
		}
		if len(res) < 1 {
			rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrResponse, fmt.Sprintf("unexpected res: %v", res)).Error())
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
			rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInternal, "get clictx failed").Error())
			return
		}
		if !rest.CheckHeightAndProve(writer, height, prove, types.DefaultCodespace) {
			return
		}

		_, page, limit, Err := rest.ParseHTTPArgsWithLimit(request, 0)
		if err != nil {
			rest.WriteErrorRes(writer, err.Error())
			return
		}

		status := request.FormValue("status")
		if status == "" {
			status = sdk.BondStatusAll
		}

		params := types.NewQueryValidatorsParams(page, limit, status)
		bz, Err := cliCtx.Cdc.MarshalJSON(params)
		if Err != nil {
			rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("cdc marshal failed: %v", Err.Error())).Error())
			return
		}

		isProve := false
		if prove == "true" {
			isProve = true
		}
		res, _, proof, err := cliCtx.Query("/custom/" + types.ModuleName + "/" + types.QueryValidators, bz, isProve)
		if err != nil {
			rest.WriteErrorRes(writer, err.Error())
			return
		}
		if len(res) < 1 {
			rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrResponse, fmt.Sprintf("get unexpected res: %v", res)).Error())
			return
		}
		var validators []types.Validator
		cliCtx.Cdc.MustUnmarshalJSON(res, &validators)
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
		if !ok || err != nil {
			rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("get clictx failed")).Error())
			return
		}
		if !rest.CheckHeightAndProve(writer, height, prove, types.DefaultCodespace) {
			return
		}

		validatorAddress := request.FormValue("validator_address")
		validatorAddr := sdk.HexToAddress(validatorAddress)
		params := types.NewQueryValidatorParams(validatorAddr)
		bz, Err := cliCtx.Cdc.MarshalJSON(params)
		if Err != nil {
			rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("cdc marshal faield: %v", Err.Error())).Error())
			return
		}

		isProve := false
		if prove == "true" {
			isProve = true
		}
		res, _, proof, err := cliCtx.Query("/custom/" + types.ModuleName + "/" + types.QueryValidator, bz, isProve)
		if err != nil {
			rest.WriteErrorRes(writer, err.Error())
			return
		}
		if len(res) < 1 {
			rest.WriteErrorRes(writer,sdkerrors.Wrap(sdkerrors.ErrResponse, fmt.Sprintf("get unexpected res: %v", res)).Error())
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
		validatorAddress := request.FormValue("validator_address")
		validatorAddr := sdk.HexToAddress(validatorAddress)
		params := types.NewQueryValidatorParams(validatorAddr)
		bz, Err := cliCtx.Cdc.MarshalJSON(params)
		if Err != nil {
			rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("cdc marhsal failed: %v", Err.Error())).Error())
			return
		}

		cliCtx, ok, err := rest.ParseQueryHeightOrReturnBadRequest(writer, cliCtx, request, "")
		if !ok || err != nil {
			rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInternal, "get clictx failed").Error())
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
			rest.WriteErrorRes(writer, err.Error())
			return
		}
		if len(res) < 1 {
			rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrResponse, fmt.Sprintf("get unexpected res: %v", res)).Error())
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
		delegatorAddress := request.FormValue("delegator_address")
		delegatorAddr := sdk.HexToAddress(delegatorAddress)
		params := types.NewQueryDelegatorParams(delegatorAddr)
		bz, Err := cliCtx.Cdc.MarshalJSON(params)
		if Err != nil {
			rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("cdc marhsal failed: %v", Err.Error())).Error())
			return
		}

		cliCtx, ok, err := rest.ParseQueryHeightOrReturnBadRequest(writer, cliCtx, request, "")
		if !ok || err != nil {
			rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInternal, "get clictx failed").Error())
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
			rest.WriteErrorRes(writer, err.Error())
			return
		}
		if len(res) < 1 {
			rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrResponse, fmt.Sprintf("get unexpected res: %v", res)).Error())
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
		validatorAddress := request.FormValue("validator_address")
		delegatorAddress := request.FormValue("delegator_address")
		height := request.FormValue("height")
		prove := request.FormValue("prove")
		validatorAddr := sdk.HexToAddress(validatorAddress)
		delegatorAddr := sdk.HexToAddress(delegatorAddress)
		params := types.NewQueryBondsParams(delegatorAddr, validatorAddr)
		bz, Err := cliCtx.Cdc.MarshalJSON(params)
		if Err != nil {
			rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("cdc marhsal failed: %v", Err.Error())).Error())
			return
		}

		cliCtx, ok, err := rest.ParseQueryHeightOrReturnBadRequest(writer, cliCtx, request, "")
		if !ok || err != nil {
			rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInternal, "get clictx failed").Error())
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
			rest.WriteErrorRes(writer, err.Error())
			return
		}
		if len(res) < 1 {
			rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrResponse, fmt.Sprintf("get unexpected res: %v", res)).Error())
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

		height := request.FormValue("height")
		prove := request.FormValue("prove")
		validatorAddress := request.FormValue("validator_address")
		delegatorAddress := request.FormValue("delegator_address")
		validatorAddr := sdk.HexToAddress(validatorAddress)
		delegatorAddr := sdk.HexToAddress(delegatorAddress)
		params := types.NewQueryBondsParams(delegatorAddr, validatorAddr)
		bz, Err := cliCtx.Cdc.MarshalJSON(params)
		if Err != nil {
			rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("cdc marhsal failed: %v", Err.Error())).Error())
			return
		}

		cliCtx, ok, err := rest.ParseQueryHeightOrReturnBadRequest(writer, cliCtx, request, "")
		if !ok || err != nil {
			rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInternal, "get clictx failed").Error())
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
			rest.WriteErrorRes(writer, err.Error())
			return
		}
		if len(res) < 1 {
			rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrResponse, fmt.Sprintf("get unexpected res: %v", res)).Error())
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
		if !ok || err != nil {
			rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrInternal, "get clictx failed").Error())
			return
		}
		height := r.FormValue("height")
		prove := r.FormValue("prove")
		delegatorAddress := r.FormValue("delegator_address")
		validatorSrc := r.FormValue("validator_src_address")
		validatorDst := r.FormValue("validator_dst_address")

		delegatorAddr := sdk.HexToAddress(delegatorAddress)
		validatorSrcAddr := sdk.HexToAddress(validatorSrc)
		validatorDstAddr := sdk.HexToAddress(validatorDst)

		params.DelegatorAddr = delegatorAddr
		params.SrcValidatorAddr = validatorSrcAddr
		params.DstValidatorAddr = validatorDstAddr
		bz, Err := cliCtx.Cdc.MarshalJSON(params)
		if Err != nil {
			rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("cdc marhsal failed: %v", Err.Error())).Error())
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
			rest.WriteErrorRes(w, err.Error())
			return
		}
		if len(res) < 1 {
			rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrResponse, fmt.Sprintf("get unexpected res: %v", res)).Error())
			return
		}
		var redelegations types.RedelegationResponses
		cliCtx.Cdc.MustUnmarshalJSON(res, &redelegations)
		value := redelegations
		resp := rest.BuildQueryRes(height, isProve, value, proof)
		rest.PostProcessResponseBare(w, cliCtx, resp)
	}
}

func operatorAddressSetQueryHandleFn(cliCtx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//
		var params types.QueryOperatorAddressesParams
		cliCtx, ok, err := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r, "")
		if !ok || err != nil {
			rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrInternal, "get clictx failed").Error())
			return
		}
		params.ConsAddresses = make([]sdk.AccAddr, 0)
		height := r.FormValue("height")
		prove := r.FormValue("prove")
		setsStr := r.FormValue("address_sets")
		sets := strings.Split(setsStr, ",")
		for _, v := range sets {
			b, err := hex.DecodeString(v)
			if err != nil {
				rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrParams, err.Error()).Error())
				return
			}
			params.ConsAddresses = append(params.ConsAddresses, sdk.AccAddr(b))
		}
		bz, Err := cliCtx.Cdc.MarshalJSON(params)
		if Err != nil {
			rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("cdc marhsal failed: %v", Err.Error())).Error())
			return
		}
		if !rest.CheckHeightAndProve(w, height, prove, types.DefaultCodespace) {
			return
		}
		isProve := false
		if prove == "true" {
			isProve = true
		}
		res, _, proof, err := cliCtx.Query("/custom/" + types.ModuleName + "/" + types.QueryOperatorAddressSet, bz, isProve)
		if err != nil {
			rest.WriteErrorRes(w, err.Error())
			return
		}
		if len(res) < 1 {
			rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrResponse, fmt.Sprintf("get unexpected res: %v", res)).Error())
			return
		}
		var result []types.ValidatorOperatorAddressResponse
		cliCtx.Cdc.MustUnmarshalJSON(res, &result)
		value := result
		resp := rest.BuildQueryRes(height, isProve, value, proof)
		rest.PostProcessResponseBare(w, cliCtx, resp)
	}
}