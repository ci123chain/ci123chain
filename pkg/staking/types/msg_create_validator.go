package types

import (
	"encoding/json"
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/transaction"
	"github.com/tendermint/tendermint/crypto"
)

type CreateValidatorTx struct {
	//
	transaction.CommonTx
	PublicKey crypto.PubKey
	Value    types.Coin
	ValidatorAddress  types.AccAddress
	DelegatorAddress  types.AccAddress
	MinSelfDelegation  types.Int
	Commission        CommissionRates
	Description       Description
}

func NewCreateValidatorTx(from types.AccAddress, gas ,nonce uint64, value types.Coin, minSelfDelegation types.Int, validatorAddr types.AccAddress, delegatorAddr types.AccAddress,
	rate, maxRate, maxChangeRate types.Dec, moniker, identity, website, securityContact, details string, pubKey crypto.PubKey) *CreateValidatorTx {
	return &CreateValidatorTx{
		CommonTx: transaction.CommonTx{
			From: from,
			Gas: 	gas,
			Nonce: nonce,
		},
		PublicKey:pubKey,
		Value:value,
		ValidatorAddress:validatorAddr,
		DelegatorAddress:delegatorAddr,
		MinSelfDelegation:minSelfDelegation,
		Commission:NewCommissionRates(rate, maxRate, maxChangeRate),
		Description: NewDescription(moniker, identity, website, securityContact, details),
	}
}

func (msg *CreateValidatorTx) ValidateBasic() types.Error {

	// note that unmarshaling from bech32 ensures either empty or valid
	if msg.DelegatorAddress.Empty() {
		return types.ErrEmptyDelegatorAddr("empty delegator address")
	}
	if msg.ValidatorAddress.Empty() {
		return types.ErrEmptyValidatorAddr("empty validator address")
	}
	if !types.AccAddress(msg.ValidatorAddress).Equals(msg.DelegatorAddress) {
		return types.ErrBadValidatorAddr("bad validator address")
	}
	if msg.PublicKey == nil {
		return types.ErrEmptyValidatorPubKey("empty validator pubkey")
	}
	if !msg.Value.Amount.IsPositive() {
		return types.ErrBadDelegationAmount("bad delegation amount")
	}
	if msg.Description == (Description{}) {
		return types.ErrEmptyDescription("empty description")
	}
	if msg.Commission == (CommissionRates{}) {
		return types.ErremptyCommission("empty commission")
	}
	if err := msg.Commission.Validate(); err != nil {
		return types.ErrInvalidCommission("invalid commission")
	}
	if !msg.MinSelfDelegation.IsPositive() {
		return types.ErrMinSelfDelegationInvalid("invalid minself delegation")
	}
	if msg.Value.Amount.LT(msg.MinSelfDelegation) {
		return types.ErrSelfDelegationBelowMinimum("self delegation below minnium")
	}

	return nil
}

func (msg *CreateValidatorTx) GetSignBytes() []byte {
	tmsg := *msg
	tmsg.Signature = nil
	signBytes := tmsg.Bytes()
	return signBytes
}
func (msg *CreateValidatorTx) SetSignature(sig []byte) {
	msg.Signature = sig
}
func (msg *CreateValidatorTx) Bytes() []byte {
	bytes, err := json.Marshal(msg)
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