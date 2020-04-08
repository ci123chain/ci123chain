package rest

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/app"
	"github.com/tanhuiya/ci123chain/pkg/client"
	"github.com/tanhuiya/ci123chain/pkg/client/helper"
	wasm "github.com/tanhuiya/ci123chain/pkg/wasm/types"
	sdk "github.com/tanhuiya/ci123chain/sdk/wasm"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

var cdc = app.MakeCodec()

func buildStoreCodeMsg(r *http.Request) ([]byte, error) {

	var wasmcode []byte
	codeStr := r.FormValue("wasmCodeStr")
	if codeStr != "" {
		Byte, err := hex.DecodeString(codeStr)
		if err != nil {
			return nil, errors.New("invalid wasmcode")
		}
		wasmcode = Byte
	}else {
		file, _, err := r.FormFile("wasmCode")
		if err != nil {
			return nil, err
		}
		wasmcode, err = ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
	}
	if wasmcode == nil {
		return nil, errors.New("wasmcode can not be empty")
	}
	ok := wasm.IsValidaWasmFile(wasmcode)
	if ok != nil {
		return nil, ok
	}

	from, gas, nonce,  priv,_, err := getArgs(r)
	if err != nil {
		return nil, err
	}
	txByte, err := sdk.SignStoreCodeMsg(from, gas, nonce, priv, from, wasmcode)
	if err != nil {
		return nil, err
	}
	return txByte, nil
}

func buildInstantiateContractMsg(r *http.Request) ([]byte, error) {

	codeId := r.FormValue("codeID")
	if codeId == "" {
		return nil, errors.New("codeID can not be empty")
	}
	codeID, err := strconv.ParseUint(codeId, 10, 64)
	if err != nil {
		return nil, err
	}
	label := r.FormValue("label")
	if label == "" {
		label = "label"
	}
	from, gas, nonce, priv, Msg, err := getArgs(r)
	if err != nil {
		return nil, err
	}

	txByte, err := sdk.SignInstantiateContractMsg(from, gas, nonce, codeID, priv, from, label, Msg)
	if err != nil {
		return nil, err
	}
	return txByte, nil
}

func buildExecuteContractMsg(r *http.Request) ([]byte, error) {

	contractAddr := r.FormValue("contractAddress")
	if contractAddr == "" {
		return nil, errors.New("contractAddress can not be empty")
	}
	contractAddress := types.HexToAddress(contractAddr)

	from, gas, nonce, priv, Msg, err := getArgs(r)
	if err != nil {
		return nil, err
	}

	txByte, err := sdk.SignExecuteContractMsg(from, gas, nonce, priv, from, contractAddress, Msg)
	if err != nil {
		return nil, err
	}
	return txByte, nil
}

func getArgs(r *http.Request) (types.AccAddress, uint64, uint64, string, json.RawMessage, error) {
	var Msg []byte

	from := r.FormValue("from")
	inputGas := r.FormValue("gas")
	inputNonce := r.FormValue("nonce")
	froms, err := helper.ParseAddrs(from)
	if err != nil {
		return types.AccAddress{}, 0, 0,  "", nil, err
	}
	if len(froms) != 1 {
		return types.AccAddress{}, 0, 0, "", nil, err
	}
	gas, err := strconv.ParseUint(inputGas, 10, 64)
	if err != nil {
		return types.AccAddress{}, 0, 0,  "", nil, err
	}
	var nonce uint64
	if inputNonce != "" {
		UserNonce, err := strconv.ParseInt(inputNonce, 10, 64)
		if err != nil || UserNonce < 0 {
			return types.AccAddress{}, 0, 0, "", nil, err
		}
		nonce = uint64(UserNonce)
	}else {
		ctx, err := client.NewClientContextFromViper(cdc)
		if err != nil {
			return types.AccAddress{}, 0, 0, "", nil, err
		}
		nonce, err = ctx.GetNonceByAddress(froms[0])
		if err != nil {
			return types.AccAddress{}, 0, 0, "", nil, err
		}
	}
	priv := r.FormValue("privateKey")
	if priv == "" {
		return types.AccAddress{}, 0, 0, "", nil, errors.New("privateKey can not be empty")
	}

	msg := r.FormValue("msg")
	if msg == "" {
		Msg = nil
	}else {
		var Args []string
		//
		args := strings.Split(msg, ",")
		method := args[0]
		for i := 1; i < len(args); i++ {
			Args = append(Args, args[i])
		}
		param := wasm.NewCallContractParams(method, Args)
		callMsg, err := json.Marshal(param)
		if err != nil {
			return types.AccAddress{}, 0, 0, "", nil, errors.New("invalid json raw message")
		}
		Msg = callMsg
	}
	JsonMsg := json.RawMessage(Msg)

	return froms[0], gas, nonce, priv,JsonMsg, nil

}