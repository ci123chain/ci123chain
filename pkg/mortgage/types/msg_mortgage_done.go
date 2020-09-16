package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/transfer"
	"github.com/ci123chain/ci123chain/pkg/util"
)

type MsgMortgageDone struct {
	FromAddress sdk.AccAddress	`json:"from_address"`
	Signature 	[]byte   		`json:"signature"`
	PubKey		[]byte			`json:"pub_key"`
	UniqueID    []byte			`json:"unique_id"`
}

func NewMsgMortgageDone(from sdk.AccAddress, gas, nonce uint64, uniqueID []byte) *MsgMortgageDone {
	msg := &MsgMortgageDone{
		FromAddress: from,
		UniqueID: 	uniqueID,
	}
	return msg
}

func (msg *MsgMortgageDone) Route() string {
	return RouterKey
}

func (msg *MsgMortgageDone) MsgType() string {
	return "mortgage_done"
}

func (msg *MsgMortgageDone) ValidateBasic() sdk.Error {
	if msg.FromAddress.Empty() {
		return transfer.ErrCheckParams(DefaultCodespace, "missing sender address")
	}
	if len(msg.UniqueID) < 1 {
		return transfer.ErrCheckParams(DefaultCodespace, "param mortgageRecord missing")
	}
	return nil
	//return msg.CommonTx.VerifySignature(msg.GetSignBytes(), true)
}

func (msg *MsgMortgageDone) GetSignBytes() []byte {
	ntx := *msg
	ntx.SetSignature(nil)
	return util.TxHash(ntx.Bytes())
}

func (msg *MsgMortgageDone) GetSignature() []byte {
	return msg.Signature
}

func (msg *MsgMortgageDone) SetSignature(sig []byte) {
	msg.SetSignature(sig)
}

func (msg *MsgMortgageDone) Bytes() []byte {
	bytes, err := MortgageCdc.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}
	return bytes
}

func (msg *MsgMortgageDone)SetPubKey(pub []byte) {
	msg.PubKey = pub
}

func (msg *MsgMortgageDone) GetFromAddress() sdk.AccAddress {
	return msg.FromAddress
}