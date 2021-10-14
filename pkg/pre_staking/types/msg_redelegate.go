package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
)

type MsgRedelegate struct {
	FromAddress   sdk.AccAddress   `json:"from_address"`
	SrcValidator  sdk.AccAddress   `json:"src_validator"`
	DstValidator  sdk.AccAddress   `json:"dst_validator"`
}


func NewMsgRedelegate(from, src, dst sdk.AccAddress) *MsgRedelegate {
	return &MsgRedelegate{
		FromAddress:  from,
		SrcValidator: src,
		DstValidator: dst,
	}
}

func (msg *MsgRedelegate) Route() string { return ModuleName }

func (msg *MsgRedelegate) MsgType() string { return "redelegate" }

func (msg *MsgRedelegate) ValidateBasic() error {
	if msg.FromAddress.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrParams, "from address can not empty")
	}
	if msg.SrcValidator.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrParams, "src_validator address can not empty")
	}
	if msg.DstValidator.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrParams, "dst_validator address can not empty")
	}
	//if !msg.Amount.IsPositive() {
	//	return sdkerrors.Wrap(sdkerrors.ErrParams, "amount can not be negative")
	//}
	return nil
}


func (msg *MsgRedelegate) GetFromAddress() sdk.AccAddress { return msg.FromAddress}

func (msg *MsgRedelegate) Bytes() []byte{
	bytes, err := PreStakingCodec.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}
	return bytes
}