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
	"github.com/ci123chain/ci123chain/pkg/vm"
	"github.com/ci123chain/ci123chain/pkg/vm/keeper"
	vmmodule "github.com/ci123chain/ci123chain/pkg/vm/moduletypes"
	wasmtypes "github.com/ci123chain/ci123chain/pkg/vm/wasmtypes"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"math/big"
	"net/http"
	"strconv"
	"strings"
)

const CAN_MIGRATE string = `{"method":"canMigrate","args": [""]}`
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

	code, err := getWasmCode(r)
	if err != nil || code == nil {
		rest.WriteErrorRes(w, wasmtypes.ErrCheckParams(vmmodule.DefaultCodespace, "get wasmcode failed"))
		return
	}
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
	res, _, _, Err := cliCtx.Query("/custom/" + vm.ModuleName + "/" + wasmtypes.QueryCodeInfo, bz, false)
	if Err != nil {
		rest.WriteErrorRes(w, Err)
		return
	}
	if len(res) > 0 { //already exists
		rest.PostProcessResponseBare(w, cliCtx, string(codeHash))
		return
	}
	msg := vm.NewUploadTx(code, from)
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
	var args []byte
	args_str := r.FormValue("args")
	if args_str == "" {
		args = nil
	}else {
		var argsStr wasmtypes.CallContractParam
		ok, err := util.CheckJsonArgs(args_str, argsStr)
		if err != nil || !ok {
			rest.WriteErrorRes(w, wasmtypes.ErrCheckParams(vmmodule.DefaultCodespace, "get initArgs failed"))
			return
		}
		var argsByte = []byte(args_str)
		args = argsByte
	}
	JsonArgs := json.RawMessage(args)
	msg := vm.NewInstantiateTx(hash, from, name, version, author, email, describe, JsonArgs)
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
	params := wasmtypes.NewQueryContractInfoParams(contractAddress)
	bz, Er := cliCtx.Cdc.MarshalJSON(params)
	if Er != nil {
		rest.WriteErrorRes(w, sdk.ErrInternal("marshal failed"))
		return
	}

	resQuery, _, _, _:= cliCtx.Query("/custom/" + vmmodule.ModuleName + "/" + wasmtypes.QueryContractInfo, bz, false)
	if resQuery == nil {
		rest.WriteErrorRes(w, wasmtypes.ErrCheckParams(vmmodule.DefaultCodespace, "contract does not exist or get contract error"))
		return
	}
	var args []byte
	args_str := r.FormValue("args")
	if args_str == "" {
		args = nil
	}else {
		var argsStr wasmtypes.CallContractParam
		ok, err := util.CheckJsonArgs(args_str, argsStr)
		if err != nil || !ok {
			rest.WriteErrorRes(w, wasmtypes.ErrCheckParams(vmmodule.DefaultCodespace, "get initArgs failed"))
			return
		}
		var argsByte = []byte(args_str)
		args = argsByte
	}
	JsonArgs := json.RawMessage(args)
	msg := wasmtypes.NewMsgExecuteContract(from, contractAddress, JsonArgs)
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
	queryParam := []byte(CAN_MIGRATE)
	params := wasmtypes.NewContractStateParam(contractAddress, sender, queryParam)
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

	var args []byte
	args_str := r.FormValue("args")
	if args_str == "" {
		args = nil
	}else {
		var argsStr wasmtypes.CallContractParam
		ok, err := util.CheckJsonArgs(args_str, argsStr)
		if err != nil || !ok {
			rest.WriteErrorRes(w, wasmtypes.ErrCheckParams(vmmodule.DefaultCodespace, "get initArgs failed"))
			return
		}
		var argsByte = []byte(args_str)
		args = argsByte
	}
	JsonArgs := json.RawMessage(args)

	msg := wasmtypes.NewMsgMigrateContract(hash, from, name, version, author, email, describe, contractAddress, JsonArgs)
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




//////////////////
func upuploadContractHandler(cliCtx context.Context, w http.ResponseWriter, r *http.Request) {
	broadcast, err := strconv.ParseBool(r.FormValue("broadcast"))
	if err != nil {
		broadcast = true
	}
	privKey, from, nonce, gas, err := rest.GetNecessaryParams(cliCtx, r, cdc, broadcast)
	if err != nil {
		rest.WriteErrorRes(w, wasmtypes.ErrCheckParams(vmmodule.DefaultCodespace, err.Error()))
		return
	}

	code, err := getEvmCode(r)
	if err != nil || code == nil {
		rest.WriteErrorRes(w, wasmtypes.ErrCheckParams(vmmodule.DefaultCodespace, "get wasmcode failed"))
		return
	}

	var to *ethcmn.Address
	to_str := r.FormValue("to")
	if len(to_str) == 0 {
		to = nil
	} else {
		to_addr := ethcmn.HexToAddress(to_str)
		to = &to_addr
	}

	amount_str := r.FormValue("amount")
	amount_int64, _ := strconv.ParseInt(amount_str, 10, 64)
	amount := big.NewInt(amount_int64)
	msg := vm.NewMsgEvmTx(from, nonce, to, amount, gas, big.NewInt(1), code)
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

func exexcuteContractHandler(cliCtx context.Context, w http.ResponseWriter, r *http.Request) {
	broadcast, err := strconv.ParseBool(r.FormValue("broadcast"))
	if err != nil {
		broadcast = true
	}
	privKey, from, nonce, gas, err := rest.GetNecessaryParams(cliCtx, r, cdc, broadcast)
	if err != nil {
		rest.WriteErrorRes(w, wasmtypes.ErrCheckParams(vmmodule.DefaultCodespace, err.Error()))
		return
	}

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

	data_str := r.FormValue("data")
	if len(to_str) == 0 {
		rest.WriteErrorRes(w, wasmtypes.ErrCheckParams(vmmodule.DefaultCodespace, "data cannot be empty"))
		return
	}
	data, _ := hex.DecodeString(data_str)
	msg := vm.NewMsgEvmTx(from, nonce, to, amount, gas, big.NewInt(1), data)
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