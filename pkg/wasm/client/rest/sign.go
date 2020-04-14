package rest

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/app"
	"github.com/tanhuiya/ci123chain/pkg/client"
	"github.com/tanhuiya/ci123chain/pkg/client/helper"
	"github.com/tanhuiya/ci123chain/pkg/util"
	wasm "github.com/tanhuiya/ci123chain/pkg/wasm/types"
	sdk "github.com/tanhuiya/ci123chain/sdk/wasm"
	"io/ioutil"
	"net/http"
	"strings"
)

var cdc = app.MakeCodec()

func buildStoreCodeMsg(r *http.Request) ([]byte, error) {
	var wasmcode []byte
	//check params
	from, gas, nonce,  priv, _, err := getArgs(r)
	if err != nil {
		return nil, err
	}
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
	txByte, err := sdk.SignStoreCodeMsg(from, gas, nonce, priv, from, wasmcode)
	if err != nil {
		return nil, err
	}
	return txByte, nil
}

func buildInstantiateContractMsg(r *http.Request) ([]byte, error) {

	codeHash := r.FormValue("codeHash")
	hash, err := hex.DecodeString(strings.ToLower(codeHash))
	if err != nil {
		return nil, errors.New("decode codeHash fail")
	}
	label := r.FormValue("label")
	if label == "" {
		label = "label"
	}else {
		err := util.CheckStringLength(1, 100, label)
		if err != nil {
			return nil, errors.New("error label")
		}
	}
	from, gas, nonce, priv, args, err := getArgs(r)
	if err != nil {
		return nil, err
	}

	txByte, err := sdk.SignInstantiateContractMsg(from, gas, nonce, hash, priv, from, label, args)
	if err != nil {
		return nil, err
	}
	return txByte, nil
}

func buildExecuteContractMsg(r *http.Request) ([]byte, error) {

	contractAddr := r.FormValue("contractAddress")
	err := util.CheckStringLength(42, 100, contractAddr)
	if err != nil {
		return nil, errors.New("error contractAddress")
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
	err := util.CheckStringLength(42, 100, from)
	if err != nil {
		return types.AccAddress{}, 0, 0,  "", nil, errors.New("error from")
	}
	froms, err := helper.ParseAddrs(from)
	if err != nil {
		return types.AccAddress{}, 0, 0,  "", nil, err
	}
	gas, err := util.CheckUint64(inputGas)
	if err != nil {
		return types.AccAddress{}, 0, 0,  "", nil, err
	}
	priv := r.FormValue("privateKey")
	err = util.CheckStringLength(1, 100, priv)
	if err != nil {
		return types.AccAddress{}, 0, 0, "", nil, errors.New("error privateKey")
	}
	msg := r.FormValue("args")
	if msg == "" {
		args = nil
	}else {
		var argsStr wasm.CallContractParam
		ok, err := util.CheckJsonArgs(msg, argsStr)
		if err != nil || !ok {
			return types.AccAddress{}, 0, 0, "", nil, errors.New("unexpected args")
		}
		var argsByte = []byte(msg)
		args = argsByte
	}
	JsonArgs := json.RawMessage(args)
	var nonce uint64
	if inputNonce != "" {
		UserNonce, err := util.CheckUint64(inputNonce)
		if err != nil || UserNonce < 0 {
			return types.AccAddress{}, 0, 0, "", nil, err
		}
		nonce = UserNonce
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

	return froms[0], gas, nonce, priv, JsonArgs, nil

}