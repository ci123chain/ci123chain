package rest

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	"github.com/ci123chain/ci123chain/pkg/account/exported"
	keeper2 "github.com/ci123chain/ci123chain/pkg/account/keeper"
	"github.com/ci123chain/ci123chain/pkg/account/types"
	types2 "github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/client/context"
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
		rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrParams, "invalid new_rate").Error())
		return
	}

	code, err := getCode(r)
	if err != nil || code == nil {
		rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrParams, "get contract code failed").Error())
		return
	}

	var msg sdk.Msg
	if keeper.IsWasm(code) {
		wasmCode, err := keeper.UnCompress(code)
		if err != nil {
			rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrParams, "uncompress code failed").Error())
			return
		}
		codeHash := keeper.MakeCodeHash(wasmCode)
		params := wasmtypes.NewQueryCodeInfoParams(strings.ToUpper(hex.EncodeToString(codeHash)))
		bz, Er := cliCtx.Cdc.MarshalJSON(params)
		if Er != nil {
			rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, Er.Error()).Error())
			return
		}
		res, _, _, Err := cliCtx.Query("/custom/" + vmmodule.ModuleName + "/" + wasmtypes.QueryCodeInfo, bz, false)
		if Err != nil {
			rest.WriteErrorRes(w, Err.Error())
			return
		}
		if len(res) > 0 { //already exists
			rest.PostProcessResponseBare(w, cliCtx, sdk.TxResponse{Code: 0, Data: strings.ToUpper(hex.EncodeToString(codeHash))})
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
		rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrInternal, err.Error()).Error())
		return
	}
	resp, err := cliCtx.BroadcastSignedTx(txByte)
	if err != nil {
		rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrInternal, err.Error()).Error())
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
		rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrParams, err.Error()).Error())
		return
	}
	codeHash := r.FormValue("code_hash")
	hash, err := hex.DecodeString(strings.ToLower(codeHash))
	if err != nil {
		rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrParams, err.Error()).Error())
	}
	name, version, author, email, describe, err := adjustInstantiateParams(r)
	if err != nil {
		rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrParams, err.Error()).Error())
		return
	}
	var args utils.CallData
	argsStr := r.FormValue("calldata")
	if argsStr == "" {
		rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrParams, "get calldata failed").Error())
		return
	}else {
		err := json.Unmarshal([]byte(argsStr), &args)
		if err != nil  {
			rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error()).Error())
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
		rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrInternal, err.Error()).Error())
		return
	}
	res, err := cliCtx.BroadcastSignedTx(txByte)
	if err != nil {
		rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrInternal, err.Error()).Error())
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
		rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrParams, err.Error()).Error())
		return
	}

	contractAddr := r.FormValue("contract_address")
	contractAddress := sdk.HexToAddress(contractAddr)
	qparams := keeper2.NewQueryAccountParams(contractAddress, -1)
	bz, err := cliCtx.Cdc.MarshalJSON(qparams)
	if err != nil {
		rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrParams, err.Error()).Error())
		return
	}

	queryRes, _, _, err := cliCtx.Query("/custom/" + types.ModuleName + "/" + types.QueryAccount, bz, false)
	if err != nil {
		rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrInternal, err.Error()).Error())
		return
	}
	if queryRes == nil{
		rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrInternal, "contract account does not exist").Error())
		return
	}

	var acc exported.Account
	err2 := cliCtx.Cdc.UnmarshalBinaryLengthPrefixed(queryRes, &acc)
	if err2 != nil {
		rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("cdc unmarshal faield: %v", err2)).Error())
		return
	}
	var msg sdk.Msg
	var args utils.CallData
	args_str := r.FormValue("calldata")
	if args_str == "" {
		rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrParams, "invalid calldata").Error())
		return
	}else {
		err := json.Unmarshal([]byte(args_str), &args)
		if err != nil  {
			rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error()).Error())
			return
		}
	}
	if acc.GetContractType() == types.WasmContractType {
		msg = wasmtypes.NewMsgExecuteContract(from, contractAddress, args)
	} else if acc.GetContractType() == types.EvmContractType {
		var to *ethcmn.Address
		to_str := r.FormValue("contract_address")
		if len(to_str) == 0 {
			rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrParams, "contract_address cannot be empty").Error())
			return
		} else {
			to_addr := ethcmn.HexToAddress(to_str)
			to = &to_addr
		}
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
		if err != nil {
			rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrParams, "encode evm callData failed").Error())
			return
		}
		msg = evm.NewMsgEvmTx(from, nonce, to, amount, gas, big.NewInt(1), payload)
	} else {
		rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrParams, "encode evm callData failed").Error())
		return
	}

	if !broadcast {
		rest.PostProcessResponseBare(w, cliCtx, hex.EncodeToString(msg.Bytes()))
		return
	}
	txByte, err := types2.SignCommonTx(from, nonce, gas, []sdk.Msg{msg}, privKey, cdc)
	if err != nil {
		rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("sign tx failed: %v", err.Error())).Error())
		return
	}
	res, err := cliCtx.BroadcastSignedTx(txByte)
	if err != nil {
		rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("boradcast tx faiedl: %v", err.Error())).Error())
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
		rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrParams, err.Error()).Error())
		return
	}

	sender := cliCtx.FromAddr
	contractAddr := r.FormValue("contract_address")
	contractAddress := sdk.HexToAddress(contractAddr)
	var arg utils.CallData
	queryParam := []byte(CAN_MIGRATE)
	err = json.Unmarshal(queryParam, &arg)
	if err != nil {
		rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error()).Error())
	}
	params := wasmtypes.NewContractStateParam(contractAddress, sender, arg)
	bz, Er := cliCtx.Cdc.MarshalJSON(params)
	if Er != nil {
		rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("cdc marhsal faield:%v", Er.Error())).Error())
		return
	}

	resQuery, _, _, _ := cliCtx.Query("/custom/" + vmmodule.ModuleName + "/" + wasmtypes.QueryContractState, bz, false)
	var contractState wasmtypes.ContractState
	cliCtx.Cdc.MustUnmarshalJSON(resQuery, &contractState)
	if contractState.Result != "true" {
		rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrResponse, "No permissions to migrate contracts").Error())
		return
	}

	codeHash := r.FormValue("code_hash")
	hash, err := hex.DecodeString(strings.ToLower(codeHash))
	if err != nil {
		rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrParams, "invalid code_hash").Error())
		return
	}
	name, version, author, email, describe, err := adjustInstantiateParams(r)
	if err != nil {
		rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrParams, err.Error()).Error())
		return
	}

	var args utils.CallData
	args_str := r.FormValue("calldata")
	if args_str == "" {
		rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrParams, "invalid calldata").Error())
		return
	}else {
		err := json.Unmarshal([]byte(args_str), &args)
		if err != nil  {
			rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error()).Error())
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
		rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("sign tx faied: %v", err.Error())).Error())
		return
	}
	res, err := cliCtx.BroadcastSignedTx(txByte)
	if err != nil {
		rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("boradcast tx failed: %v", err.Error())).Error())
		return
	}
	rest.PostProcessResponseBare(w, cliCtx, res)
}