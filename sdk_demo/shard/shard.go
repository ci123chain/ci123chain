package main

import (
	"encoding/hex"
	sdk "github.com/tanhuiya/ci123chain/sdk/shard"
)

func signAddShardTxDemo() (string, error) {

	txByte, err := sdk.SignAddShardMsg(from, offlineGas, offlineNonce, t, name, offlineHeight, priv)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(txByte), nil
}

