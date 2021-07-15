package types

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	types2 "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/ci123chain/ci123chain/pkg/cryptosuit"
	"github.com/ci123chain/ci123chain/pkg/util"
	"github.com/ethereum/go-ethereum/crypto"
)

type CommonTx struct {
	From      types2.AccAddress `json:"from"`
	Nonce     uint64            `json:"nonce"`
	Gas       uint64            `json:"gas"`
	Msgs      []types2.Msg      `json:"msgs"`
	PubKey    []byte            `json:"pub_key"`
	Signature []byte            `json:"signature"`
}

func NewCommonTx(from types2.AccAddress, nonce, gas uint64, msgs []types2.Msg) *CommonTx {
	return &CommonTx{
		From:      from,
		Nonce:     nonce,
		Gas:       gas,
		Msgs:      msgs,
	}
}

func (tx CommonTx) ValidateBasic() error {
	if tx.From.Empty() {
		return ErrInvalidParam("empty from address")
	}
	// TODO Currently we don't support a gas system.
	if len(tx.Msgs) == 0 {
		return ErrInvalidParam("empty messagees")
	}
	if len(tx.Signature) == 0 {
		return ErrInvalidParam("message with no signature")

	}
	return nil
}

func (tx *CommonTx) SetSignature(sig []byte) {
	tx.Signature = sig
}

func (tx *CommonTx) GetSignature() []byte {
	return tx.Signature
}

func (tx *CommonTx) SetPubKey(pub []byte) {
	tx.PubKey = pub
}

func (msg *CommonTx) GetSignBytes() []byte{
	ntx := *msg
	ntx.SetSignature(nil)
	return util.TxHash(ntx.Bytes())
}

func (msg *CommonTx) Bytes() []byte {
	bytes, err := GetCodec().MarshalBinaryBare(msg)
	if err != nil {
		panic(err)
	}
	return bytes
}

func (msg *CommonTx) GetGas() uint64 {
	return msg.Gas
}

func (msg *CommonTx) GetNonce() uint64 {
	return msg.Nonce
}

func (msg *CommonTx) GetMsgs() []types2.Msg{
	return msg.Msgs
}

func (msg *CommonTx) GetFromAddress() types2.AccAddress{
	return msg.From
}


func SignCommonTx(from types2.AccAddress, nonce, gas uint64, msgs []types2.Msg, priv string, cdc *codec.Codec) ([]byte, error){
	tx := NewCommonTx(from, nonce, gas, msgs)
	var signature []byte
	privPub, err := hex.DecodeString(priv)
	eth := cryptosuit.NewETHSignIdentity()
	if !IsValidPrivateKey(from, privPub){
		return nil, ErrInvalidParam("invalid private_key, the private key does not match the from account")
	}
	signature, err = eth.Sign(tx.GetSignBytes(), privPub)
	if err != nil {
		return nil, err
	}
	tx.SetSignature(signature)
	return cdc.MarshalBinaryBare(tx)
}

// DefaultTxDecoder logic for standard transfer decoding
func DefaultTxDecoder(cdc *codec.Codec) types2.TxDecoder {
	return func(txBytes []byte) (types2.Tx, error) {
		var transfer *CommonTx
		err := codec.GetLegacyAminoByCodec(cdc).UnmarshalBinaryBare(txBytes, &transfer)
		if err != nil {
			var pbTx PbTx
			err = GetEncodingConfig().Marshaler.UnmarshalBinaryBare(txBytes, &pbTx)
			if err != nil {
				var ethTx *MsgEthereumTx
				err := cdc.UnmarshalBinaryBare(txBytes, &ethTx)
				if err != nil {
					return nil, sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("decode msg failed: %v", err.Error()))
				}
				return ethTx, nil
			}
			return &pbTx, nil
		}
		return transfer, nil
	}
}

func IsValidPrivateKey(from types2.AccAddress,identity []byte) bool {
	key := crypto.ToECDSAUnsafe(identity)
	by1 := crypto.PubkeyToAddress(key.PublicKey).Bytes()
	by2 := from.Bytes()
	if bytes.Equal(by1, by2) {
		return true
	}else {
		return false
	}
}