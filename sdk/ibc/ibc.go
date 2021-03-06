package ibc

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/client/helper"
	"github.com/ci123chain/ci123chain/pkg/cryptosuit"
	"github.com/ci123chain/ci123chain/pkg/ibc"
)
func SignIBCTransferMsg(from string, to string, amount uint64, priv []byte) ([]byte, error) {
	tx, err := buildIBCTransferMsg(from, to, amount)
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

func buildIBCTransferMsg(from, to string, amount uint64) (sdk.Msg, error) {
	fromAddr, err := helper.StrToAddress(from)
	if err != nil {
		return nil, err
	}
	toAddr, err := helper.StrToAddress(to)
	if err != nil {
		return nil, err
	}
	ibcMsg := ibc.NewIBCTransfer(fromAddr, toAddr, sdk.NewUInt64Coin(amount))
	return ibcMsg, nil
}

func NewIBCTransferMsg(from, to sdk.AccAddress, amount uint64) []byte{
	msg := ibc.NewIBCTransfer(from, to, sdk.NewUInt64Coin(amount))
	return msg.Bytes()
}

