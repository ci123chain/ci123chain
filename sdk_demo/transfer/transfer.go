package main

import (
	"encoding/hex"
	sdk "github.com/ci123chain/ci123chain/sdk/transfer"
)


func SignTransferTxDemo() (string, error) {
	isFabric = false
	//isFabric = true
	txByte, err := sdk.SignTransferMsg(from, to, offlineAmount, offlineGas, offlineNonce, priv, isFabric)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(txByte), nil
}


