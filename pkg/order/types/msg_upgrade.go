package types


import (
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/client/helper"
	"github.com/tanhuiya/ci123chain/pkg/cryptosuit"
	"github.com/tanhuiya/ci123chain/pkg/transaction"
	"github.com/tanhuiya/ci123chain/pkg/util"
)

type UpgradeTx struct {
	transaction.CommonTx
	Type      string   `json:"type"`
	Height    int64    `json:"height"`
	Name      string   `json:"name"`
}

func NewUpgradeTx(from sdk.AccAddress, gas ,nonce uint64, t, name string, height int64) UpgradeTx{

	return UpgradeTx{
		CommonTx: transaction.CommonTx{
			From: from,
			Gas: 	gas,
			Nonce: nonce,
		},
		Type:t,
		Height:height,
		Name:name,
	}
}

func SignUpgradeTx(from string, gas, nonce uint64,t, name string, height int64, priv []byte) ([]byte, error){

	fromAddr, err := helper.StrToAddress(from)
	if err != nil {
		return nil, err
	}
	tx := NewUpgradeTx(fromAddr, gas, nonce, t, name, height)
	sid := cryptosuit.NewFabSignIdentity()
	pub, err  := sid.GetPubKey(priv)

	tx.SetPubKey(pub)
	signbyte := tx.GetSignBytes()
	signature, err := sid.Sign(signbyte, priv)
	tx.SetSignature(signature)
	return tx.Bytes(), nil
}


func (msg *UpgradeTx) ValidateBasic() sdk.Error{
	return nil
}

func (msg *UpgradeTx) SetPubKey(pub []byte) {
	msg.CommonTx.PubKey = pub
}

func (msg *UpgradeTx) SetSignature(sig []byte) {
	msg.CommonTx.Signature = sig
}

func (msg *UpgradeTx) Route() string {
	return RouteKey
}

func (msg *UpgradeTx) GetSignBytes() []byte{
	ntx := *msg
	ntx.SetSignature(nil)
	return util.TxHash(ntx.Bytes())
}

func (msg *UpgradeTx)Bytes() []byte {
	bytes, err := OrCdc.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}

	return bytes
}

func (msg *UpgradeTx) GetNonce() uint64 {
	return msg.CommonTx.Nonce
}

func (msg *UpgradeTx) GetGas() uint64 {
	return msg.CommonTx.Gas
}

func (msg *UpgradeTx) GetFromAddress() sdk.AccAddress{
	return msg.CommonTx.From
}