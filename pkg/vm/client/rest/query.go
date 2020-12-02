package rest

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	"github.com/ci123chain/ci123chain/pkg/account/exported"
	"github.com/ci123chain/ci123chain/pkg/account/keeper"
	types2 "github.com/ci123chain/ci123chain/pkg/account/types"
	types3 "github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/client"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/transfer"
	"github.com/ci123chain/ci123chain/pkg/vm"
	evm "github.com/ci123chain/ci123chain/pkg/vm/evmtypes"
	"github.com/ci123chain/ci123chain/pkg/vm/moduletypes"
	"github.com/ci123chain/ci123chain/pkg/vm/moduletypes/utils"
	"github.com/ci123chain/ci123chain/pkg/vm/wasmtypes"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"math/big"
	"net/http"
	"strconv"
)

const DefaultRPCGasLimit = 10000000

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
		contractAddr := r.FormValue("contract_address")
		contractAddress := sdk.HexToAddress(contractAddr)
		qparams := keeper.NewQueryAccountParams(contractAddress)
		bz, err := cliCtx.Cdc.MarshalJSON(qparams)
		if err != nil {
			rest.WriteErrorRes(w, types.ErrCheckParams(vm.DefaultCodespace, err.Error()))
			return
		}
		from := r.FormValue("from")
		fromAddr := sdk.HexToAddress(from)
		ctx, err2 := client.NewClientContextFromViper(cdc)
		if err2 != nil {
			rest.WriteErrorRes(w, types.ErrCheckParams(vm.DefaultCodespace, err2.Error()))
			return
		}
		nonce, _, err := ctx.GetNonceByAddress(fromAddr, false)
		if err != nil {
			return
		}
		var gas uint64
		gasStr := r.FormValue("gas")
		if gasStr != "" {
			gas, err = strconv.ParseUint(gasStr, 10, 64)
			if err != nil {
				rest.WriteErrorRes(w, types.ErrCheckParams(vm.DefaultCodespace, err.Error()))
				return
			}
		} else {
			gas = uint64(DefaultRPCGasLimit)
		}

		var msg sdk.Msg
		var args utils.CallData
		args_str := r.FormValue("calldata")
		if args_str == "" {
			rest.WriteErrorRes(w, types.ErrCheckParams(vm.DefaultCodespace, "get callData failed"))
			return
		}else {
			err := json.Unmarshal([]byte(args_str), &args)
			if err != nil  {
				rest.WriteErrorRes(w, types.ErrCheckParams(vm.DefaultCodespace, "get callData failed"))
				return
			}
		}
		queryRes, _, _, err := cliCtx.Query("/custom/" + types2.ModuleName + "/" + types2.QueryAccount, bz, false)
		if err != nil {
			rest.WriteErrorRes(w, types.ErrCheckParams(vm.DefaultCodespace,"query contract account failed"))
			return
		}
		if queryRes == nil{
			rest.WriteErrorRes(w, types.ErrCheckParams(vm.DefaultCodespace,"contract account does not exist"))
			return
		}
		var acc exported.Account
		err2 = cliCtx.Cdc.UnmarshalBinaryLengthPrefixed(queryRes, &acc)
		if err2 != nil {
			rest.WriteErrorRes(w, sdk.ErrInternal("unmarshal query response to account failed"))
			return
		}
		if acc.GetContractType() == types2.WasmContractType {
			msg = types.NewMsgExecuteContract(fromAddr, contractAddress, args)
		} else if acc.GetContractType() == types2.EvmContractType {
			var to *ethcmn.Address
			to_addr := ethcmn.HexToAddress(contractAddr)
			to = &to_addr
			amount_str := r.FormValue("amount")
			amount := new(big.Int)
			if len(amount_str) < 2 {
				amount.SetString(amount_str, 10)
			} else if amount_str[:2] == "0x" {
				amount.SetString(amount_str[2:], 16)
			} else {
				amount.SetString(amount_str, 10)
			}

			payload, err := evm.EVMEncode(args)
			s := hex.EncodeToString(payload)
			fmt.Println(s)
			if err != nil {
				rest.WriteErrorRes(w, types.ErrCheckParams(vm.DefaultCodespace, "encode evm callData failed"))
				return
			}
			msg = evm.NewMsgEvmTx(fromAddr, nonce, to, amount, gas, big.NewInt(1), payload)
		} else {
			rest.WriteErrorRes(w, types.ErrCheckParams(vm.DefaultCodespace,"not contract account"))
			return
		}

		txBytes, err := cdc.MarshalBinaryBare(types3.NewCommonTx(fromAddr, nonce, gas, []sdk.Msg{msg}))
		if err != nil {
			rest.WriteErrorRes(w, types.ErrCheckParams(vm.DefaultCodespace,err.Error()))
			return
		}

		res, _, _, err := cliCtx.Query("app/simulate", txBytes, false)
		if err != nil {
			rest.WriteErrorRes(w, transfer.ErrQueryTx(types.DefaultCodespace, err.Error()))
			return
		}
		var simResponse sdk.QureyAppResponse
		if err := cdc.UnmarshalBinaryBare(res, &simResponse); err != nil {
			rest.WriteErrorRes(w, transfer.ErrQueryTx(types.DefaultCodespace, err.Error()))
			return
		}

		rest.PostProcessResponseBare(w, cliCtx, simResponse)
	}
}
