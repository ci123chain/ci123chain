package sdk

import (
	"gitlab.oneitfarm.com/blockchain/ci123chain/pkg/client/helper"
	"gitlab.oneitfarm.com/blockchain/ci123chain/pkg/cryptosuit"
	"gitlab.oneitfarm.com/blockchain/ci123chain/pkg/transaction"
)

func SignTx(from, to string, amount, gas uint64, priv []byte, isFabric bool) ([]byte, error) {
	tx, err := buildTransferTx(from, to, gas, amount, isFabric)
	if err != nil {
		return nil, err
	}

	sid := cryptosuit.NewFabSignIdentity()
	signature, err := sid.Sign(tx.GetSignBytes(), priv)
	tx.SetSignature(signature)
	return tx.Bytes(), err
}

func Verifier(digest, signature, pubKey []byte, addr []byte) (bool, error) {
	sid := cryptosuit.NewFabSignIdentity()
	return sid.Verifier(digest, signature, pubKey, addr)
}

func buildTransferTx(from, to string, gas, amount uint64, isFabric bool) (transaction.Transaction, error) {
	fromAddr, err := helper.StrToAddress(from)
	if err != nil {
		return nil, err
	}
	toAddr, err := helper.StrToAddress(to)
	if err != nil {
		return nil, err
	}
	nonce, err := transaction.GetNonceByAddress(fromAddr)
	tx := transaction.NewTransferTx(fromAddr, toAddr, gas, nonce, amount, true)
	return tx, nil
}