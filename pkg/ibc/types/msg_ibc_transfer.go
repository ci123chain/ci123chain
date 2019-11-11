package types

import (
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/transaction"
	"github.com/tanhuiya/ci123chain/pkg/util"
)

type IBCTransfer struct {
	transaction.CommonTx
	ToAddress 	 sdk.AccAddress `json:"to_address"`
	UniqueID 	 []byte 		`json:"unique_id"`
	Coin 	 sdk.Coin			`json:"coin"`
}

func NewIBCTransferMsg(from, to sdk.AccAddress, amout sdk.Coin, gas uint64, nonce uint64) *IBCTransfer {
	return &IBCTransfer{
		CommonTx:transaction.CommonTx{
			From:  from,
			Gas: 	gas,
			Nonce: nonce,
		},
		ToAddress: to,
		Coin: amout,
	}
}

func (msg *IBCTransfer) ValidateBasic() sdk.Error {
	if err := msg.CommonTx.ValidateBasic(); err != nil {
		return err
	}
	if msg.ToAddress.Empty() {
		return sdk.ErrInvalidAddress("missing sender address")
	}
	if !msg.Coin.IsValid() {
		return sdk.ErrInvalidCoins("coin is invalid" + msg.Coin.String())
	}
	return msg.CommonTx.VerifySignature(msg.GetSignBytes(), true)
}

func NewIBCTransfer(from, to sdk.AccAddress, gas, nonce uint64, coin sdk.Coin, uniqueID []byte) *IBCTransfer {
	msg := &IBCTransfer{
		CommonTx: transaction.CommonTx{
			From: from,
			Nonce: nonce,
			Gas:  gas,
		},
		ToAddress: 	to,
		UniqueID: 	uniqueID,
		Coin: 		coin,
	}
	return msg
}

func (msg *IBCTransfer)GetSignBytes() []byte {
	ntx := *msg
	ntx.SetSignature(nil)
	return util.TxHash(ntx.Bytes())
}


func (msg *IBCTransfer)SetSignature(sig []byte) {
	msg.CommonTx.SetSignature(sig)
}

func (msg *IBCTransfer)Bytes() []byte {
	bytes, err := IbcCdc.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}

	return bytes
}

func (msg *IBCTransfer)SetPubKey(pub []byte) {
	msg.CommonTx.PubKey = pub
}

func (msg *IBCTransfer) Route() string {
	return RouterKey
}
