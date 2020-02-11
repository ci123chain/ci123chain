package transfer

import (
	"encoding/hex"
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/client/helper"
	"github.com/tanhuiya/ci123chain/pkg/cryptosuit"
	"github.com/tanhuiya/ci123chain/pkg/transfer"
)

func SignTransferMsg(from, to string, amount, gas, nonce uint64, priv string, isfabric bool) ([]byte, error) {

	var signature []byte
	fromAddr, err := helper.StrToAddress(from)
	if err != nil {
		return nil, err
	}
	toAddr, err := helper.StrToAddress(to)
	if err != nil {
		return nil, err
	}
	privPub, err := hex.DecodeString(priv)
	if err != nil {
		return nil, err
	}
	tx := transfer.NewTransferTx(fromAddr, toAddr, gas, nonce, sdk.NewUInt64Coin(amount), isfabric)

	if isfabric {
		fab := cryptosuit.NewFabSignIdentity()
		pubkey, err := fab.GetPubKey(privPub)
		if err != nil {
			return nil, err
		}
		tx.SetPubKey(pubkey)
		signature, err = fab.Sign(tx.GetSignBytes(), privPub)
		if err != nil {
			return nil, err
		}
	} else {
		eth := cryptosuit.NewETHSignIdentity()
		signature, err = eth.Sign(tx.GetSignBytes(), privPub)
		if err != nil {
			return nil, err
		}
	}
	tx.SetSignature(signature)
	return tx.Bytes(), nil
}
