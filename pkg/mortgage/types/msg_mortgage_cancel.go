package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/ci123chain/ci123chain/pkg/util"
)

type MsgMortgageCancel struct {
	FromAddress sdk.AccAddress	`json:"from_address"`
	Signature 	[]byte   		`json:"signature"`
	PubKey		[]byte			`json:"pub_key"`
	UniqueID  	[]byte			`json:"unique_id"`
}

func NewMsgMortgageCancel(from sdk.AccAddress, uniqueID []byte) *MsgMortgageCancel {
	msg := &MsgMortgageCancel{
		FromAddress: from,
		UniqueID: 	uniqueID,
	}
	return msg
}

func (msg *MsgMortgageCancel) Route() string {
	return RouterKey
}

func (msg *MsgMortgageCancel) MsgType() string {
	return "mortgage_cancel"
}

func (msg *MsgMortgageCancel) ValidateBasic() error {
	if msg.FromAddress .Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "empty from address")
	}
	if len(msg.UniqueID) < 1 {
		return sdkerrors.Wrap(sdkerrors.ErrParams, "param mortgageRecord missing")
	}
	return nil
}

func (msg *MsgMortgageCancel) GetSignBytes() []byte {
	ntx := *msg
	ntx.SetSignature(nil)
	return util.TxHash(ntx.Bytes())
}

func (msg *MsgMortgageCancel) SetSignature(sig []byte) {
	msg.SetSignature(sig)
}

func (msg *MsgMortgageCancel) GetSignature() []byte {
	return msg.Signature
}

func (msg *MsgMortgageCancel) Bytes() []byte {
	bytes, err := MortgageCdc.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}
	return bytes
}

func (msg *MsgMortgageCancel) SetPubKey(pub []byte) {
	msg.PubKey = pub
}

func (msg *MsgMortgageCancel) GetFromAddress() sdk.AccAddress {
	return msg.FromAddress
}
