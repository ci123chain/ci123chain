package ibc

import (
	"github.com/tanhuiya/ci123chain/pkg/client/helper"
	"github.com/tanhuiya/ci123chain/pkg/cryptosuit"
	"github.com/tanhuiya/ci123chain/pkg/ibc"
	"github.com/tanhuiya/ci123chain/pkg/transaction"
)

// 生成 MortgageDone 完成交易
func SignApplyIBCMsg(from string, uniqueID, observerID []byte, gas uint64, priv []byte, nonce uint64) ([]byte, error) {
	tx, err := buildApplyIBCMsg(from, uniqueID, observerID, gas, nonce)
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


func buildApplyIBCMsg (from string, uniqueID, observerID []byte, gas uint64, nonce uint64) (transaction.Transaction, error) {
	fromAddr, err := helper.StrToAddress(from)
	if err != nil {
		return nil, err
	}

	ibcMsg := ibc.NewApplyIBCTx(fromAddr, gas, nonce, uniqueID, observerID)
	return ibcMsg, nil
}


