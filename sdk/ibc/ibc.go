package ibc

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/client/helper"
	"github.com/ci123chain/ci123chain/pkg/cryptosuit"
	"github.com/ci123chain/ci123chain/pkg/ibc"
)
func SignIBCTransferMsg(from string, to string, amount uint64, priv []byte, denom string, gas,  nonce uint64) ([]byte, error) {
	tx, err := buildIBCTransferMsg(from, to, amount, denom)
	if err != nil {
		return nil, err
	}
	sid := cryptosuit.NewFabSignIdentity()
	pub, err  := sid.GetPubKey(priv)

	signbyte := tx.Bytes()
	signature, err := sid.Sign(signbyte, priv)
	msg := types.CommonTx{sdk.HexToAddress(from), nonce, gas,  []sdk.Msg{tx}, pub, signature}
	return msg.Bytes(), nil
}

func buildIBCTransferMsg(from, to string, amount uint64, denom string) (sdk.Msg, error) {
	fromAddr, err := helper.StrToAddress(from)
	if err != nil {
		return nil, err
	}
	toAddr, err := helper.StrToAddress(to)
	if err != nil {
		return nil, err
	}
	ibcMsg := ibc.NewIBCTransfer(fromAddr, toAddr, sdk.NewUInt64Coin(denom, amount))
	return ibcMsg, nil
}

func NewIBCTransferMsg(from, to sdk.AccAddress, amount uint64, denom string) []byte{
	msg := ibc.NewIBCTransfer(from, to, sdk.NewUInt64Coin(denom, amount))
	return msg.Bytes()
}

