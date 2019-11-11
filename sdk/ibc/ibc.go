package ibc

import (
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/client/helper"
	"github.com/tanhuiya/ci123chain/pkg/cryptosuit"
	"github.com/tanhuiya/ci123chain/pkg/ibc"
	"github.com/tanhuiya/ci123chain/pkg/transaction"
)

// 生成 MortgageDone 完成交易
func SignIBCTransferMsg(from string, to string, amount, gas uint64, priv []byte) ([]byte, error) {
	tx, err := buildIBCTransferMsg(from, to, amount, gas)
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


func buildIBCTransferMsg (from, to string, amount, gas uint64) (transaction.Transaction, error) {
	fromAddr, err := helper.StrToAddress(from)
	if err != nil {
		return nil, err
	}
	toAddr, err := helper.StrToAddress(to)
	if err != nil {
		return nil, err
	}
	nonce, err := transaction.GetNonceByAddress(fromAddr)
	ibcMsg := ibc.NewIBCTransfer(fromAddr, toAddr, sdk.Coin(amount),  gas, nonce)
	return ibcMsg, nil
}

