package wasm

import (
	"encoding/hex"
	"encoding/json"
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/cryptosuit"
	"github.com/tanhuiya/ci123chain/pkg/wasm"
)


func SignStoreCodeMsg(from sdk.AccAddress, gas, nonce uint64, priv string, sender sdk.AccAddress, wasmcode []byte, source, builder string) ([]byte, error) {

	privateKey, err := hex.DecodeString(priv)
	if err != nil {
		return nil, err
	}

	tx := wasm.NewStoreCodeTx(from, gas, nonce, sender, wasmcode, source, builder)
	sid := cryptosuit.NewFabSignIdentity()
	pub, err  := sid.GetPubKey(privateKey)

	tx.SetPubKey(pub)
	signbyte := tx.GetSignBytes()
	signature, err := sid.Sign(signbyte, privateKey)
	tx.SetSignature(signature)
	return tx.Bytes(), nil
}

func SignInstantiateContractMsg(from sdk.AccAddress, gas, nonce, codeID uint64, priv string, sender sdk.AccAddress, label string,
	initMsg json.RawMessage, initFunds sdk.Coin) ([]byte, error) {

	privateKey, err := hex.DecodeString(priv)
	if err != nil {
		return nil, err
	}

	tx := wasm.NewInstantiateTx(from, gas, nonce, codeID, sender, label, initMsg, initFunds)
	sid := cryptosuit.NewFabSignIdentity()
	pub, err  := sid.GetPubKey(privateKey)

	tx.SetPubKey(pub)
	signbyte := tx.GetSignBytes()
	signature, err := sid.Sign(signbyte, privateKey)
	tx.SetSignature(signature)
	return tx.Bytes(), nil
}


func SignExecuteContractMsg(from sdk.AccAddress, gas, nonce uint64, priv string, sender, contractAddress sdk.AccAddress, msg json.RawMessage, sendFunds sdk.Coin) ([]byte, error) {

	privateKey, err := hex.DecodeString(priv)
	if err != nil {
		return nil, err
	}

	tx := wasm.NewExecuteTx(from, gas, nonce, sender, contractAddress, msg, sendFunds)
	sid := cryptosuit.NewFabSignIdentity()
	pub, err  := sid.GetPubKey(privateKey)

	tx.SetPubKey(pub)
	signbyte := tx.GetSignBytes()
	signature, err := sid.Sign(signbyte, privateKey)
	tx.SetSignature(signature)
	return tx.Bytes(), nil
}