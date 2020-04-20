package ibc

import (
	"github.com/ci123chain/ci123chain/pkg/client/helper"
	"github.com/ci123chain/ci123chain/pkg/cryptosuit"
	"github.com/ci123chain/ci123chain/pkg/ibc"
	"github.com/ci123chain/ci123chain/pkg/transaction"
)

// 生成 MortgageDone 完成交易

func SignIBCReceiptMsg(from string, raw []byte, gas, nonce uint64, priv []byte) ([]byte, error) {

	tx, err := buildIBCReceiptMsg(from, raw, gas, nonce)
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



func buildIBCReceiptMsg (from string, raw []byte, gas, nonce uint64) (transaction.Transaction, error) {

	fromAddr, err := helper.StrToAddress(from)
	if err != nil {
		return nil, err
	}
	ibcMsg := ibc.NewIBCReceiveReceiptMsg(fromAddr, raw, gas, nonce)
	return ibcMsg, nil
}



