package types

import (
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/transaction"
	"github.com/ci123chain/ci123chain/pkg/util"
	"github.com/tendermint/tendermint/crypto"
)


type CreateValidatorTx struct {
	//
	transaction.CommonTx
	PublicKey         crypto.PubKey      `json:"public_key"`
	Value             types.Coin         `json:"value"`
	ValidatorAddress  types.AccAddress   `json:"validator_address"`
	DelegatorAddress  types.AccAddress   `json:"delegator_address"`
	MinSelfDelegation  types.Int         `json:"min_self_delegation"`
	Commission        CommissionRates    `json:"commission"`
	Description       Description        `json:"description"`
}

func NewCreateValidatorTx(from types.AccAddress, gas ,nonce uint64, value types.Coin, minSelfDelegation types.Int, validatorAddr types.AccAddress, delegatorAddr types.AccAddress,
	rate, maxRate, maxChangeRate types.Dec, moniker, identity, website, securityContact, details string, publicKey crypto.PubKey ) CreateValidatorTx {
	return CreateValidatorTx{
		CommonTx: transaction.CommonTx{
			From: from,
			Gas: 	gas,
			Nonce: nonce,
		},
		PublicKey:publicKey,
		Value:value,
		ValidatorAddress:validatorAddr,
		DelegatorAddress:delegatorAddr,
		MinSelfDelegation:minSelfDelegation,
		Commission:NewCommissionRates(rate, maxRate, maxChangeRate),
		Description: NewDescription(moniker, identity, website, securityContact, details),
	}
}

func (msg *CreateValidatorTx) ValidateBasic() types.Error {

	err := msg.VerifySignature(msg.GetSignBytes(), false)
	if err != nil {
		return ErrCheckParams(DefaultCodespace, err.Error())
	}

	// note that unmarshaling from bech32 ensures either empty or valid
	if msg.DelegatorAddress.Empty() {
		return ErrInvalidAddress(DefaultCodespace, fmt.Sprintf("empty delegator address"))
	}
	if msg.ValidatorAddress.Empty() {
		return ErrInvalidAddress(DefaultCodespace, fmt.Sprintf("empty validator address"))
	}
	if !msg.ValidatorAddress.Equals(msg.DelegatorAddress) {
		return ErrInvalidAddress(DefaultCodespace, fmt.Sprintf("expected %s, got %s", msg.DelegatorAddress, msg.ValidatorAddress))
	}
	if msg.PublicKey == nil {
		return ErrEmptyPublicKey(DefaultCodespace, "empty publicKey")
	}
	if !msg.Value.Amount.IsPositive() {
		return ErrCheckParams(DefaultCodespace, "invalid amount")
	}
	if msg.Description == (Description{}) {
		return ErrCheckParams(DefaultCodespace, "description can not be empty")
	}
	if msg.Commission == (CommissionRates{}) {
		return ErrCheckParams(DefaultCodespace, "commission can not be empty")
	}
	if err := msg.Commission.Validate(); err != nil {
		return ErrCheckParams(DefaultCodespace, err.Error())
	}
	if !msg.MinSelfDelegation.IsPositive() {
		return ErrCheckParams(DefaultCodespace, "invalid minSelfDelegation")
	}
	if msg.Value.Amount.LT(msg.MinSelfDelegation) {
		return ErrCheckParams(DefaultCodespace, "self delegation below minnium")
	}

	return nil
}

func (msg *CreateValidatorTx) GetSignBytes() []byte {
	tmsg := *msg
	tmsg.Signature = nil
	signBytes := tmsg.Bytes()
	return util.TxHash(signBytes)
}
func (msg *CreateValidatorTx) SetSignature(sig []byte) {
	msg.Signature = sig
}
func (msg *CreateValidatorTx) Bytes() []byte {
	bytes, err := StakingCodec.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}

	return bytes
}
func (msg *CreateValidatorTx) SetPubKey(pub []byte) {
	msg.PubKey = pub
}
func (msg *CreateValidatorTx) Route() string {return RouteKey}
func (msg *CreateValidatorTx) GetGas() uint64 { return msg.Gas}

func (msg *CreateValidatorTx) GetNonce() uint64 { return msg.Nonce}
func (msg *CreateValidatorTx) GetFromAddress() types.AccAddress { return msg.From}