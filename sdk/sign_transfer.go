package sdk

import (
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/client/helper"
	"github.com/tanhuiya/ci123chain/pkg/cryptosuit"
	"github.com/tanhuiya/ci123chain/pkg/transfer"
)

// 签名消息
// 返回 []byte, 需要 转为 hex 类型后广播
func SignTx(from, to string, amount, gas uint64, priv []byte, isFabric bool) ([]byte, error) {
	tx, err := buildTransferTx(from, to, gas, amount, isFabric)
	if err != nil {
		return nil, err
	}

	sid := cryptosuit.NewFabSignIdentity()
	pub, err  := sid.GetPubKey(priv)
	if err != nil {
		return nil, err
	}
	tx.SetPubKey(pub)
	signature, err := sid.Sign(tx.GetSignBytes(), priv)
	tx.SetSignature(signature)
	return tx.Bytes(), err
}

func Verifier(digest, signature, pubKey []byte, addr []byte) (bool, error) {
	sid := cryptosuit.NewFabSignIdentity()
	return sid.Verifier(digest, signature, pubKey, addr)
}

func buildTransferTx(from, to string, gas, amount uint64, isFabric bool) (transfer.Transaction, error) {
	fromAddr, err := helper.StrToAddress(from)
	if err != nil {
		return nil, err
	}
	toAddr, err := helper.StrToAddress(to)
	if err != nil {
		return nil, err
	}
	nonce, err := transfer.GetNonceByAddress(fromAddr)
	tx := transfer.NewTransferTx(fromAddr, toAddr, gas, nonce, types.Coin(amount), isFabric)
	return tx, nil
}