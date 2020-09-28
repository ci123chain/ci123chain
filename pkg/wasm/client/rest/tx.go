package rest

import (
	"encoding/hex"
	"encoding/json"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	types2 "github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/client"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/util"
	wasm2 "github.com/ci123chain/ci123chain/pkg/wasm"
	"github.com/ci123chain/ci123chain/pkg/wasm/types"
	"net/http"
	"strconv"
	"strings"
)

const CAN_MIGRATE string = `{"method":"canMigrate","args": [""]}`
func uploadContractHandler(cliCtx context.Context,w http.ResponseWriter, r *http.Request) {
	broadcast, err := strconv.ParseBool(r.FormValue("broadcast"))
	if err != nil {
		broadcast = true
	}
	privKey, from, nonce, gas, err := rest.GetNecessaryParams(cliCtx, r, cdc, broadcast)
	if err != nil {
		rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace, err.Error()))
		return
	}
	wasmCode, err := getWasmCode(r)
	if err != nil || wasmCode == nil {
		rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace, "get wasmcode failed"))
		return
	}
	msg := wasm2.NewUploadTx(wasmCode, from)
	if !broadcast {
		rest.PostProcessResponseBare(w, cliCtx, hex.EncodeToString(msg.Bytes()))
		return
	}

	txByte, err := types2.SignCommonTx(from, nonce, gas, []sdk.Msg{msg}, privKey, cdc)
	if err != nil {
		rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace,err.Error()))
		return
	}
	res, err := cliCtx.BroadcastSignedTx(txByte)
	if err != nil {
		rest.WriteErrorRes(w, client.ErrBroadcast(types.DefaultCodespace, err))
		return
	}
	rest.PostProcessResponseBare(w, cliCtx, res)
}

func instantiateContractHandler(cliCtx context.Context,w http.ResponseWriter, r *http.Request) {
	broadcast, err := strconv.ParseBool(r.FormValue("broadcast"))
	if err != nil {
		broadcast = true
	}
	privKey, from, nonce, gas, err := rest.GetNecessaryParams(cliCtx, r, cdc, broadcast)
	if err != nil {
		rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace, err.Error()))
		return
	}
	codeHash := r.FormValue("code_hash")
	hash, err := hex.DecodeString(strings.ToLower(codeHash))
	if err != nil {
		rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace, err.Error()))
	}
	name, version, author, email, describe, err := adjustInstantiateParams(r)
	if err != nil {
		rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace, "get params failed"))
		return
	}
	var args []byte
	args_str := r.FormValue("args")
	if args_str == "" {
		args = nil
	}else {
		var argsStr types.CallContractParam
		ok, err := util.CheckJsonArgs(args_str, argsStr)
		if err != nil || !ok {
			rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace, "get initArgs failed"))
			return
		}
		var argsByte = []byte(args_str)
		args = argsByte
	}
	JsonArgs := json.RawMessage(args)
	msg := wasm2.NewInstantiateTx(hash, from, name, version, author, email, describe, JsonArgs)
	if !broadcast {
		rest.PostProcessResponseBare(w, cliCtx, hex.EncodeToString(msg.Bytes()))
		return
	}

	txByte, err := types2.SignCommonTx(from, nonce, gas, []sdk.Msg{msg}, privKey, cdc)
	if err != nil {
		rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace,err.Error()))
		return
	}
	res, err := cliCtx.BroadcastSignedTx(txByte)
	if err != nil {
		rest.WriteErrorRes(w, client.ErrBroadcast(types.DefaultCodespace, err))
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
		rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace, err.Error()))
		return
	}

	contractAddr := r.FormValue("contract_address")
	contractAddress := sdk.HexToAddress(contractAddr)
	params := types.NewQueryContractInfoParams(contractAddress)
	bz, Er := cliCtx.Cdc.MarshalJSON(params)
	if Er != nil {
		rest.WriteErrorRes(w, sdk.ErrInternal("marshal failed"))
		return
	}

	resQuery, _, _, _:= cliCtx.Query("/custom/" + types.ModuleName + "/" + types.QueryContractInfo, bz, false)
	if resQuery == nil {
		rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace, "contract does not exist or get contract error"))
		return
	}
	var args []byte
	args_str := r.FormValue("args")
	if args_str == "" {
		args = nil
	}else {
		var argsStr types.CallContractParam
		ok, err := util.CheckJsonArgs(args_str, argsStr)
		if err != nil || !ok {
			rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace, "get initArgs failed"))
			return
		}
		var argsByte = []byte(args_str)
		args = argsByte
	}
	JsonArgs := json.RawMessage(args)
	msg := types.NewMsgExecuteContract(from, contractAddress, JsonArgs)
	if !broadcast {
		rest.PostProcessResponseBare(w, cliCtx, hex.EncodeToString(msg.Bytes()))
		return
	}

	txByte, err := types2.SignCommonTx(from, nonce, gas, []sdk.Msg{msg}, privKey, cdc)
	if err != nil {
		rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace,err.Error()))
		return
	}
	res, err := cliCtx.BroadcastSignedTx(txByte)
	if err != nil {
		rest.WriteErrorRes(w, client.ErrBroadcast(types.DefaultCodespace, err))
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
		rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace, err.Error()))
		return
	}

	sender := cliCtx.FromAddr
	contractAddr := r.FormValue("contract_address")
	contractAddress := sdk.HexToAddress(contractAddr)
	queryParam := []byte(CAN_MIGRATE)
	params := types.NewContractStateParam(contractAddress, sender, queryParam)
	bz, Er := cliCtx.Cdc.MarshalJSON(params)
	if Er != nil {
		rest.WriteErrorRes(w, sdk.ErrInternal("marshal failed"))
		return
	}

	resQuery, _, _, _ := cliCtx.Query("/custom/" + types.ModuleName + "/" + types.QueryContractState, bz, false)
	var contractState types.ContractState
	cliCtx.Cdc.MustUnmarshalJSON(resQuery, &contractState)
	if contractState.Result != "true" {
		rest.WriteErrorRes(w, sdk.ErrInternal("No permissions to migrate contracts"))
		return
	}

	codeHash := r.FormValue("code_hash")
	hash, err := hex.DecodeString(strings.ToLower(codeHash))
	if err != nil {
		rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace, "codeHash error"))
		return
	}
	name, version, author, email, describe, err := adjustInstantiateParams(r)
	if err != nil {
		rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace, "get params failed"))
		return
	}

	var args []byte
	args_str := r.FormValue("args")
	if args_str == "" {
		args = nil
	}else {
		var argsStr types.CallContractParam
		ok, err := util.CheckJsonArgs(args_str, argsStr)
		if err != nil || !ok {
			rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace, "get initArgs failed"))
			return
		}
		var argsByte = []byte(args_str)
		args = argsByte
	}
	JsonArgs := json.RawMessage(args)

	msg := types.NewMsgMigrateContract(hash, from, name, version, author, email, describe, contractAddress, JsonArgs)
	if !broadcast {
		rest.PostProcessResponseBare(w, cliCtx, hex.EncodeToString(msg.Bytes()))
		return
	}

	txByte, err := types2.SignCommonTx(from, nonce, gas, []sdk.Msg{msg}, privKey, cdc)
	if err != nil {
		rest.WriteErrorRes(w, types.ErrCheckParams(types.DefaultCodespace,err.Error()))
		return
	}
	res, err := cliCtx.BroadcastSignedTx(txByte)
	if err != nil {
		rest.WriteErrorRes(w, client.ErrBroadcast(types.DefaultCodespace, err))
		return
	}
	rest.PostProcessResponseBare(w, cliCtx, res)
}