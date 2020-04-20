package ibc

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/app"
	"github.com/ci123chain/ci123chain/pkg/client/helper"
	"github.com/ci123chain/ci123chain/pkg/cryptosuit"
	"github.com/ci123chain/ci123chain/pkg/ibc"
	"github.com/ci123chain/ci123chain/pkg/transaction"
)
var cdc = app.MakeCodec()

// 生成 MortgageDone 完成交易

func SignIBCTransferMsg(from string, to string, amount, gas, nonce uint64, priv []byte) ([]byte, error) {
	tx, err := buildIBCTransferMsg(from, to, amount, gas, nonce)
	if err != nil {
		return nil, err
	}
	sid := cryptosuit.NewFabSignIdentity()
	pub, err  := sid.GetPubKey(priv)

	tx.SetPubKey(pub)
	signbyte := tx.GetSignBytes()
	signature, err := sid.Sign(signbyte, priv)
	tx.SetSignature(signature)
	return tx.Bytes(), nil
}



func buildIBCTransferMsg (from, to string, amount, gas, nonce uint64) (transaction.Transaction, error) {
	fromAddr, err := helper.StrToAddress(from)
	if err != nil {
		return nil, err
	}
	toAddr, err := helper.StrToAddress(to)
	if err != nil {
		return nil, err
	}
	ibcMsg := ibc.NewIBCTransfer(fromAddr, toAddr, sdk.NewUInt64Coin(amount),  gas, nonce)
	return ibcMsg, nil
}

