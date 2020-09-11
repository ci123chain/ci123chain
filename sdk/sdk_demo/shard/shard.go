package main

import (
	"encoding/hex"
	sdk "github.com/ci123chain/ci123chain/sdk/shard"
)

func signUpgradeMsgDemo() (string, error) {
	msg, err := sdk.SignUpgradeMsg(t, name, offlineHeight, priv)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(msg.Bytes()), nil
}

