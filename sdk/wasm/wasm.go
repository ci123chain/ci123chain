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
	tx := wasm.NewInstantiateTx(code, from, gas, nonce, sender, name, version, author, email, describe, initMsg)
	var signature []byte
	privPub, err := hex.DecodeString(priv)
	eth := cryptosuit.NewETHSignIdentity()
	signature, err = eth.Sign(tx.GetSignBytes(), privPub)
	if err != nil {
		return nil, err
	}
	tx.SetSignature(signature)
	return tx.Bytes(), nil
}

func SignExecuteContractMsg(from sdk.AccAddress, gas, nonce uint64, priv string, sender, contractAddress sdk.AccAddress, msg json.RawMessage) ([]byte, error) {
	tx := wasm.NewExecuteTx(from, gas, nonce, sender, contractAddress, msg)
	var signature []byte
	privPub, err := hex.DecodeString(priv)
	eth := cryptosuit.NewETHSignIdentity()
	signature, err = eth.Sign(tx.GetSignBytes(), privPub)
	if err != nil {
		return nil, err
	}
	tx.SetSignature(signature)
	return tx.Bytes(), nil
}

func SignMigrateContractMsg(code []byte, from sdk.AccAddress, gas, nonce uint64, priv string, sender sdk.AccAddress, name, version, author, email, describe string,
	contractAddr sdk.AccAddress, initMsg json.RawMessage) ([]byte, error) {
	tx := wasm.NewMigrateTx(code, from, gas, nonce, sender, name, version, author, email, describe, contractAddr, initMsg)
	var signature []byte
	privPub, err := hex.DecodeString(priv)
	eth := cryptosuit.NewETHSignIdentity()
	signature, err = eth.Sign(tx.GetSignBytes(), privPub)
	if err != nil {
		return nil, err
	}

	tx.SetSignature(signature)
	return tx.Bytes(), nil
}