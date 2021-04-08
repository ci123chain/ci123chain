package types

import (
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
)

type MsgCreateValidator struct {
	FromAddress		  types.AccAddress	 `json:"from_address"`
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

func (msg *MsgCreateValidator) ValidateBasic() error {
	// note that unmarshaling from bech32 ensures either empty or valid
	if msg.DelegatorAddress.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "empty delegator address")
	}
	if msg.ValidatorAddress.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "empty validator address")
	}
	if !msg.ValidatorAddress.Equals(msg.DelegatorAddress) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("expected %s, got %s", msg.DelegatorAddress, msg.ValidatorAddress))
	}
	if msg.PublicKey == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidPubKey, "empty publickey")
	}
	if !msg.Value.Amount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrParams, "amount can not be negative")
	}
	if msg.Description == (Description{}) {
		return sdkerrors.Wrap(sdkerrors.ErrParams, "description can not be empty")
	}
	if msg.Commission == (CommissionRates{}) {
		return sdkerrors.Wrap(sdkerrors.ErrParams, "commission can not be empty")
	}
	if err := msg.Commission.Validate(); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrParams, err.Error())
	}
	if !msg.MinSelfDelegation.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrParams, "invalid minSelfDelegation")
	}
	if msg.Value.Amount.LT(msg.MinSelfDelegation) {
		return sdkerrors.Wrap(sdkerrors.ErrParams, "self delegation must greater than minnium")
	}

	return nil
}

func (msg *MsgCreateValidator) GetFromAddress() types.AccAddress { return msg.FromAddress}

func (msg *MsgCreateValidator) Bytes() []byte {
	bytes, err := StakingCodec.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}
	return bytes
}