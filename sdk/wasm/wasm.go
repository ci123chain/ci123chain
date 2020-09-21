package wasm

import (
	"encoding/json"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/app"
	"github.com/ci123chain/ci123chain/pkg/wasm"
)

var cdc = app.MakeCodec()

func SignInstantiateContractMsg(code []byte,from sdk.AccAddress, gas, nonce uint64, priv string, name, version, author, email, describe string,
	initMsg json.RawMessage) ([]byte, error) {
	msg := wasm.NewInstantiateTx(code, from, name, version, author, email, describe, initMsg)
	txByte, err := app.SignCommonTx(from, nonce, gas, []sdk.Msg{msg}, priv, cdc)
	if err != nil {
		return nil, err
	}
	return txByte, nil
}

func NewInstantiateMsg(code []byte, from sdk.AccAddress, name, version, author, email, describe string, initMsg json.RawMessage) []byte {
	msg := wasm.NewInstantiateTx(code, from, name, version, author, email, describe, initMsg)
	return msg.Bytes()
}

func SignExecuteContractMsg(from sdk.AccAddress, gas, nonce uint64, priv string, contractAddress sdk.AccAddress, args json.RawMessage) ([]byte, error) {
	msg := wasm.NewExecuteTx(from, contractAddress, args)
	txByte, err := app.SignCommonTx(from, nonce, gas, []sdk.Msg{msg}, priv, cdc)
	if err != nil {
		return nil, err
	}
	return txByte, nil
}

func SignMigrateContractMsg(code []byte, from sdk.AccAddress, gas, nonce uint64, priv string, name, version, author, email, describe string,
	contractAddr sdk.AccAddress, initMsg json.RawMessage) ([]byte, error) {
	msg := wasm.NewMigrateTx(code, from, name, version, author, email, describe, contractAddr, initMsg)
	txByte, err := app.SignCommonTx(from, nonce, gas, []sdk.Msg{msg}, priv, cdc)
	if err != nil {
		return nil, err
	}
	return txByte, nil
}