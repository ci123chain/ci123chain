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
			return nil, errors.New("wasmCodeStr; cannot get wasm file: " + err.Error())
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

	from, gas, nonce,  priv, _, err := getArgs(r)
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
	from, gas, nonce, priv, args, err := getArgs(r)
	if err != nil {
		return nil, err
	}

	txByte, err := sdk.SignInstantiateContractMsg(from, gas, nonce, codeID, priv, from, label, args)
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

	from, gas, nonce, priv, args, err := getArgs(r)
	if err != nil {
		return nil, err
	}

	txByte, err := sdk.SignExecuteContractMsg(from, gas, nonce, priv, from, contractAddress, args)
	if err != nil {
		return nil, err
	}
	return txByte, nil
}

func getArgs(r *http.Request) (types.AccAddress, uint64, uint64, string, json.RawMessage, error) {
	var args []byte

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

	msg := r.FormValue("args")
	if msg == "" {
		args = nil
	}else {
		var Args []string
		var method string
		//
		str := strings.Split(msg, ",")
		if len(str) == 1 {
			method = str[0]
			Args = []string{}
		}else {
			method = str[0]
			for i := 1; i < len(str); i++ {
				Args = append(Args, str[i])
			}
		}
		param := wasm.NewCallContractParams(method, Args)
		callMsg, err := json.Marshal(param)
		if err != nil {
			return types.AccAddress{}, 0, 0, "", nil, errors.New("invalid json raw message")
		}
		args = callMsg
	}
	JsonArgs := json.RawMessage(args)

	return froms[0], gas, nonce, priv, JsonArgs, nil

}