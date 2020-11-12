package wasm

import (
	"encoding/json"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/vm"
)

var cdc = types.MakeCodec()

func SignUploadContractMsg(code []byte,from sdk.AccAddress, gas, nonce uint64, priv string) ([]byte, error) {
	msg := vm.NewUploadTx(code, from)
	txByte, err := types.SignCommonTx(from, nonce, gas, []sdk.Msg{msg}, priv, cdc)
	if err != nil {
		return nil, err
	}
	return txByte, nil
}

func SignInstantiateContractMsg(codeHash []byte, from sdk.AccAddress, gas, nonce uint64, priv string, name, version, author, email, describe string,
	initMsg json.RawMessage) ([]byte, error) {
	msg := vm.NewInstantiateTx(codeHash, from, name, version, author, email, describe, initMsg)
	txByte, err := types.SignCommonTx(from, nonce, gas, []sdk.Msg{msg}, priv, cdc)
	if err != nil {
		return nil, err
	}
	return txByte, nil
}

func NewInstantiateMsg(codeHash []byte, from sdk.AccAddress, name, version, author, email, describe string, initMsg json.RawMessage) []byte {
	msg := vm.NewInstantiateTx(codeHash, from, name, version, author, email, describe, initMsg)
	return msg.Bytes()
}

func SignExecuteContractMsg(from sdk.AccAddress, gas, nonce uint64, priv string, contractAddress sdk.AccAddress, args json.RawMessage) ([]byte, error) {
	msg := vm.NewExecuteTx(from, contractAddress, args)
	txByte, err := types.SignCommonTx(from, nonce, gas, []sdk.Msg{msg}, priv, cdc)
	if err != nil {
		return nil, err
	}
	return txByte, nil
}

func SignMigrateContractMsg(codeHash []byte, from sdk.AccAddress, gas, nonce uint64, priv string, name, version, author, email, describe string,
	contractAddr sdk.AccAddress, initMsg json.RawMessage) ([]byte, error) {
	msg := vm.NewMigrateTx(codeHash, from, name, version, author, email, describe, contractAddr, initMsg)
	txByte, err := types.SignCommonTx(from, nonce, gas, []sdk.Msg{msg}, priv, cdc)
	if err != nil {
		return nil, err
	}
	return txByte, nil
}