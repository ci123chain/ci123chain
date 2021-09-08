package main

import (
	"encoding/hex"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
 	sdk "github.com/ci123chain/ci123chain/sdk/infrastructure"
)

func SignStoreContent(from types.AccAddress, gas, nonce uint64, priv string, key, content string) (string, error) {

	tx, err := sdk.SignInfrastructureStoreContent(from, gas, nonce, priv, key, content)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(tx), nil
}
