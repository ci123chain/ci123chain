package types

import (
	"encoding/hex"
	"errors"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	codectypes "github.com/ci123chain/ci123chain/pkg/abci/codec/types"
	types2 "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/ci123chain/ci123chain/pkg/cryptosuit"
	"github.com/ci123chain/ci123chain/pkg/util"
)


func NewPbTx(from types2.AccAddress, nonce, gas uint64, msgs []types2.PbMsg) *PbTx {

	return &PbTx{
		From:      from.Bytes(),
		Nonce:     nonce,
		Gas:       gas,
		Msgs:      PackPbMsgs(msgs),
	}
}

func PackPbMsgs(msgs []types2.PbMsg) []*codectypes.Any {
	ms := []*codectypes.Any{}
	for _, v := range msgs {
		any, err := codectypes.NewAnyWithValue(v)
		if err != nil {
			panic(err)
		}
		ms = append(ms, any)
	}
	return ms
}

func (tx PbTx) ValidateBasic() error {
	if len(tx.From) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "empty from address")
	}
	// TODO Currently we don't support a gas system.
	if len(tx.Msgs) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrParams, "empty messagees")
	}
	if len(tx.Signature) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrNoSignatures, "message with no signature")
	}
	return nil
}

func (tx *PbTx) SetSignature(sig []byte) {
	tx.Signature = sig
}

func (tx *PbTx) SetPubKey(pub []byte) {
	tx.PubKey = pub
}

func (msg *PbTx) GetSignBytes() []byte{
	ntx := *msg
	ntx.SetSignature(nil)
	return util.TxHash(ntx.Bytes())
}

func (msg *PbTx) Bytes() []byte {

	bytes, err := GetEncodingConfig().Marshaler.MarshalBinaryBare(msg)
	if err != nil {
		panic(err)
	}
	return bytes
}

func (msg *PbTx) GetFromAddress() types2.AccAddress{
	return types2.ToAccAddress(msg.From)
}

func (m *PbTx) GetMsgs() []types2.Msg {
	if m != nil {
		var msgs []types2.Msg
		for _, v := range m.Msgs{
			var msg types2.Msg
			err := GetEncodingConfig().InterfaceRegistry.UnpackAny(v, &msg)
			if err != nil {
				panic(err)
			}
			msgs = append(msgs, msg)
		}
		return msgs
	}
	return nil
}

func SignPbTx(from types2.AccAddress, nonce, gas uint64, msgs []types2.PbMsg, priv string, cdc codec.BinaryMarshaler) ([]byte, error){
	tx := NewPbTx(from, nonce, gas, msgs)
	var signature []byte
	privPub, err := hex.DecodeString(priv)
	eth := cryptosuit.NewETHSignIdentity()
	if !IsValidPrivateKey(from, privPub){
		return nil, errors.New("invalid private_key, the private key does not match the from account")
	}
	signature, err = eth.Sign(tx.GetSignBytes(), privPub)
	if err != nil {
		return nil, err
	}
	tx.SetSignature(signature)
	return cdc.MarshalBinaryBare(tx)
}


func (tx *PbTx) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	for _, m := range tx.Msgs {
		err := codectypes.UnpackInterfaces(m, unpacker)
		if err != nil {
			return err
		}
	}
	return nil
}