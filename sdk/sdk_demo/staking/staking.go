package main

import (
	"encoding/hex"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	sdk "github.com/ci123chain/ci123chain/sdk/staking"
)

func SignCreateValidatorTx(from types.AccAddress, amount uint64, priv string, minSelfDelegation int64,
	validatorAddress, delegatorAddress types.AccAddress, rate, maxRate, maxChangeRate int64,
	moniker, identity, website, securityContact, details string, publicKey string) (string, error) {
	//

	tx, err := sdk.SignCreateValidatorMSg(from, gas, nonce, amount, priv, minSelfDelegation, validatorAddress,
		delegatorAddress, rate, maxRate, maxChangeRate, moniker, identity, website, securityContact, details, publicKey)

	if err != nil {
		return "", err
	}
	return hex.EncodeToString(tx), nil
}

func SignDelegateTx(from types.AccAddress, gas, nonce, amount uint64, priv string, validatorAddress, delegatorAddress types.AccAddress) (string, error) {

	tx, err := sdk.SignDelegateMsg(from, gas, nonce, amount, priv, validatorAddress, delegatorAddress)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(tx), nil
}

func SignRelegateTx(from types.AccAddress, gas, nonce, amount uint64, priv string, validatorSrcAddress, validatorDstAddress, delegatorAddress types.AccAddress) (string, error) {

	//
	tx, err := sdk.SignRedelegateMsg(from, gas, nonce ,amount, priv, validatorSrcAddress, validatorDstAddress, delegatorAddress)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(tx), nil
}

func SignUndelegate(from types.AccAddress, gas, nonce, amount uint64, priv string, validatorAddress, delegatorAddress types.AccAddress) (string, error) {
	tx, err := sdk.SignUndelegateMsg(from, gas, nonce, amount, priv, validatorAddress, delegatorAddress)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(tx), nil
}