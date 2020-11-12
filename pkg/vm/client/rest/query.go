package rest

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/transfer"
	"github.com/ci123chain/ci123chain/pkg/util"
	"github.com/ci123chain/ci123chain/pkg/vm/moduletypes"
	"github.com/ci123chain/ci123chain/pkg/vm/wasmtypes"
	"net/http"
)

//func listCodesHandlerFn(cliCtx context.Context) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		rest.WriteErrorRes(w, sdk.ErrInternal("Implement me"))
//	}
//}


func queryCodeHandlerFn(cliCtx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		height := r.FormValue("height")
		prove := r.FormValue("prove")
		codeHash := r.FormValue("code_hash")
		cliCtx, ok, Err := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r, "")
		if !ok {
			rest.WriteErrorRes(w, Err)
			return
		}
		if !rest.CheckHeightAndProve(w, height, prove, moduletypes.DefaultCodespace) {
			return
		}
		params := types.NewQueryCodeInfoParams(codeHash)
		bz, Er := cliCtx.Cdc.MarshalJSON(params)
		if Er != nil {
			rest.WriteErrorRes(w, sdk.ErrInternal("marshal failed"))
			return
		}

		isProve := false
		if prove == "true" {
			isProve = true
		}
		res, _, proof, Err := cliCtx.Query("/custom/" + moduletypes.ModuleName + "/" + types.QueryCodeInfo, bz, isProve)
		if Err != nil {
			rest.WriteErrorRes(w, Err)
			return
		}
		if len(res) < 1 {
			rest.WriteErrorRes(w, transfer.ErrQueryTx(types.DefaultCodespace, "no excepted code"))
			return
		}
		var codeInfo types.CodeInfo
		cliCtx.Cdc.MustUnmarshalBinaryBare(res, &codeInfo)
		value := codeInfo
		resp := rest.BuildQueryRes(height, isProve, value, proof)
		rest.PostProcessResponseBare(w, cliCtx, resp)
	}
}


func listContractsByCodeHandlerFn(cliCtx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		height := r.FormValue("height")
		prove := r.FormValue("prove")
		accountAddr := r.FormValue("account_address")
		accountAddress := sdk.HexToAddress(accountAddr)
		params := types.NewContractListParams(accountAddress)
		bz, Er := cliCtx.Cdc.MarshalJSON(params)
		if Er != nil {
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
		res, _, _, err := cliCtx.Query("/custom/" + moduletypes.ModuleName + "/" + types.QueryContractList, bz, isProve)
		if err != nil {
			rest.WriteErrorRes(w, err)
			return
		}
		if len(res) < 1 {
			rest.WriteErrorRes(w, transfer.ErrQueryTx(types.DefaultCodespace, "the length of contract list is 0"))
			return
		}
		var contractList types.ContractListResponse
		cliCtx.Cdc.MustUnmarshalBinaryBare(res, &contractList)

		rest.PostProcessResponseBare(w, cliCtx, contractList)
	}
}

func queryContractHandlerFn(cliCtx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		height := r.FormValue("height")
		prove := r.FormValue("prove")
		contractAddr := r.FormValue("contract_address")
		contractAddress := sdk.HexToAddress(contractAddr)

		cliCtx, ok, err := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r, "")
		if !ok {
			rest.WriteErrorRes(w, err)
			return
		}
		if !rest.CheckHeightAndProve(w, height, prove, types.DefaultCodespace) {
			return
		}

		params := types.NewQueryContractInfoParams(contractAddress)
		bz, Er := cliCtx.Cdc.MarshalJSON(params)
		if Er != nil {
			rest.WriteErrorRes(w, sdk.ErrInternal("marshal failed"))
			return
		}

		isProve := false
		if prove == "true" {
			isProve = true
		}
		res, _, proof, err := cliCtx.Query("/custom/" + moduletypes.ModuleName + "/" + types.QueryContractInfo, bz, isProve)
		if err != nil {
			rest.WriteErrorRes(w, err)
			return
		}
		if len(res) < 1 {
			rest.WriteErrorRes(w, transfer.ErrQueryTx(types.DefaultCodespace, "no expected contract"))
			return
		}
		var contractInfo types.ContractInfo
		cliCtx.Cdc.MustUnmarshalBinaryBare(res, &contractInfo)
		value := contractInfo
		resp := rest.BuildQueryRes(height, isProve, value, proof)
		rest.PostProcessResponseBare(w, cliCtx, resp)
	}
}

func queryContractStateAllHandlerFn(cliCtx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var queryParam []byte
		contractAddr := r.FormValue("contract_address")
		contractAddress := sdk.HexToAddress(contractAddr)
		msg := r.FormValue("args")
		height := r.FormValue("height")
		prove := r.FormValue("prove")
		if msg == "" {
			queryParam = nil
		}else {
			var argsStr types.CallContractParam
			ok, err := util.CheckJsonArgs(msg, argsStr)
			if err != nil || !ok {
				//return types.AccAddress{}, 0, 0, "", nil, errors.New("unexpected args")
			}
			queryParam = []byte(msg)

		}

		cliCtx, ok, err := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r, "")
		if !ok {
			rest.WriteErrorRes(w, err)
			return
		}
		if !rest.CheckHeightAndProve(w, height, prove, types.DefaultCodespace) {
			return
		}
		params := types.NewContractStateParam(contractAddress, sdk.AccAddress{}, queryParam)
		bz, Er := cliCtx.Cdc.MarshalJSON(params)
		if Er != nil {
			rest.WriteErrorRes(w, sdk.ErrInternal("marshal failed"))
			return
		}

		isProve := false
		if prove == "true" {
			isProve = true
		}
		res, _, proof, err := cliCtx.Query("/custom/" + moduletypes.ModuleName + "/" + types.QueryContractState, bz, isProve)
		if err != nil {
			rest.WriteErrorRes(w, err)
			return
		}
		if len(res) < 1 {
			rest.WriteErrorRes(w, transfer.ErrQueryTx(types.DefaultCodespace, "query response length less than 1"))
			return
		}
		var contractState types.ContractState
		cliCtx.Cdc.MustUnmarshalJSON(res, &contractState)
		value := contractState
		resp := rest.BuildQueryRes(height, isProve, value, proof)
		rest.PostProcessResponseBare(w, cliCtx, resp)
	}
}
