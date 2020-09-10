package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/transfer"
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

func (msg *MsgMortgage) ValidateBasic() sdk.Error {
	if msg.FromAddress.Empty() {
		return transfer.ErrCheckParams(DefaultCodespace, "missing from address")
	}
	if msg.ToAddress.Empty() {
		return transfer.ErrCheckParams(DefaultCodespace, "missing to address")
	}
	if len(msg.UniqueID) < 1 {
		return transfer.ErrCheckParams(DefaultCodespace, "param mortgageRecord missing")
	}
	if !msg.Coin.IsValid() {
		return transfer.ErrCheckParams(DefaultCodespace, "coin is invalid" + msg.Coin.String())
	}
	return nil
	//return msg.CommonTx.VerifySignature(msg.GetSignBytes(), true)
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