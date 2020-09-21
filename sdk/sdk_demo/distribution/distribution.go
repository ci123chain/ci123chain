package main

import (
	"encoding/hex"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	sdk "github.com/ci123chain/ci123chain/sdk/distribution"
)


func SignCommunityPoolTx(from string, amount int64, gas, nonce uint64, priv string) (string, error) {

	txBytes, err := sdk.SignFundCommunityPoolTx(from, amount, gas, nonce, priv)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(txBytes), nil
}

func SignWithdrawCommissionTx(from, validator string, gas, nonce uint64, priv string) (string, error) {
	fromAddr := types.HexToAddress(from)
	validatorAddress := types.HexToAddress(validator)
	txBytes, err := sdk.SignWithdrawValidatorCommissionTx(fromAddr, validatorAddress, gas, nonce, priv)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(txBytes), nil
}

func SignWithdrawRewardsTx(from, validatorAddress, delegatorAddress string, gas, nonce uint64, priv string) (string, error) {
	fromAddr := types.HexToAddress(from)
	validatorAddr := types.HexToAddress(validatorAddress)
	delegatorAddr := types.HexToAddress(delegatorAddress)
	txBytes, err := sdk.SignWithdrawDelegatorRewardTx(fromAddr, validatorAddr, delegatorAddr, gas, nonce, priv)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(txBytes), nil
}



func SignSetWithdrawAddressTx(from, withdrawAddress string, gas, nonce uint64, priv string) (string, error) {
	fromAddr := types.HexToAddress(from)
	withdrawAddr := types.HexToAddress(withdrawAddress)
	txBytes, err := sdk.SignMsgSetWithdrawAddress(fromAddr, withdrawAddr, gas, nonce, priv)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(txBytes), nil
}