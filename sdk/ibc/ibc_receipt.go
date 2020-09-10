package ibc

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/client/helper"
	"github.com/ci123chain/ci123chain/pkg/cryptosuit"
	"github.com/ci123chain/ci123chain/pkg/ibc"
)

// 生成 MortgageDone 完成交易

func SignIBCReceiptMsg(from string, raw, priv []byte) ([]byte, error) {

	tx, err := buildIBCReceiptMsg(from, raw)
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



func buildIBCReceiptMsg (from string, raw []byte) (sdk.Msg, error) {

	fromAddr, err := helper.StrToAddress(from)
	if err != nil {
		return nil, err
	}
	ibcMsg := ibc.NewIBCReceiveReceiptMsg(fromAddr, raw)
	return ibcMsg, nil
}



