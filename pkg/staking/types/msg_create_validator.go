package types

import (
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/util"
)

type MsgCreateValidator struct {
	FromAddress		  types.AccAddress	 `json:"from_address"`
	Signature 		  []byte   			 `json:"signature"`
	PubKey			  []byte			 `json:"pub_key"`
	PublicKey         string      		 `json:"public_key"`
	Value             types.Coin         `json:"value"`
	ValidatorAddress  types.AccAddress   `json:"validator_address"`
	DelegatorAddress  types.AccAddress   `json:"delegator_address"`
	MinSelfDelegation types.Int          `json:"min_self_delegation"`
	Commission        CommissionRates    `json:"commission"`
	Description       Description        `json:"description"`
}

func NewMsgCreateValidator(from types.AccAddress, value types.Coin, minSelfDelegation types.Int, validatorAddr types.AccAddress, delegatorAddr types.AccAddress,
	rate, maxRate, maxChangeRate types.Dec, moniker, identity, website, securityContact, details string, publicKey string ) *MsgCreateValidator {
	return &MsgCreateValidator{
		FromAddress: from,
		PublicKey:publicKey,
		Value:value,
		ValidatorAddress:validatorAddr,
		DelegatorAddress:delegatorAddr,
		MinSelfDelegation:minSelfDelegation,
		Commission:NewCommissionRates(rate, maxRate, maxChangeRate),
		Description: NewDescription(moniker, identity, website, securityContact, details),
	}
}

func (msg *MsgCreateValidator) Route() string {return RouteKey}

func (msg *MsgCreateValidator) MsgType() string {return "create-validator"}

func (msg *MsgCreateValidator) ValidateBasic() types.Error {
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
	if msg.PublicKey == "" {
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

func (msg *MsgCreateValidator) GetSignBytes() []byte {
	ntx := *msg
	ntx.SetSignature(nil)
	return util.TxHash(ntx.Bytes())
}
func (msg *MsgCreateValidator) SetSignature(sig []byte) {
	msg.Signature = sig
}

func (msg *MsgCreateValidator) SetPubKey(pub []byte) {
	msg.PubKey = pub
}

func (msg *MsgCreateValidator) GetFromAddress() types.AccAddress { return msg.FromAddress}

func (msg *MsgCreateValidator) Bytes() []byte {
	bytes, err := StakingCodec.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}
	return bytes
}

func (msg *MsgCreateValidator) GetSignature() []byte {
	return msg.Signature
}