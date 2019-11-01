package transaction

import (
	"github.com/tanhuiya/ci123chain/pkg/cryptosuit"
	"io"

	"github.com/tanhuiya/ci123chain/pkg/abci/types"
	perrors "github.com/tanhuiya/ci123chain/pkg/error"
	"github.com/ethereum/go-ethereum/rlp"
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
	if isEmptyAddr(tx.From) {
		return ErrInvalidTx(DefaultCodespace, "tx.From == nil")
	}
	// TODO Currently we don't support a gas system.
	// if tx.Gas == 0 {
	// 	return ErrInvalidTx(DefaultCodespace, "tx.Gas == 0")
	// }
	if len(tx.Signature) == 0 {
		return ErrInvalidTx(DefaultCodespace, "len(tx.Signature) == 0")
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
			return perrors.ErrInvalidSignature(perrors.DefaultCodespace, err.Error())
		}
	} else {
		eth := cryptosuit.NewETHSignIdentity()
		valid, err := eth.Verifier(hash, tx.Signature, nil, tx.From.Bytes())
		if !valid || err != nil {
			return perrors.ErrInvalidSignature(perrors.DefaultCodespace, err.Error())
		}
	}
	return nil
}


func (tx CommonTx) EncodeRLP(w io.Writer) error {
	if err := tx.EncodeNoSig(w); err != nil {
		return err
	}
	if err := rlp.Encode(w, tx.Signature); err != nil {
		return err
	}
	return nil
}

func (tx CommonTx) EncodeNoSig(w io.Writer) error {
	if err := rlp.Encode(w, tx.Code); err != nil {
		return err
	}
	if err := rlp.Encode(w, tx.From); err != nil {
		return err
	}
	if err := rlp.Encode(w, tx.Nonce); err != nil {
		return err
	}
	if err := rlp.Encode(w, tx.Gas); err != nil {
		return err
	}
	if err := rlp.Encode(w, tx.PubKey); err != nil {
		return err
	}
	return nil
}

func (tx *CommonTx) DecodeRLP(s *rlp.Stream) error {

	if err := s.Decode(&tx.Code); err != nil {
		return err
	}
	if err := s.Decode(&tx.From); err != nil {
		return err
	}
	if err := s.Decode(&tx.Nonce); err != nil {
		return err
	}
	if err := s.Decode(&tx.Gas); err != nil {
		return err
	}
	if err := s.Decode(&tx.PubKey); err != nil {
		return err
	}
	b, err := s.Bytes()
	if err != nil {
		return err
	}
	tx.Signature = b
	return nil
}
