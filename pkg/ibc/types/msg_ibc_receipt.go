package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/util"
)

// bank to account
type IBCReceiveReceiptMsg struct {
	FromAddress sdk.AccAddress	`json:"from_address"`
	Signature 	[]byte   		`json:"signature"`
	PubKey		[]byte			`json:"pub_key"`

	RawMessage 	[]byte 			`json:"raw_message"`
}

func NewIBCReceiveReceiptMsg(from sdk.AccAddress, raw []byte) *IBCReceiveReceiptMsg {
	return &IBCReceiveReceiptMsg{
		FromAddress: from,
		RawMessage:  raw,
	}
}

func (msg *IBCReceiveReceiptMsg) ValidateBasic() sdk.Error {
	if msg.FromAddress.Empty(){
		return sdk.ErrInvalidAddress("from address is empty")
	}
	return nil
	//return msg.CommonTx.VerifySignature(msg.GetSignBytes(), true)
}

func (msg *IBCReceiveReceiptMsg) GetSignBytes() []byte {
	ntx := *msg
	ntx.SetSignature(nil)
	return util.TxHash(ntx.Bytes())
}

func (msg *IBCReceiveReceiptMsg) SetSignature(sig []byte) {
	msg.SetSignature(sig)
}

func (msg *IBCReceiveReceiptMsg) GetSignature() []byte {
	return msg.Signature
}

func (msg *IBCReceiveReceiptMsg) Bytes() []byte {
	bytes, err := IbcCdc.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}
	return bytes
}

func (msg *IBCReceiveReceiptMsg) SetPubKey(pub []byte) {
	msg.PubKey = pub
}

func (msg *IBCReceiveReceiptMsg) Route() string {
	return RouterKey
}

func (msg *IBCReceiveReceiptMsg) MsgType() string {
	return "ibc_receive_receipt"
}

func (msg *IBCReceiveReceiptMsg) GetFromAddress() sdk.AccAddress {
	return msg.FromAddress
}