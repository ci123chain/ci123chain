package types

import (
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	stakeingtypes "github.com/ci123chain/ci123chain/pkg/staking/types"
	"time"
)

type MsgPrestakingCreateValidatorDirect struct {
	FromAddress		  types.AccAddress	 `json:"from_address"`
	PublicKey         string      		 `json:"public_key"`
	ValidatorAddress  types.AccAddress   `json:"validator_address"`
	DelegatorAddress  types.AccAddress   `json:"delegator_address"`
	MinSelfDelegation types.Int          `json:"min_self_delegation"`
	Commission        stakeingtypes.CommissionRates    `json:"commission"`
	Description       stakeingtypes.Description        `json:"description"`
	Amount    		types.Coin		     `json:"amount"`
	DelegateTime  	time.Duration         `json:"delegate_time"`
}

func NewMsgCreateValidatorDirect(from types.AccAddress, minSelfDelegation types.Int, validatorAddr types.AccAddress, delegatorAddr types.AccAddress,
	rate, maxRate, maxChangeRate types.Dec, moniker, identity, website, securityContact, details string, publicKey string, amount types.Coin, dt time.Duration) *MsgPrestakingCreateValidatorDirect {
	return &MsgPrestakingCreateValidatorDirect{
		FromAddress: from,
		PublicKey:publicKey,
		ValidatorAddress:validatorAddr,
		DelegatorAddress:delegatorAddr,
		MinSelfDelegation:minSelfDelegation,
		Commission: stakeingtypes.NewCommissionRates(rate, maxRate, maxChangeRate),
		Description: stakeingtypes.NewDescription(moniker, identity, website, securityContact, details),
		Amount: amount,
		DelegateTime: dt,
	}
}

func (msg MsgPrestakingCreateValidatorDirect) Route() string {return RouteKey}

func (msg MsgPrestakingCreateValidatorDirect) MsgType() string {return "create-validator-direct"}

func (msg MsgPrestakingCreateValidatorDirect) ValidateBasic() error {
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
	if msg.Description == (stakeingtypes.Description{}) {
		return sdkerrors.Wrap(sdkerrors.ErrParams, "description can not be empty")
	}
	if msg.Commission == (stakeingtypes.CommissionRates{}) {
		return sdkerrors.Wrap(sdkerrors.ErrParams, "commission can not be empty")
	}
	if err := msg.Commission.Validate(); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrParams, err.Error())
	}
	if !msg.MinSelfDelegation.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrParams, "invalid minSelfDelegation")
	}
	if !msg.Amount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrParams, "amount can not be negative")
	}
	if msg.DelegateTime.Seconds() <= (time.Hour * 24 * 7).Seconds(){
		return sdkerrors.Wrap(sdkerrors.ErrParams, "the time should longer than 168h(1 week)")
	}
	return nil
}

func (msg MsgPrestakingCreateValidatorDirect) GetFromAddress() types.AccAddress { return msg.FromAddress}

func (msg MsgPrestakingCreateValidatorDirect) Bytes() []byte {
	bytes, err := PreStakingCodec.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}
	return bytes
}