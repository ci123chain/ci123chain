package types

import (
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/order/keeper"
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



func (msg *UpgradeTx) ValidateBasic() sdk.Error{
	//return msg.CommonTx.VerifySignature(msg.Bytes(), msg.FabricMode)
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
	bytes, err := keeper.ModuleCdc.MarshalBinaryLengthPrefixed(msg)
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