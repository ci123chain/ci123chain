package types

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/ci123chain/ci123chain/pkg/util"
)

const (
	StateMortgaged = "StateMortgaged"
	StateSuccess = "StateSuccess"
	StateCancel = "StateCancel"
)

//var _ transaction.Transaction = (*MsgMortgage)(nil)

type MsgMortgage struct {
	FromAddress sdk.AccAddress	`json:"from_address"`
	Signature 	[]byte   		`json:"signature"`
	PubKey		[]byte			`json:"pub_key"`

	ToAddress 	 sdk.AccAddress `json:"to_address"`
	UniqueID 	 []byte 		`json:"unique_id"`
	Coin 	 	 sdk.Coin		`json:"coin"`
}

func NewMsgMortgage(from, to sdk.AccAddress, coin sdk.Coin, uniqueID []byte) *MsgMortgage {
	msg := &MsgMortgage{
		FromAddress: from,
		ToAddress: 	 to,
		UniqueID: 	 uniqueID,
		Coin: 		 coin,
	}
	return msg
}

func (msg *MsgMortgage) ValidateBasic() error {
	if msg.FromAddress.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "empty from address")
	}
	if msg.ToAddress.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "empty to address")
	}
	if len(msg.UniqueID) < 1 {
		return sdkerrors.Wrap(sdkerrors.ErrParams, "param mortgageRecord missing")
	}
	if !msg.Coin.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, fmt.Sprintf("coin is invalid" + msg.Coin.String()))
	}
	return nil
}

func (msg *MsgMortgage) GetSignBytes() []byte {
	ntx := *msg
	ntx.SetSignature(nil)
	return util.TxHash(ntx.Bytes())
}

func (msg *MsgMortgage) SetSignature(sig []byte) {
	msg.SetSignature(sig)
}

func (msg *MsgMortgage) Bytes() []byte {
	bytes, err := MortgageCdc.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}

	return bytes
}

func (msg *MsgMortgage) SetPubKey(pub []byte) {
	msg.PubKey = pub
}

func (msg *MsgMortgage) Route() string {
	return RouterKey
}

func (msg *MsgMortgage) MsgType() string {
	return "mortgage"
}

func (msg *MsgMortgage) GetFromAddress() sdk.AccAddress {
	return msg.FromAddress
}

func (msg *MsgMortgage) GetSignature() []byte {
	return msg.Signature
}

type Mortgage struct {
	MsgMortgage

	State  string `json:"state"`
}