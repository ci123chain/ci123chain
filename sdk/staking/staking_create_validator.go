package staking

import (
	"github.com/tanhuiya/ci123chain/pkg/cryptosuit"
	"github.com/tanhuiya/ci123chain/pkg/staking"
)

func SignCreateValidatorMSg(from string, amount, gas, nonce uint64, priv []byte, minSelfDelegation int64,
	validatorAddress, delegatorAddress string, rate, maxRate, maxChangeRate int64,
	moniker, identity, website, securityContact, details string, pubKeyTp, pubKeyVal string) ([]byte, error) {

	fromAddr, amt, validatorAddr, delegatorAddr, err := CommonParseArgs(from, amount, validatorAddress, delegatorAddress)
	if err != nil {
		return nil, err
	}

	selfDelegation, r, mr, mxr := CreateParseArgs(minSelfDelegation, rate, maxRate, maxChangeRate)
	tx := staking.NewCreateValidatorMsg(fromAddr, gas, nonce, amt,selfDelegation,validatorAddr, delegatorAddr,r,mr,mxr,
	moniker, identity, website, securityContact, details, pubKeyTp, pubKeyVal)

	sid := cryptosuit.NewFabSignIdentity()
	pub, err  := sid.GetPubKey(priv)

	tx.SetPubKey(pub)
	signbyte := tx.GetSignBytes()
	signature, err := sid.Sign(signbyte, priv)
	tx.SetSignature(signature)

	return tx.Bytes(), nil
}