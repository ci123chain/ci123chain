package ibc

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/client/helper"
	"github.com/ci123chain/ci123chain/pkg/cryptosuit"
	"github.com/ci123chain/ci123chain/pkg/ibc"
)

// 生成 MortgageDone 完成交易

func SignApplyIBCMsg(from string, uniqueID, observerID, priv []byte) ([]byte, error) {

	tx, err := buildApplyIBCMsg(from, uniqueID, observerID)
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



func buildApplyIBCMsg (from string, uniqueID, observerID []byte) (sdk.Msg, error) {
	fromAddr, err := helper.StrToAddress(from)
	if err != nil {
		return nil, err
	}
	ibcMsg := ibc.NewApplyIBCTx(fromAddr, uniqueID, observerID)
	return ibcMsg, nil
}


