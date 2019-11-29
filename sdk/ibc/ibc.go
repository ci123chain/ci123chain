package ibc

import (
	"github.com/spf13/viper"
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/app"
	"github.com/tanhuiya/ci123chain/pkg/client"
	"github.com/tanhuiya/ci123chain/pkg/client/helper"
	"github.com/tanhuiya/ci123chain/pkg/cryptosuit"
	"github.com/tanhuiya/ci123chain/pkg/ibc"
	"github.com/tanhuiya/ci123chain/pkg/transaction"
)
var cdc = app.MakeCodec()

// 生成 MortgageDone 完成交易
func SignIBCTransferMsg(from string, to string, amount, gas uint64, priv []byte, node string) ([]byte, error) {
	tx, err := buildIBCTransferMsg(from, to, amount, gas, node)
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


func buildIBCTransferMsg (from, to string, amount, gas uint64, node string) (transaction.Transaction, error) {
	fromAddr, err := helper.StrToAddress(from)
	if err != nil {
		return nil, err
	}
	toAddr, err := helper.StrToAddress(to)
	if err != nil {
		return nil, err
	}
	viper.Set("node", "tcp://" + node)
	viper.Set("address", fromAddr)
	ctx, err := client.NewClientContextFromViper(cdc)
	if err != nil {
		return nil,err
	}
	nonce, err := ctx.GetNonceByAddress(fromAddr)
	ibcMsg := ibc.NewIBCTransfer(fromAddr, toAddr, sdk.Coin(amount),  gas, nonce)
	return ibcMsg, nil
}

