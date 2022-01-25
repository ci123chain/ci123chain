package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
)

type MsgSetStakingToken struct {
	FromAddress   sdk.AccAddress   `json:"from_address"`
	TokenAddress  sdk.AccAddress   `json:"token_address"`
}


func NewMsgSetStakingToken(from, token sdk.AccAddress) *MsgSetStakingToken {
	return &MsgSetStakingToken{
		FromAddress:  from,
		TokenAddress: token,
	}
}

func (msg *MsgSetStakingToken) Route() string { return ModuleName }

func (msg *MsgSetStakingToken) MsgType() string { return "set_staking_token" }

func (msg *MsgSetStakingToken) ValidateBasic() error {
	if msg.FromAddress.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrParams, "from address can not empty")
	}
	if msg.TokenAddress.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrParams, "token_address address can not empty")
	}
	//if !msg.Amount.IsPositive() {
	//	return sdkerrors.Wrap(sdkerrors.ErrParams, "amount can not be negative")
	//}
	return nil
}


func (msg *MsgSetStakingToken) GetFromAddress() sdk.AccAddress { return msg.FromAddress}

func (msg *MsgSetStakingToken) Bytes() []byte{
	bytes, err := PreStakingCodec.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}
	return bytes
}