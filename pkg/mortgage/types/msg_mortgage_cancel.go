package types

import (
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/transaction"
	"github.com/tanhuiya/ci123chain/pkg/util"
)

type MsgMortgageCancel struct {
	transaction.CommonTx
	UniqueID  			[]byte			`json:"unique_id"`
}



func NewMsgMortgageCancel(from sdk.AccAddress, gas, nonce uint64, uniqueID []byte) *MsgMortgageCancel {
	msg := &MsgMortgageCancel{
		CommonTx: transaction.CommonTx{
			From: from,
			Nonce: nonce,
			Gas:  gas,
		},
		UniqueID: 	uniqueID,
	}
	return msg
}

func (msg *MsgMortgageCancel) Route() string {
	return RouterKey
}

func (msg *MsgMortgageCancel) ValidateBasic() sdk.Error {
	if msg.CommonTx.From.Empty() {
		return sdk.ErrInvalidAddress("missing sender address")
	}
	if len(msg.UniqueID) < 1 {
		return sdk.ErrInternal("param mortgageRecord missing")
	}
	return msg.CommonTx.VerifySignature(msg.GetSignBytes(), true)
}


func (msg *MsgMortgageCancel)GetSignBytes() []byte {
	ntx := *msg
	ntx.SetSignature(nil)
	return util.TxHash(ntx.Bytes())
}


func (msg *MsgMortgageCancel)SetSignature(sig []byte) {
	msg.CommonTx.SetSignature(sig)
}

func (msg *MsgMortgageCancel)Bytes() []byte {
	bytes, err := MortgageCdc.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}
	return bytes
}

func (msg *MsgMortgageCancel)SetPubKey(pub []byte) {
	msg.CommonTx.PubKey = pub
}
