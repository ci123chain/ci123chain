package wasm

import (
	"encoding/hex"
	"encoding/json"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/cryptosuit"
	"github.com/ci123chain/ci123chain/pkg/wasm"
)

func SignInstantiateContractMsg(code []byte,from sdk.AccAddress, gas, nonce uint64, priv string, sender sdk.AccAddress, name, version, author, email, describe string,
	initMsg json.RawMessage) ([]byte, error) {
	privateKey, err := hex.DecodeString(priv)
	if err != nil {
		return nil, err
	}

	tx := wasm.NewInstantiateTx(code, from, gas, nonce, sender, name, version, author, email, describe, initMsg)
	sid := cryptosuit.NewFabSignIdentity()
	pub, err  := sid.GetPubKey(privateKey)

	tx.SetPubKey(pub)
	signbyte := tx.GetSignBytes()
	signature, err := sid.Sign(signbyte, privateKey)
	tx.SetSignature(signature)
	return tx.Bytes(), nil
}

func SignExecuteContractMsg(from sdk.AccAddress, gas, nonce uint64, priv string, sender, contractAddress sdk.AccAddress, msg json.RawMessage) ([]byte, error) {
	privateKey, err := hex.DecodeString(priv)
	if err != nil {
		return nil, err
	}

	tx := wasm.NewExecuteTx(from, gas, nonce, sender, contractAddress, msg)
	sid := cryptosuit.NewFabSignIdentity()
	pub, err  := sid.GetPubKey(privateKey)

	tx.SetPubKey(pub)
	signbyte := tx.GetSignBytes()
	signature, err := sid.Sign(signbyte, privateKey)
	tx.SetSignature(signature)
	return tx.Bytes(), nil
}

func SignMigrateContractMsg(code []byte, from sdk.AccAddress, gas, nonce uint64, priv string, sender sdk.AccAddress, name, version, author, email, describe string,
	contractAddr sdk.AccAddress, initMsg json.RawMessage) ([]byte, error) {
	privateKey, err := hex.DecodeString(priv)
	if err != nil {
		return nil, err
	}

	tx := wasm.NewMigrateTx(code, from, gas, nonce, sender, name, version, author, email, describe, contractAddr, initMsg)
	sid := cryptosuit.NewFabSignIdentity()
	pub, err  := sid.GetPubKey(privateKey)

	tx.SetPubKey(pub)
	signbyte := tx.GetSignBytes()
	signature, err := sid.Sign(signbyte, privateKey)
	tx.SetSignature(signature)
	return tx.Bytes(), nil
}