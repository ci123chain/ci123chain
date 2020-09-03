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
<<<<<<< HEAD
=======

	tx := wasm.NewStoreCodeTx(from, gas, nonce, sender, wasmcode)
	/*sid := cryptosuit.NewFabSignIdentity()
	pub, err  := sid.GetPubKey(privateKey)

	tx.SetPubKey(pub)
	signbyte := tx.GetSignBytes()
	signature, err := sid.Sign(signbyte, privateKey)*/
	eth := cryptosuit.NewETHSignIdentity()
	signature, err := eth.Sign(tx.GetSignBytes(), privateKey)
	if err != nil {
		return nil, err
	}
>>>>>>> mint
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
<<<<<<< HEAD
=======

	tx := wasm.NewInstantiateTx(from, gas, nonce, codeHash, sender, label, initMsg)
	/*sid := cryptosuit.NewFabSignIdentity()
	pub, err  := sid.GetPubKey(privateKey)

	tx.SetPubKey(pub)
	signbyte := tx.GetSignBytes()
	signature, err := sid.Sign(signbyte, privateKey)*/
	eth := cryptosuit.NewETHSignIdentity()
	signature, err := eth.Sign(tx.GetSignBytes(), privateKey)
	if err != nil {
		return nil, err
	}
>>>>>>> mint
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
<<<<<<< HEAD
=======

	tx := wasm.NewExecuteTx(from, gas, nonce, sender, contractAddress, msg)
	/*sid := cryptosuit.NewFabSignIdentity()
	pub, err  := sid.GetPubKey(privateKey)

	tx.SetPubKey(pub)
	signbyte := tx.GetSignBytes()
	signature, err := sid.Sign(signbyte, privateKey)*/
	eth := cryptosuit.NewETHSignIdentity()
	signature, err := eth.Sign(tx.GetSignBytes(), privateKey)
	if err != nil {
		return nil, err
	}
>>>>>>> mint
	tx.SetSignature(signature)
	return tx.Bytes(), nil
}