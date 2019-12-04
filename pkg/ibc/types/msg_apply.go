package types

import (
	"encoding/hex"
	"github.com/tanhuiya/ci123chain/pkg/transaction"
	"github.com/tanhuiya/ci123chain/pkg/transfer"
	"github.com/tanhuiya/ci123chain/pkg/util"
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
)

type ApplyIBCTx struct {
	// 跨链交易ID
	transaction.CommonTx
	UniqueID []byte		`json:"unique_id"`
	ObserverID []byte	`json:"observer_id"`
}

func NewApplyIBCTx(from sdk.AccAddress, gas ,nonce uint64, uniqueID, observerID []byte) *ApplyIBCTx {
	return &ApplyIBCTx{
		CommonTx: transaction.CommonTx{
			From: from,
			Gas: 	gas,
			Nonce: nonce,
		},
		UniqueID: uniqueID,
		ObserverID: observerID,
	}
}

func (msg *ApplyIBCTx) ValidateBasic() sdk.Error {
	if err := msg.CommonTx.ValidateBasic(); err != nil {
		return err
	}
	if len(msg.UniqueID) < 1 {
		return transfer.ErrCheckParams(DefaultCodespace, "UniqueID is invalid " + hex.EncodeToString(msg.UniqueID))
	}
	if len(msg.ObserverID) < 1 {
		return transfer.ErrCheckParams(DefaultCodespace, "ObserverID is invalid " + hex.EncodeToString(msg.ObserverID))
	}
	return msg.CommonTx.VerifySignature(msg.GetSignBytes(), true)
}


func (msg *ApplyIBCTx)GetSignBytes() []byte {
	ntx := *msg
	ntx.SetSignature(nil)
	return util.TxHash(ntx.Bytes())
}


func (msg *ApplyIBCTx)SetSignature(sig []byte) {
	msg.CommonTx.SetSignature(sig)
}

func (msg *ApplyIBCTx)Bytes() []byte {
	bytes, err := IbcCdc.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}

	return bytes
}

func (msg *ApplyIBCTx) SetPubKey(pub []byte) {
	msg.CommonTx.PubKey = pub
}

func (msg *ApplyIBCTx) Route() string {
	return RouterKey
}

func (msg *ApplyIBCTx) GetGas() uint64 {
	return msg.CommonTx.Gas
}

func (msg *ApplyIBCTx) GetNonce() uint64 {
	return msg.CommonTx.Nonce
}

func (msg *ApplyIBCTx) GetFromAddress() sdk.AccAddress {
	return msg.CommonTx.From
}