package main

import (
	"encoding/hex"
	sdk "github.com/ci123chain/ci123chain/sdk/distribution"
)


func SignCommunityPoolTx(from string, amount int64, gas, nonce uint64, priv string) (string, error) {

	txBytes, err := sdk.SignFundCommunityPoolTx(from, amount, gas, nonce, priv)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(txBytes), nil
}




func SignWithdrawCommissionTx(from, validatorAddress string, gas, nonce uint64, priv string) (string, error) {

	txBytes, err := sdk.SignWithdrawValidatorCommissionTx(from, validatorAddress, gas, nonce, priv)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(txBytes), nil
}



func SignWithdrawRewardsTx(from, validatorAddress, delegatorAddress string, gas, nonce uint64, priv string) (string, error) {

	txBytes, err := sdk.SignWithdrawDelegatorRewardTx(from, validatorAddress, delegatorAddress, gas, nonce, priv)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(txBytes), nil
}



func SignSetWithdrawAddressTx(from, withdrawAddress string, gas, nonce uint64, priv string) (string, error) {
	txBytes, err := sdk.SignSetWithdrawAddressTx(from, withdrawAddress, gas, nonce, priv)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(txBytes), nil
}