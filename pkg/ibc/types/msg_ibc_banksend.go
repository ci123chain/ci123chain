package types

import (
	"github.com/tanhuiya/ci123chain/pkg/transaction"
	"github.com/tanhuiya/ci123chain/pkg/util"
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
)

// bank to account
type IBCMsgBankSend struct {
	transaction.CommonTx
	RawMessage 	[]byte 	`json:"raw_message"`
}



func NewIBCMsgBankSendMsg(from sdk.AccAddress, raw []byte, gas uint64, nonce uint64) *IBCMsgBankSend {
	return &IBCMsgBankSend{
		CommonTx: transaction.CommonTx{
			From:  from,
			Gas: 	gas,
			Nonce: nonce,
		},
		RawMessage: raw,
	}
}

func (msg *IBCMsgBankSend) ValidateBasic() sdk.Error {
	if err := msg.CommonTx.ValidateBasic(); err != nil {
		return err
	}
	// todo unmarshal to signedIBCMsg
	return msg.CommonTx.VerifySignature(msg.GetSignBytes(), true)
}


func (msg *IBCMsgBankSend)GetSignBytes() []byte {
	ntx := *msg
	ntx.SetSignature(nil)
	return util.TxHash(ntx.Bytes())
}


func (msg *IBCMsgBankSend)SetSignature(sig []byte) {
	msg.CommonTx.SetSignature(sig)
}

func (msg *IBCMsgBankSend)Bytes() []byte {
	bytes, err := IbcCdc.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}

	return bytes
}

func (msg *IBCMsgBankSend)SetPubKey(pub []byte) {
	msg.CommonTx.PubKey = pub
}

func (msg *IBCMsgBankSend) Route() string {
	return RouterKey
}

func (msg *IBCMsgBankSend) GetGas() uint64 {
	return msg.CommonTx.Gas
}

func (msg *IBCMsgBankSend) GetNonce() uint64 {
	return msg.CommonTx.Nonce
}

func (msg *IBCMsgBankSend) GetFromAddress() sdk.AccAddress {
	return msg.CommonTx.From
}