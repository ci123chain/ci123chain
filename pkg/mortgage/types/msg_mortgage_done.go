package types

import (
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/transaction"
	"github.com/tanhuiya/ci123chain/pkg/util"
)

type MsgMortgageDone struct {
	transaction.CommonTx
	UniqueID  			[]byte			`json:"unique_id"`
}

func NewMsgMortgageDone(from sdk.AccAddress, gas, nonce uint64, uniqueID []byte) *MsgMortgageDone {
	msg := &MsgMortgageDone{
		CommonTx: transaction.CommonTx{
			From: from,
			Nonce: nonce,
			Gas:  gas,
		},
		UniqueID: 	uniqueID,
	}
	return msg
}

func (msg *MsgMortgageDone) Route() string {
	return RouterKey
}

func (msg *MsgMortgageDone) ValidateBasic() sdk.Error {
	if msg.CommonTx.From.Empty() {
		return sdk.ErrInvalidAddress("missing sender address")
	}
	if len(msg.UniqueID) < 1 {
		return sdk.ErrInternal("param mortgageRecord missing")
	}
	return msg.CommonTx.VerifySignature(msg.GetSignBytes(), true)
}


func (msg *MsgMortgageDone)GetSignBytes() []byte {
	ntx := *msg
	ntx.SetSignature(nil)
	return util.TxHash(ntx.Bytes())
}


func (msg *MsgMortgageDone)SetSignature(sig []byte) {
	msg.CommonTx.SetSignature(sig)
}

func (msg *MsgMortgageDone)Bytes() []byte {
	bytes, err := MortgageCdc.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}
	return bytes
}

func (msg *MsgMortgageDone)SetPubKey(pub []byte) {
	msg.CommonTx.PubKey = pub
}

func (msg *MsgMortgageDone) GetGas() uint64 {
	return msg.CommonTx.Gas
}

func (msg *MsgMortgageDone) GetNonce() uint64 {
	return msg.CommonTx.Nonce
}

func (msg *MsgMortgageDone) GetFromAddress() sdk.AccAddress {
	return msg.CommonTx.From
}