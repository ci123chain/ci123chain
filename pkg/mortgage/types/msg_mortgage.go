package types

import (
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/transaction"
	"github.com/tanhuiya/ci123chain/pkg/util"
)

const (
	StateMortgaged = "StateMortgaged"
	StateSuccess = "StateSuccess"
	StateCancel = "StateCancel"
)

var _ transaction.Transaction = (*MsgMortgage)(nil)

type MsgMortgage struct {
	transaction.CommonTx
	//FromAddress  sdk.AccAddress `json:"from_address"`
	ToAddress 	 sdk.AccAddress `json:"to_address"`
	UniqueID 	 []byte 		`json:"unique_id"`
	Coin 	 sdk.Coin			`json:"coin"`
}

func (msg *MsgMortgage) ValidateBasic() sdk.Error {
	if err := msg.CommonTx.ValidateBasic(); err != nil {
		return err
	}
	if msg.ToAddress.Empty() {
		return sdk.ErrInvalidAddress("missing sender address")
	}
	if len(msg.UniqueID) < 1 {
		return sdk.ErrInternal("param mortgageRecord missing")
	}
	if !msg.Coin.IsValid() {
		return sdk.ErrInvalidCoins("coin is invalid" + msg.Coin.String())
	}
	return msg.CommonTx.VerifySignature(msg.GetSignBytes(), true)
}

func NewMsgMortgage(from, to sdk.AccAddress, gas, nonce uint64, coin sdk.Coin, uniqueID []byte) *MsgMortgage {
	msg := &MsgMortgage{
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
func (msg *MsgMortgage)GetSignBytes() []byte {
	ntx := *msg
	ntx.SetSignature(nil)
	return util.TxHash(ntx.Bytes())
}


func (msg *MsgMortgage)SetSignature(sig []byte) {
	msg.CommonTx.SetSignature(sig)
}

func (msg *MsgMortgage)Bytes() []byte {
	bytes, err := MortgageCdc.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}

	return bytes
}

func (msg *MsgMortgage)SetPubKey(pub []byte) {
	msg.CommonTx.PubKey = pub
}

func (msg *MsgMortgage) Route() string {
	return RouterKey
}

func (msg *MsgMortgage) GetGas() uint64 {
	return msg.CommonTx.Gas
}

func (msg *MsgMortgage) GetNonce() uint64 {
	return msg.CommonTx.Nonce
}

func (msg *MsgMortgage) GetFromAddress() sdk.AccAddress {
	return msg.CommonTx.From
}

type Mortgage struct {
	MsgMortgage

	State  string `json:"state"`
}