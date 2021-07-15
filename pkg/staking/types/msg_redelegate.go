package types

import (
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
)

type MsgRedelegate struct {
	FromAddress			 types.AccAddress	 `json:"from_address"`
	DelegatorAddress     types.AccAddress    `json:"delegator_address"`
	ValidatorSrcAddress  types.AccAddress	 `json:"validator_src_address"`
	ValidatorDstAddress  types.AccAddress	 `json:"validator_dst_address"`
	Amount               types.Coin	 		 `json:"amount"`
}

func NewMsgRedelegate(fromAddr types.AccAddress, delegateAddr types.AccAddress, validatorSrcAddr,
	validatorDstAddr types.AccAddress, amount types.Coin) *MsgRedelegate {
	//
	return &MsgRedelegate{
		FromAddress: 		 fromAddr,
		DelegatorAddress:    delegateAddr,
		ValidatorSrcAddress: validatorSrcAddr,
		ValidatorDstAddress: validatorDstAddr,
		Amount:              amount,
	}
}

func (msg *MsgRedelegate) ValidateBasic() error {
	if msg.DelegatorAddress.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "empty delegator address")
	}
	if msg.ValidatorSrcAddress.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "empty validatorSrc address")
	}
	if msg.ValidatorDstAddress.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "empty validatorDst address")
	}
	if !msg.Amount.Amount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrParams, "amount can not be negative")
	}
	if !msg.FromAddress.Equal(msg.DelegatorAddress) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("expected %s, got %s", msg.FromAddress.String(), msg.DelegatorAddress.String()))
	}
	return nil
}

func (msg *MsgRedelegate) Bytes() []byte {
	bytes, err := StakingCodec.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}

	return bytes
}
func (msg *MsgRedelegate) Route() string {return RouteKey}
func (msg *MsgRedelegate) MsgType() string {return "redelegate"}
func (msg *MsgRedelegate) GetFromAddress() types.AccAddress { return msg.FromAddress}