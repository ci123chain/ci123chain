package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"time"
)

type MsgPreStaking struct {
	FromAddress sdk.AccAddress   `json:"from_address"`
	Amount    sdk.Coin		     `json:"amount"`
	Contract  sdk.AccAddress     `json:"contract"`
	DelegateTime  time.Duration         `json:"delegate_time"`
}

func NewMsgPreStaking(from sdk.AccAddress, amount sdk.Coin, c sdk.AccAddress, dt time.Duration) *MsgPreStaking {
	return &MsgPreStaking{
		FromAddress: from,
		Amount:      amount,
		Contract:    c,
		DelegateTime: dt,
	}
}

func (msg *MsgPreStaking) Route() string { return ModuleName }

func (msg *MsgPreStaking) MsgType() string { return "pre-staking" }

func (msg *MsgPreStaking) ValidateBasic() error {
	if msg.FromAddress.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrParams, "from address can not empty")
	}
	if !msg.Amount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrParams, "amount can not be negative")
	}
	if msg.Contract.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrParams, "contract can not empty")
	}
	if msg.DelegateTime.Seconds() <= (time.Hour * 24 * 7).Seconds(){
		return sdkerrors.Wrap(sdkerrors.ErrParams, "the time should longer than 168h(1 week)")
	}
	return nil
}


func (msg *MsgPreStaking) GetFromAddress() sdk.AccAddress { return msg.FromAddress}

func (msg *MsgPreStaking) Bytes() []byte{
	bytes, err := PreStakingCodec.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}
	return bytes
}