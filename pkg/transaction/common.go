package transaction

import (
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/cryptosuit"
)

type CommonTx struct {
	Code      uint8
	From      types.AccAddress
	Nonce     uint64
	Gas       uint64
	PubKey 	  []byte
	Signature []byte
}


func (tx CommonTx) ValidateBasic() types.Error {
	if EmptyAddr(tx.From) {
		return ErrInvalidTransfer(DefaultCodespace)
	}
	// TODO Currently we don't support a gas system.
	// if tx.Gas == 0 {
	// 	return ErrInvalidTx(DefaultCodespace, "tx.Gas == 0")
	// }
	if len(tx.Signature) == 0 {
		return ErrInvalidSignature(DefaultCodespace)
	}
	return nil
}

func (tx *CommonTx) SetSignature(sig []byte) {
	tx.Signature = sig
}

func (tx *CommonTx) SetPubKey(pub []byte) {
	tx.PubKey = pub
}


func (tx *CommonTx) VerifySignature(hash []byte, fabricMode bool) types.Error  {

	if fabricMode {
		fab := cryptosuit.NewFabSignIdentity()
		valid, err := fab.Verifier(hash, tx.Signature, tx.PubKey, tx.From.Bytes())
		if !valid || err != nil {
			return ErrInvalidSignature(DefaultCodespace)
		}
	} else {
		eth := cryptosuit.NewETHSignIdentity()
		valid, err := eth.Verifier(hash, tx.Signature, nil, tx.From.Bytes())
		if !valid || err != nil {
			return ErrInvalidSignature(DefaultCodespace)
		}
	}
	return nil
}



func EmptyAddr(addr types.AccAddress) bool {
	return addr == types.AccAddress{}
}
