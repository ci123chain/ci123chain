package main

import (
	"encoding/hex"
	sdk "github.com/tanhuiya/ci123chain/sdk/shard"
)

func signAddShardTxDemo() (string, error) {
	from := "0x3F43E75Aaba2c2fD6E227C10C6E7DC125A93DE3c"
	gas := uint64(20000)
	nonce := uint64(2)
	t := "ADD"
	name := "ciChain-2"
	height := int64(900)
	priv := "2b452434ac4f7cf9c5d61d62f23834f34e851fb6efdb8d4a8c6e214a8bc93d70"
	isFabric := false
	txByte, err := sdk.SignAddShardMsg(from, gas, nonce, t, name, height, priv, isFabric)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(txByte), nil
}

