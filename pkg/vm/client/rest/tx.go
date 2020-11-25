package rest

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	"github.com/ci123chain/ci123chain/pkg/account/exported"
	keeper2 "github.com/ci123chain/ci123chain/pkg/account/keeper"
	"github.com/ci123chain/ci123chain/pkg/account/types"
	types2 "github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/client"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/util"
	evm "github.com/ci123chain/ci123chain/pkg/vm/evmtypes"
	"github.com/ci123chain/ci123chain/pkg/vm/keeper"
	vmmodule "github.com/ci123chain/ci123chain/pkg/vm/moduletypes"
	"github.com/ci123chain/ci123chain/pkg/vm/moduletypes/utils"
	wasmtypes "github.com/ci123chain/ci123chain/pkg/vm/wasmtypes"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"math/big"
	"net/http"
	"strconv"
	"strings"
)
const CAN_MIGRATE string = `{"method":"canMigrate()"}`

func uploadContractHandler(cliCtx context.Context, w http.ResponseWriter, r *http.Request) {
	broadcast, err := strconv.ParseBool(r.FormValue("broadcast"))
	if err != nil {
		broadcast = true
	}

	privKey, from, nonce, gas, err := rest.GetNecessaryParams(cliCtx, r, cdc, broadcast)
	if err != nil {
		rest.WriteErrorRes(w, wasmtypes.ErrCheckParams(vmmodule.DefaultCodespace, err.Error()))
		return
	}

	code, err := getCode(r)
	if err != nil || code == nil {
		rest.WriteErrorRes(w, wasmtypes.ErrCheckParams(vmmodule.DefaultCodespace, "get wasmcode failed"))
		return
	}

	var msg sdk.Msg
	if keeper.IsWasm(code) {
		wasmCode, err := keeper.UnCompress(code)
		if err != nil {
			rest.WriteErrorRes(w, sdk.ErrInternal("UnCompress code failed"))
			return
		}
		codeHash := keeper.MakeCodeHash(wasmCode)
		params := wasmtypes.NewQueryCodeInfoParams(string(codeHash))
		bz, Er := cliCtx.Cdc.MarshalJSON(params)
		if Er != nil {
			rest.WriteErrorRes(w, sdk.ErrInternal("marshal failed"))
			return
		}
		res, _, _, Err := cliCtx.Query("/custom/" + vmmodule.ModuleName + "/" + wasmtypes.QueryCodeInfo, bz, false)
		if Err != nil {
			rest.WriteErrorRes(w, Err)
			return
		}
		if len(res) > 0 { //already exists
			rest.PostProcessResponseBare(w, cliCtx, string(codeHash))
			return
		}
		msg = wasmtypes.NewMsgUploadContract(code, from)
	} else {
		amount_str := r.FormValue("amount")
		amount_int64, _ := strconv.ParseInt(amount_str, 10, 64)
		amount := big.NewInt(amount_int64)
		msg = evm.NewMsgEvmTx(from, nonce, nil, amount, gas, big.NewInt(1), code)
	}

	if !broadcast {
		rest.PostProcessResponseBare(w, cliCtx, hex.EncodeToString(msg.Bytes()))
		return
	}
	txByte, err := types2.SignCommonTx(from, nonce, gas, []sdk.Msg{msg}, privKey, cdc)
	if err != nil {
		rest.WriteErrorRes(w, wasmtypes.ErrCheckParams(vmmodule.DefaultCodespace,err.Error()))
		return
	}
	resp, err := cliCtx.BroadcastSignedTx(txByte)
	if err != nil {
		rest.WriteErrorRes(w, client.ErrBroadcast(vmmodule.DefaultCodespace, err))
		return
	}
	rest.PostProcessResponseBare(w, cliCtx, resp)
}

func instantiateContractHandler(cliCtx context.Context,w http.ResponseWriter, r *http.Request) {
	broadcast, err := strconv.ParseBool(r.FormValue("broadcast"))
	if err != nil {
		broadcast = true
	}
	privKey, from, nonce, gas, err := rest.GetNecessaryParams(cliCtx, r, cdc, broadcast)
	if err != nil {
		rest.WriteErrorRes(w, wasmtypes.ErrCheckParams(vmmodule.DefaultCodespace, err.Error()))
		return
	}
	codeHash := r.FormValue("code_hash")
	hash, err := hex.DecodeString(strings.ToLower(codeHash))
	if err != nil {
		rest.WriteErrorRes(w, wasmtypes.ErrCheckParams(vmmodule.DefaultCodespace, err.Error()))
	}
	name, version, author, email, describe, err := adjustInstantiateParams(r)
	if err != nil {
		rest.WriteErrorRes(w, wasmtypes.ErrCheckParams(vmmodule.DefaultCodespace, "get params failed"))
		return
	}
	var args utils.CallData
	args_str := r.FormValue("args")
	if args_str == "" {
		rest.WriteErrorRes(w, wasmtypes.ErrCheckParams(vmmodule.DefaultCodespace, "get callData failed"))
		return
	}else {
		ok, err := util.CheckJsonArgs(args_str, args)
		if err != nil || !ok {
			rest.WriteErrorRes(w, wasmtypes.ErrCheckParams(vmmodule.DefaultCodespace, "get callData failed"))
			return
		}
	}
	msg := wasmtypes.NewMsgInstantiateContract(hash, from, name, version, author, email, describe, args)
	if !broadcast {
		rest.PostProcessResponseBare(w, cliCtx, hex.EncodeToString(msg.Bytes()))
		return
	}

	txByte, err := types2.SignCommonTx(from, nonce, gas, []sdk.Msg{msg}, privKey, cdc)
	if err != nil {
		rest.WriteErrorRes(w, wasmtypes.ErrCheckParams(vmmodule.DefaultCodespace,err.Error()))
		return
	}
	res, err := cliCtx.BroadcastSignedTx(txByte)
	if err != nil {
		rest.WriteErrorRes(w, client.ErrBroadcast(vmmodule.DefaultCodespace, err))
		return
	}
	rest.PostProcessResponseBare(w, cliCtx, res)
}

func executeContractHandler(cliCtx context.Context,w http.ResponseWriter, r *http.Request) {
	broadcast, err := strconv.ParseBool(r.FormValue("broadcast"))
	if err != nil {
		broadcast = true
	}
	privKey, from, nonce, gas, err := rest.GetNecessaryParams(cliCtx, r, cdc, broadcast)
	if err != nil {
		rest.WriteErrorRes(w, wasmtypes.ErrCheckParams(vmmodule.DefaultCodespace, err.Error()))
		return
	}

	contractAddr := r.FormValue("contract_address")
	contractAddress := sdk.HexToAddress(contractAddr)
	qparams := keeper2.NewQueryAccountParams(contractAddress)
	bz, err := cliCtx.Cdc.MarshalJSON(qparams)
	if err != nil {
		rest.WriteErrorRes(w, wasmtypes.ErrCheckParams(vmmodule.DefaultCodespace, err.Error()))
		return
	}

	queryRes, _, _, err := cliCtx.Query("/custom/" + types.ModuleName + "/" + types.QueryAccount, bz, false)
	if err != nil {
		rest.WriteErrorRes(w, wasmtypes.ErrCheckParams(vmmodule.DefaultCodespace,"query contract account failed"))
		return
	}
	if queryRes == nil{
		rest.WriteErrorRes(w, wasmtypes.ErrCheckParams(vmmodule.DefaultCodespace,"contract account does not exist"))
		return
	}
	var acc exported.Account
	err2 := cliCtx.Cdc.UnmarshalBinaryLengthPrefixed(queryRes, &acc)
	if err2 != nil {
		rest.WriteErrorRes(w, sdk.ErrInternal("unmarshal query response to account failed"))
		return
	}
	var msg sdk.Msg
	if acc.GetContractType() == types.WasmContractType {
		var args utils.CallData
		args_str := r.FormValue("args")
		if args_str == "" {
			rest.WriteErrorRes(w, wasmtypes.ErrCheckParams(vmmodule.DefaultCodespace, "get callData failed"))
			return
		}else {
			ok, err := util.CheckJsonArgs(args_str, args)
			if err != nil || !ok {
				rest.WriteErrorRes(w, wasmtypes.ErrCheckParams(vmmodule.DefaultCodespace, "get callData failed"))
				return
			}
		}
		msg = wasmtypes.NewMsgExecuteContract(from, contractAddress, args)
	} else if acc.GetContractType() == types.EvmContractType {
		var to *ethcmn.Address
		to_str := r.FormValue("contract_address")
		if len(to_str) == 0 {
			rest.WriteErrorRes(w, wasmtypes.ErrCheckParams(vmmodule.DefaultCodespace, "contract_address cannot be empty"))
			return
		} else {
			to_addr := ethcmn.HexToAddress(to_str)
			to = &to_addr
		}
		amount_str := r.FormValue("amount")
		amount_int64, _ := strconv.ParseInt(amount_str, 10, 64)
		amount := big.NewInt(amount_int64)

		var args utils.CallData
		args_str := r.FormValue("args")
		if args_str == "" {
			rest.WriteErrorRes(w, wasmtypes.ErrCheckParams(vmmodule.DefaultCodespace, "get callData failed"))
			return
		}else {
			ok, err := util.CheckJsonArgs(args_str, args)
			if err != nil || !ok {
				rest.WriteErrorRes(w, wasmtypes.ErrCheckParams(vmmodule.DefaultCodespace, "get callData failed"))
				return
			}
		}
		payload, err := evm.EVMEncode(args)
		s := hex.EncodeToString(payload)
		fmt.Println(s)
		if err != nil {
			rest.WriteErrorRes(w, wasmtypes.ErrCheckParams(vmmodule.DefaultCodespace, "encode evm callData failed"))
			return
		}
		msg = evm.NewMsgEvmTx(from, nonce, to, amount, gas, big.NewInt(1), payload)
	} else {
		rest.WriteErrorRes(w, wasmtypes.ErrCheckParams(vmmodule.DefaultCodespace,"not contract account"))
		return
	}

	if !broadcast {
		rest.PostProcessResponseBare(w, cliCtx, hex.EncodeToString(msg.Bytes()))
		return
	}
	txByte, err := types2.SignCommonTx(from, nonce, gas, []sdk.Msg{msg}, privKey, cdc)
	if err != nil {
		rest.WriteErrorRes(w, wasmtypes.ErrCheckParams(vmmodule.DefaultCodespace,err.Error()))
		return
	}
	res, err := cliCtx.BroadcastSignedTx(txByte)
	if err != nil {
		rest.WriteErrorRes(w, client.ErrBroadcast(vmmodule.DefaultCodespace, err))
		return
	}
	rest.PostProcessResponseBare(w, cliCtx, res)
}

func migrateContractHandler(cliCtx context.Context,w http.ResponseWriter, r *http.Request) {
	broadcast, err := strconv.ParseBool(r.FormValue("broadcast"))
	if err != nil {
		broadcast = true
	}
	privKey, from, nonce, gas, err := rest.GetNecessaryParams(cliCtx, r, cdc, broadcast)
	if err != nil {
		rest.WriteErrorRes(w, wasmtypes.ErrCheckParams(vmmodule.DefaultCodespace, err.Error()))
		return
	}

	sender := cliCtx.FromAddr
	contractAddr := r.FormValue("contract_address")
	contractAddress := sdk.HexToAddress(contractAddr)
	var arg utils.CallData
	queryParam := []byte(CAN_MIGRATE)
	err = json.Unmarshal(queryParam, &arg)
	if err != nil {
		rest.WriteErrorRes(w, wasmtypes.ErrCheckParams(vmmodule.DefaultCodespace, err.Error()))
	}
	params := wasmtypes.NewContractStateParam(contractAddress, sender, arg)
	bz, Er := cliCtx.Cdc.MarshalJSON(params)
	if Er != nil {
		rest.WriteErrorRes(w, sdk.ErrInternal("marshal failed"))
		return
	}

	resQuery, _, _, _ := cliCtx.Query("/custom/" + vmmodule.ModuleName + "/" + wasmtypes.QueryContractState, bz, false)
	var contractState wasmtypes.ContractState
	cliCtx.Cdc.MustUnmarshalJSON(resQuery, &contractState)
	if contractState.Result != "true" {
		rest.WriteErrorRes(w, sdk.ErrInternal("No permissions to migrate contracts"))
		return
	}

	codeHash := r.FormValue("code_hash")
	hash, err := hex.DecodeString(strings.ToLower(codeHash))
	if err != nil {
		rest.WriteErrorRes(w, wasmtypes.ErrCheckParams(vmmodule.DefaultCodespace, "codeHash error"))
		return
	}
	name, version, author, email, describe, err := adjustInstantiateParams(r)
	if err != nil {
		rest.WriteErrorRes(w, wasmtypes.ErrCheckParams(vmmodule.DefaultCodespace, "get params failed"))
		return
	}

	var args utils.CallData
	args_str := r.FormValue("args")
	if args_str == "" {
		rest.WriteErrorRes(w, wasmtypes.ErrCheckParams(vmmodule.DefaultCodespace, "get callData failed"))
		return
	}else {
		ok, err := util.CheckJsonArgs(args_str, args)
		if err != nil || !ok {
			rest.WriteErrorRes(w, wasmtypes.ErrCheckParams(vmmodule.DefaultCodespace, "get callData failed"))
			return
		}
	}
	msg := wasmtypes.NewMsgMigrateContract(hash, from, name, version, author, email, describe, contractAddress, args)
	if !broadcast {
		rest.PostProcessResponseBare(w, cliCtx, hex.EncodeToString(msg.Bytes()))
		return
	}

	txByte, err := types2.SignCommonTx(from, nonce, gas, []sdk.Msg{msg}, privKey, cdc)
	if err != nil {
		rest.WriteErrorRes(w, wasmtypes.ErrCheckParams(vmmodule.DefaultCodespace,err.Error()))
		return
	}
	res, err := cliCtx.BroadcastSignedTx(txByte)
	if err != nil {
		rest.WriteErrorRes(w, client.ErrBroadcast(vmmodule.DefaultCodespace, err))
		return
	}
	rest.PostProcessResponseBare(w, cliCtx, res)
}