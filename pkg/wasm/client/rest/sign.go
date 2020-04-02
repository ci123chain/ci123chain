package rest

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/app"
	"github.com/tanhuiya/ci123chain/pkg/client"
	"github.com/tanhuiya/ci123chain/pkg/client/helper"
	sdk "github.com/tanhuiya/ci123chain/sdk/wasm"
	"io/ioutil"
	"net/http"
	"strconv"
)

var cdc = app.MakeCodec()

func buildStoreCodeMsg(r *http.Request) ([]byte, error) {

	file, _, err := r.FormFile("wasmCode")
	if err != nil {
		return nil, err
	}
	wasmcode, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	if wasmcode == nil {
		return nil, errors.New("wasmcode can not be empty")
	}

	from, gas, nonce,  priv,_, _, err := getArgs(r)
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
	from, gas, nonce, priv, Msg, Funds, err := getArgs(r)
	if err != nil {
		return nil, err
	}

	txByte, err := sdk.SignInstantiateContractMsg(from, gas, nonce, codeID, priv, from, label, Msg, Funds)
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

	from, gas, nonce, priv, Msg, Funds, err := getArgs(r)
	if err != nil {
		return nil, err
	}

	txByte, err := sdk.SignExecuteContractMsg(from, gas, nonce, priv, from, contractAddress, Msg, Funds)
	if err != nil {
		return nil, err
	}
	return txByte, nil
}

func getArgs(r *http.Request) (types.AccAddress, uint64, uint64, string, json.RawMessage, types.Coin, error) {
	var Msg []byte
	var Str interface{}
	var Funds types.Coin

	from := r.FormValue("from")
	inputGas := r.FormValue("gas")
	inputNonce := r.FormValue("nonce")
	froms, err := helper.ParseAddrs(from)
	if err != nil {
		return types.AccAddress{}, 0, 0,  "", nil, types.Coin{}, err
	}
	if len(froms) != 1 {
		return types.AccAddress{}, 0, 0, "", nil, types.Coin{}, err
	}
	gas, err := strconv.ParseUint(inputGas, 10, 64)
	if err != nil {
		return types.AccAddress{}, 0, 0,  "", nil, types.Coin{}, err
	}
	var nonce uint64
	if inputNonce != "" {
		UserNonce, err := strconv.ParseInt(inputNonce, 10, 64)
		if err != nil || UserNonce < 0 {
			return types.AccAddress{}, 0, 0, "", nil, types.Coin{}, err
		}
		nonce = uint64(UserNonce)
	}else {
		ctx, err := client.NewClientContextFromViper(cdc)
		if err != nil {
			return types.AccAddress{}, 0, 0, "", nil, types.Coin{}, err
		}
		nonce, err = ctx.GetNonceByAddress(froms[0])
		if err != nil {
			return types.AccAddress{}, 0, 0, "", nil, types.Coin{}, err
		}
	}
	priv := r.FormValue("privateKey")
	if priv == "" {
		return types.AccAddress{}, 0, 0, "", nil, types.Coin{},errors.New("privateKey can not be empty")
	}

	msg := r.FormValue("msg")
	if msg == "" {
		Msg = nil
	}else {
		Msg, err = hex.DecodeString(msg)
		if err != nil {
			return types.AccAddress{}, 0, 0, "", nil, types.Coin{}, errors.New("error initMsg")
		}
		err = json.Unmarshal(Msg, &Str)
		if err != nil {
			return types.AccAddress{}, 0, 0, "", nil, types.Coin{}, errors.New("invalid json raw message")
		}
	}
	JsonMsg := json.RawMessage(Msg)
	funds := r.FormValue("funds")
	if funds == "" {
		Funds = types.NewCoin(types.NewInt(0))
	}else {
		fs, err := strconv.ParseInt(funds, 10, 64)
		if err != nil {
			return types.AccAddress{}, 0, 0, "", nil, types.Coin{}, err
		}
		Funds = types.NewCoin(types.NewInt(fs))
	}

	return froms[0], gas, nonce, priv,JsonMsg, Funds, nil

}