package main

import (
	"encoding/hex"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	sdk "github.com/ci123chain/ci123chain/sdk/transfer"
)


func SignTransferTxDemo() (string, error) {
	isFabric = false
	//isFabric = true
	msg, err := sdk.SignMsgTransfer(types.HexToAddress(from), types.HexToAddress(to), offlineAmount, priv, isFabric)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(msg.Bytes()), nil
}


