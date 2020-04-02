package transaction

import (
	"errors"
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/cryptosuit"
	"github.com/tanhuiya/ci123chain/pkg/transaction/types"
)

type CommonTx struct {
	Code      uint8				`json:"code"`
	From      sdk.AccAddress	`json:"from"`
	Nonce     uint64			`json:"nonce"`
	Gas       uint64			`json:"gas"`
	PubKey 	  []byte			`json:"pub_key"`
	Signature []byte			`json:"signature"`
}


func (tx CommonTx) ValidateBasic() sdk.Error {
	if EmptyAddr(tx.From) {
		return types.ErrInvalidTransfer(types.DefaultCodespace, errors.New("empty from address"))
	}
	// TODO Currently we don't support a gas system.
	// if tx.Gas == 0 {
	// 	return ErrInvalidTx(DefaultCodespace, "tx.Gas == 0")
	// }
	if len(tx.Signature) == 0 {
		return types.ErrSignature(types.DefaultCodespace, errors.New("no signature"))
	}
	return nil
}

func (tx *CommonTx) SetSignature(sig []byte) {
	tx.Signature = sig
}

func (tx *CommonTx) SetPubKey(pub []byte) {
	tx.PubKey = pub
}


func (tx *CommonTx) VerifySignature(hash []byte, fabricMode bool) sdk.Error  {

	if fabricMode {
		fab := cryptosuit.NewFabSignIdentity()
		valid, err := fab.Verifier(hash, tx.Signature, tx.PubKey, tx.From.Bytes())
		if !valid || err != nil {
			return types.ErrSignature(types.DefaultCodespace, errors.New("verifier failed"))
		}
	} else {
		eth := cryptosuit.NewETHSignIdentity()
		valid, err := eth.Verifier(hash, tx.Signature, nil, tx.From.Bytes())
		if !valid || err != nil {
			return types.ErrSignature(types.DefaultCodespace, errors.New("verifier failed"))
		}
	}
	return nil
}

func EmptyAddr(addr sdk.AccAddress) bool {
	return addr == sdk.AccAddress{}
}
