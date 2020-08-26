package transfer

import (
	"errors"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/client/helper"
	"github.com/ci123chain/ci123chain/pkg/cryptosuit"
	"github.com/ci123chain/ci123chain/pkg/transaction"
	"github.com/ci123chain/ci123chain/pkg/transfer/types"
	"github.com/ci123chain/ci123chain/pkg/util"
)

const RouteKey = "Transfer"

func NewTransferTx(from, to sdk.AccAddress, gas, nonce uint64, amount sdk.Coin, isFabric bool ) transaction.Transaction {
	tx := &TransferTx{
		CommonTx: transaction.CommonTx{
			From: from,
			Gas:  gas,
			Nonce:nonce,
		},
		To: 		to,
		Amount: 	amount,
		FabricMode: isFabric,
	}
	return tx
}

func SignTransferTx(from string, to string, amount, gas, nonce uint64, priv []byte) ([]byte, error) {
	fromAddr, err := helper.StrToAddress(from)
	if err != nil {
		return nil, err
	}
	toAddr, err := helper.StrToAddress(to)
	if err != nil {
		return nil, err
	}
	tx := NewTransferTx(fromAddr, toAddr, gas, nonce, sdk.NewUInt64Coin(amount), true)
	sid := cryptosuit.NewFabSignIdentity()
	pub, err  := sid.GetPubKey(priv)

	tx.SetPubKey(pub)
	signbyte := tx.GetSignBytes()
	signature, err := sid.Sign(signbyte, priv)
	tx.SetSignature(signature)
	return tx.Bytes(), nil
}


type TransferTx struct {
	transaction.CommonTx
	To     sdk.AccAddress   	`json:"to"`
	Amount sdk.Coin         	`json:"amount"`
	FabricMode bool         	`json:"fabric_mode"`
}

//func DecodeTransferTx(b []byte) (*TransferTx, error) {
//	var transfer TransferTx
//	err := transferCdc.UnmarshalBinaryLengthPrefixed(b, &transfer)
//	if err != nil {
//		return nil, err
//	}
//	return &transfer, nil
//
//	//tx := new(TransferTx)
//	//return tx, rlp.DecodeBytes(b, tx)
//}

func (tx *TransferTx) SetPubKey(pub []byte) {
	tx.CommonTx.PubKey = pub
}

func (tx *TransferTx) SetSignature(sig []byte) {
	tx.CommonTx.SetSignature(sig)
}

func (tx *TransferTx) GetSignature() []byte{
	return tx.CommonTx.GetSignature()
}


func (tx *TransferTx) ValidateBasic() sdk.Error {
	if err := tx.CommonTx.ValidateBasic(); err != nil {
		return err
	}
	if tx.Amount.IsEqual(sdk.NewCoin(sdk.NewInt(0)))  {
		return types.ErrBadAmount(types.DefaultCodespace, errors.New("amount = 0"))
	}
	if transaction.EmptyAddr(tx.To) {
		return types.ErrBadReceiver(types.DefaultCodespace, errors.New("empty to address"))
	}
	return nil
	//return tx.VerifySignature(tx.GetSignBytes(), tx.FabricMode)
}

func (tx *TransferTx) Route() string {
	return RouteKey
}

func (tx *TransferTx) GetSignBytes() []byte {
	ntx := *tx
	ntx.SetSignature(nil)
	return util.TxHash(ntx.Bytes())
}

func (tx *TransferTx) Bytes() []byte {

	bytes, err := transferCdc.MarshalBinaryLengthPrefixed(tx)
	if err != nil {
		panic(err)
	}

	return bytes
}

func (tx *TransferTx) GetGas() uint64 {
	return tx.Gas
}

func (tx *TransferTx) GetNonce() uint64 {
	return tx.Nonce
}

func (tx *TransferTx) GetFromAddress() sdk.AccAddress {
	return tx.From
}
