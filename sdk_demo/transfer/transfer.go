package main

import (
	"encoding/hex"
	sdk "github.com/tanhuiya/ci123chain/sdk/transfer"
)


func SignTransferTxDemo() (string, error) {
	var isFabric bool
	from := "0x3F43E75Aaba2c2fD6E227C10C6E7DC125A93DE3c"
	to := "0x505A74675dc9C71eF3CB5DF309256952917E801e"
	amount := uint64(2)
	gas := uint64(20000)
	nonce := uint64(2)
	priv := "2b452434ac4f7cf9c5d61d62f23834f34e851fb6efdb8d4a8c6e214a8bc93d70"
	isFabric = false
	//isFabric = true
	txByte, err := sdk.SignTransferMsg(from, to, amount, gas, nonce, priv, isFabric)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(txByte), nil
}


