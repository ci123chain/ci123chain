package main

import (
	"encoding/hex"
	sdk "github.com/ci123chain/ci123chain/sdk/transfer"
)


func SignTransferTxDemo() (string, error) {
	isFabric = false
	//isFabric = true
	msg, err := sdk.SignMsgTransfer(from, to, offlineGas, offlineNonce, offlineAmount, priv, isFabric)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(msg), nil
}


