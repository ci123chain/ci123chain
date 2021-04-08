package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/util"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
)

// bank to account
type IBCMsgBankSend struct {
	FromAddress sdk.AccAddress	`json:"from_address"`
	Signature 	[]byte   		`json:"signature"`
	PubKey		[]byte			`json:"pub_key"`

	RawMessage 	[]byte 	        `json:"raw_message"`
}

func NewIBCMsgBankSendMsg(from sdk.AccAddress, raw []byte) *IBCMsgBankSend {
	return &IBCMsgBankSend{
		FromAddress: from,
		RawMessage:  raw,
	}
}

func (msg *IBCMsgBankSend) ValidateBasic() error {
	if msg.FromAddress.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "empty from address")
	}
	// todo unmarshal to signedIBCMsg
	return nil
}

func (msg *IBCMsgBankSend) GetSignBytes() []byte {
	ntx := *msg
	ntx.SetSignature(nil)
	return util.TxHash(ntx.Bytes())
}

func (msg *IBCMsgBankSend) SetSignature(sig []byte) {
	msg.SetSignature(sig)
}

func (msg *IBCMsgBankSend) GetSignature() []byte{
	return msg.Signature
}

func (msg *IBCMsgBankSend) Bytes() []byte {
	bytes, err := IbcCdc.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}

	return bytes
}

func (msg *IBCMsgBankSend) SetPubKey(pub []byte) {
	msg.PubKey = pub
}

func (msg *IBCMsgBankSend) Route() string {
	return RouterKey
}

func (msg *IBCMsgBankSend) MsgType() string {
	return "IBC_banksend"
}

func (msg *IBCMsgBankSend) GetFromAddress() sdk.AccAddress {
	return msg.FromAddress
}