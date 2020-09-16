package wasm

import (
	"encoding/hex"
	"encoding/json"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/cryptosuit"
	"github.com/ci123chain/ci123chain/pkg/wasm"
)

func SignInstantiateContractMsg(code []byte,from sdk.AccAddress, priv string, name, version, author, email, describe string,
	initMsg json.RawMessage) (sdk.Msg, error) {
	tx := wasm.NewInstantiateTx(code, from, name, version, author, email, describe, initMsg)
	var signature []byte
	privPub, err := hex.DecodeString(priv)
	eth := cryptosuit.NewETHSignIdentity()
	signature, err = eth.Sign(tx.GetSignBytes(), privPub)
	if err != nil {
		return nil, err
	}
	tx.SetSignature(signature)
	return tx, nil
}

func SignExecuteContractMsg(from sdk.AccAddress, priv string, contractAddress sdk.AccAddress, msg json.RawMessage) (sdk.Msg, error) {
	tx := wasm.NewExecuteTx(from, contractAddress, msg)
	var signature []byte
	privPub, err := hex.DecodeString(priv)
	eth := cryptosuit.NewETHSignIdentity()
	signature, err = eth.Sign(tx.GetSignBytes(), privPub)
	if err != nil {
		return nil, err
	}
	tx.SetSignature(signature)
	return tx, nil
}

func SignMigrateContractMsg(code []byte, from sdk.AccAddress, priv string, name, version, author, email, describe string,
	contractAddr sdk.AccAddress, initMsg json.RawMessage) (sdk.Msg, error) {
	tx := wasm.NewMigrateTx(code, from, name, version, author, email, describe, contractAddr, initMsg)
	var signature []byte
	privPub, err := hex.DecodeString(priv)
	eth := cryptosuit.NewETHSignIdentity()
	signature, err = eth.Sign(tx.GetSignBytes(), privPub)
	if err != nil {
		return nil, err
	}

	tx.SetSignature(signature)
	return tx, nil
}