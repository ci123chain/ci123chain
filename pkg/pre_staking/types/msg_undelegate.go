package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
)


type MsgUndelegate struct {
	FromAddress sdk.AccAddress `json:"from_address"`
	Amount    sdk.Coin		  `json:"amount"`
}

func NewMsgUndelegate(from sdk.AccAddress, amount sdk.Coin) *MsgUndelegate {
	return &MsgUndelegate{
		FromAddress: from,
		Amount:      amount,
	}
}

func (msg *MsgUndelegate) Route() string { return ModuleName }

func (msg *MsgUndelegate) MsgType() string { return "pre-staking" }

func (msg *MsgUndelegate) ValidateBasic() error {
	if msg.FromAddress.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrParams, "from address can not empty")
	}
	if !msg.Amount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrParams, "amount can not be negative")
	}
	return nil
}


func (msg *MsgUndelegate) GetFromAddress() sdk.AccAddress { return msg.FromAddress}

func (msg *MsgUndelegate) Bytes() []byte{
	bytes, err := PreStakingCodec.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}
	return bytes
}