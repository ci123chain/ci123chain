package staking

import (
	"encoding/hex"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/cryptosuit"
	"github.com/ci123chain/ci123chain/pkg/staking"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
)

func SignCreateValidatorMSg(from sdk.AccAddress, amount uint64, priv string, minSelfDelegation int64,
	validatorAddress, delegatorAddress sdk.AccAddress, rate, maxRate, maxChangeRate int64,
	moniker, identity, website, securityContact, details string, publicKey string) (sdk.Msg, error) {

	by, err := hex.DecodeString(publicKey)
	if err != nil {
		return nil, err
	}
	var public crypto.PubKey
	err = cdc.UnmarshalJSON(by, &public)
	if err != nil {
		return nil, err
	}

	amt := sdk.NewUInt64Coin(amount)
	selfDelegation, r, mr, mxr := CreateParseArgs(minSelfDelegation, rate, maxRate, maxChangeRate)
	msg := staking.NewCreateValidatorMsg(from, amt, selfDelegation, validatorAddress, delegatorAddress, r, mr, mxr,
	moniker, identity, website, securityContact, details, public)

	var signature []byte
	privPub, err := hex.DecodeString(priv)
	eth := cryptosuit.NewETHSignIdentity()
	signature, err = eth.Sign(msg.GetSignBytes(), privPub)

	if err != nil {
		return nil, err
	}
	msg.SetSignature(signature)

	return msg, nil
}

var cdc = amino.NewCodec()

func init() {
	cryptoAmino.RegisterAmino(cdc)
}