package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
)

type MsgDeploy struct {
	From   sdk.AccAddress   `json:"from"`
}


func NewMsgDeploy(from sdk.AccAddress) *MsgDeploy {
	return &MsgDeploy{
		From: from,
	}
}



func (msg *MsgDeploy) Route() string { return ModuleName }

func (msg *MsgDeploy) MsgType() string { return "pre-staking" }

func (msg *MsgDeploy) ValidateBasic() error {
	//if msg.FromAddress.Empty() {
	//	return sdkerrors.Wrap(sdkerrors.ErrParams, "from address can not empty")
	//}
	//if !msg.Amount.IsPositive() {
	//	return sdkerrors.Wrap(sdkerrors.ErrParams, "amount can not be negative")
	//}
	//if msg.Contract.Empty() {
	//	return sdkerrors.Wrap(sdkerrors.ErrParams, "contract can not empty")
	//}
	//if msg.DelegateTime.Seconds() <= (time.Second * 3600 * 24 * 3 ).Seconds(){
	//	return sdkerrors.Wrap(sdkerrors.ErrParams, "delegate_time can not be zero")
	//}
	return nil
}


func (msg *MsgDeploy) GetFromAddress() sdk.AccAddress { return msg.From}

func (msg *MsgDeploy) Bytes() []byte{
	bytes, err := PreStakingCodec.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}
	return bytes
}