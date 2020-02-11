package shard

import (
	"encoding/hex"
	"github.com/tanhuiya/ci123chain/pkg/client/helper"
	"github.com/tanhuiya/ci123chain/pkg/cryptosuit"

	"github.com/tanhuiya/ci123chain/pkg/order"
)

func SignAddShardMsg(from string, gas, nonce uint64,t, name string, height int64, priv string) ([]byte, error){

	fromAddr, err := helper.StrToAddress(from)
	if err != nil {
		return nil, err
	}
	tx := order.NewAddShardTx(fromAddr, gas, nonce, t, name, height)
	sid := cryptosuit.NewFabSignIdentity()
	privPub, err := hex.DecodeString(priv)
	if err != nil {
		return nil, err
	}
	pub, err  := sid.GetPubKey(privPub)
	if err != nil {
		return nil, err
	}

	tx.SetPubKey(pub)
	signbyte := tx.GetSignBytes()
	signature, err := sid.Sign(signbyte, privPub)
	tx.SetSignature(signature)
	return tx.Bytes(), nil
}