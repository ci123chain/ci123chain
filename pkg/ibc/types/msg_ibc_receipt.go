package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/transaction"
	"github.com/ci123chain/ci123chain/pkg/util"
)

// bank to account
type IBCReceiveReceiptMsg struct {
	transaction.CommonTx
	RawMessage 	[]byte 	`json:"raw_message"`
}


func NewIBCReceiveReceiptMsg(from sdk.AccAddress, raw []byte, gas uint64, nonce uint64) *IBCReceiveReceiptMsg {
	return &IBCReceiveReceiptMsg{
		CommonTx: transaction.CommonTx{
			From:  from,
			Gas: 	gas,
			Nonce: nonce,
		},
		RawMessage: raw,
	}
}

func (msg *IBCReceiveReceiptMsg) ValidateBasic() sdk.Error {
	if err := msg.CommonTx.ValidateBasic(); err != nil {
		return err
	}
	// todo unmarshal to signedIBCMsg
	return msg.CommonTx.VerifySignature(msg.GetSignBytes(), true)
}


func (msg *IBCReceiveReceiptMsg)GetSignBytes() []byte {
	ntx := *msg
	ntx.SetSignature(nil)
	return util.TxHash(ntx.Bytes())
}


func (msg *IBCReceiveReceiptMsg)SetSignature(sig []byte) {
	msg.CommonTx.SetSignature(sig)
}

func (msg *IBCReceiveReceiptMsg)Bytes() []byte {
	bytes, err := IbcCdc.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}

	return bytes
}

func (msg *IBCReceiveReceiptMsg)SetPubKey(pub []byte) {
	msg.CommonTx.PubKey = pub
}

func (msg *IBCReceiveReceiptMsg) Route() string {
	return RouterKey
}

func (msg *IBCReceiveReceiptMsg) GetGas() uint64 {
	return msg.CommonTx.Gas
}

func (msg *IBCReceiveReceiptMsg) GetNonce() uint64 {
	return msg.CommonTx.Nonce
}

func (msg *IBCReceiveReceiptMsg) GetFromAddress() sdk.AccAddress {
	return msg.CommonTx.From
}