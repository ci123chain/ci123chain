package staking

import (
	"encoding/hex"
	"github.com/ci123chain/ci123chain/pkg/cryptosuit"
	"github.com/ci123chain/ci123chain/pkg/staking"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
)

func SignCreateValidatorMSg(from string, amount int64, gas, nonce uint64, priv string, minSelfDelegation int64,
	validatorAddress, delegatorAddress string, rate, maxRate, maxChangeRate int64,
	moniker, identity, website, securityContact, details string, publicKey string) ([]byte, error) {

	by, err := hex.DecodeString(publicKey)
	if err != nil {
		return nil, err
	}
	var public crypto.PubKey
	err = cdc.UnmarshalJSON(by, &public)
	if err != nil {
		return nil, err
	}

	fromAddr, amt, validatorAddr, delegatorAddr, err := CommonParseArgs(from, amount, validatorAddress, delegatorAddress)
	if err != nil {
		return nil, err
	}

	selfDelegation, r, mr, mxr := CreateParseArgs(minSelfDelegation, rate, maxRate, maxChangeRate)
	tx := staking.NewCreateValidatorMsg(fromAddr, gas, nonce, amt,selfDelegation,validatorAddr, delegatorAddr,r,mr,mxr,
	moniker, identity, website, securityContact, details, public)
<<<<<<< HEAD
	var signature []byte
	privPub, err := hex.DecodeString(priv)
	eth := cryptosuit.NewETHSignIdentity()
	signature, err = eth.Sign(tx.GetSignBytes(), privPub)
=======

	/*sid := cryptosuit.NewFabSignIdentity()
	pub, err  := sid.GetPubKey(privateKey)

	tx.SetPubKey(pub)
	signbyte := tx.GetSignBytes()
	signature, err := sid.Sign(signbyte, privateKey)*/
	eth := cryptosuit.NewETHSignIdentity()
	signature, err := eth.Sign(tx.GetSignBytes(), privateKey)
>>>>>>> mint
	if err != nil {
		return nil, err
	}
	tx.SetSignature(signature)

	return tx.Bytes(), nil
}

var cdc = amino.NewCodec()

func init() {
	cryptoAmino.RegisterAmino(cdc)
}