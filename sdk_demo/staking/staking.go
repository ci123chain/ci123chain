package main

import (
	"encoding/hex"
	sdk "github.com/ci123chain/ci123chain/sdk/staking"
)

func SignCreateValidatorTx(from string, amount int64, gas, nonce uint64, priv string, minSelfDelegation int64,
	validatorAddress, delegatorAddress string, rate, maxRate, maxChangeRate int64,
	moniker, identity, website, securityContact, details string, publicKey string) (string, error) {
	//

	txBytes, err := sdk.SignCreateValidatorMSg(from, amount, gas, nonce, priv, minSelfDelegation, validatorAddress,
		delegatorAddress, rate, maxRate, maxChangeRate, moniker, identity, website, securityContact, details, publicKey)

	if err != nil {
		return "", err
	}
	return hex.EncodeToString(txBytes), nil
}

func SignDelegateTx(from string, amount int64, gas, nonce uint64, priv string, validatorAddress, delegatorAddress string) (string, error) {

	txBytes, err := sdk.SignDelegateMsg(from, amount, gas, nonce, priv, validatorAddress, delegatorAddress)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(txBytes), nil
}

func SignRelegateTx(from string, amount int64, gas, nonce uint64, priv string, validatorSrcAddress, validatorDstAddress, delegatorAddress string) (string, error) {

	//
	txBytes, err := sdk.SignRedelegateMsg(from, amount, gas, nonce, priv, validatorSrcAddress, validatorDstAddress, delegatorAddress)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(txBytes), nil
}

func SignUndelegate(from string, amount int64, gas, nonce uint64, priv string, validatorAddress, delegatorAddress string) (string, error) {


	txBytes, err := sdk.SignUndelegateMsg(from, amount, gas, nonce, priv, validatorAddress, delegatorAddress)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(txBytes), nil
}