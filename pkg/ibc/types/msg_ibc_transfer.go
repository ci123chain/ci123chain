package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/util"
)

type IBCTransfer struct {
	FromAddress sdk.AccAddress	`json:"from_address"`
	Signature 	[]byte   		`json:"signature"`
	PubKey		[]byte			`json:"pub_key"`

	ToAddress 	sdk.AccAddress  `json:"to_address"`
	UniqueID 	[]byte 			`json:"unique_id"`
	Coin 	 	sdk.Coin		`json:"coin"`
}

func NewIBCTransferMsg(from, to sdk.AccAddress, amount sdk.Coin) *IBCTransfer {
	return &IBCTransfer{
		FromAddress: from,
		ToAddress:   to,
		UniqueID:    nil,
		Coin:        amount,
	}
}

func (msg *IBCTransfer) ValidateBasic() sdk.Error {
	if msg.ToAddress.Empty() {
		return sdk.ErrInvalidAddress("missing sender address")
	}
	if !msg.Coin.IsValid() {
		return sdk.ErrInvalidCoins("coin is invalid" + msg.Coin.String())
	}
	if msg.FromAddress.Empty() {
		return sdk.ErrInvalidAddress("missing from address")
	}
	return nil
	//return msg.CommonTx.VerifySignature(msg.GetSignBytes(), true)
}

func (msg *IBCTransfer) GetSignBytes() []byte {
	ntx := *msg
	ntx.SetSignature(nil)
	return util.TxHash(ntx.Bytes())
}

func (msg *IBCTransfer) SetSignature(sig []byte) {
	msg.SetSignature(sig)
}

func (msg *IBCTransfer) GetSignature() []byte {
	return msg.Signature
}

func (msg *IBCTransfer) Bytes() []byte {
	bytes, err := IbcCdc.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}

	return bytes
}

func (msg *IBCTransfer) SetPubKey(pub []byte) {
	msg.PubKey = pub
}

func (msg *IBCTransfer) Route() string {
	return RouterKey
}

func (msg *IBCTransfer) MsgType() string {
	return "ibc_transfer"
}

func (msg *IBCTransfer) GetFromAddress() sdk.AccAddress {
	return msg.FromAddress
}