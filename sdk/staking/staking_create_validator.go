package staking

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/app"
	"github.com/ci123chain/ci123chain/pkg/staking"
	"github.com/tendermint/go-amino"
	//"github.com/tendermint/tendermint/crypto"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
)

func SignCreateValidatorMSg(from sdk.AccAddress, gas, nonce, amount uint64, priv string, minSelfDelegation int64,
	validatorAddress, delegatorAddress sdk.AccAddress, rate, maxRate, maxChangeRate int64,
	moniker, identity, website, securityContact, details string, publicKey string) ([]byte, error) {

	/*by, err := hex.DecodeString(publicKey)
	if err != nil {
		return nil, err
	}
	var public crypto.PubKey
	err = cdc.UnmarshalJSON(by, &public)
	if err != nil {
		return nil, err
	}*/

	amt := sdk.NewUInt64Coin(amount)
	selfDelegation, r, mr, mxr := CreateParseArgs(minSelfDelegation, rate, maxRate, maxChangeRate)
	msg := staking.NewCreateValidatorMsg(from, amt, selfDelegation, validatorAddress, delegatorAddress, r, mr, mxr,
	moniker, identity, website, securityContact, details, publicKey)

	txByte, err := app.SignCommonTx(from, nonce, gas, []sdk.Msg{msg}, priv, cdc)
	if err != nil {
		return nil, err
	}
	return txByte, nil
}

/*func NewCreateValidatorMsg(from sdk.AccAddress, amt sdk.Coin, selfDelegation sdk.Int, validatorAddress, delegatorAddress sdk.AccAddress, r, mr, mxr sdk.Dec,
	moniker, identity, website, securityContact, details string, public crypto.PubKey) []byte {
	msg := staking.NewCreateValidatorMsg(from, amt, selfDelegation, validatorAddress, delegatorAddress, r, mr, mxr,
		moniker, identity, website, securityContact, details, public)
	return msg.Bytes()
}*/

var cdc = amino.NewCodec()

func init() {
	cryptoAmino.RegisterAmino(cdc)
}