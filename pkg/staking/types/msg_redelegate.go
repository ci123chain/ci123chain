package types

import (
	"encoding/json"
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/transaction"
)

type RedelegateTx struct {
	transaction.CommonTx
	DelegatorAddress     types.AccAddress    `json:"delegator_address"`
	ValidatorSrcAddress  types.AccAddress	 `json:"validator_src_address"`
	ValidatorDstAddress  types.AccAddress	 `json:"validator_dst_address"`
	Amount               types.Coin	 		 `json:"amount"`
}

func NewRedelegateTx(from types.AccAddress, gas ,nonce uint64, delegateAddr types.AccAddress, validatorSrcAddr,
	validatorDstAddr types.AccAddress, amount types.Coin) *RedelegateTx {
	//
	return &RedelegateTx{
		CommonTx:            transaction.CommonTx{
			From: from,
			Gas: 	gas,
			Nonce: nonce,
		},
		DelegatorAddress:    delegateAddr,
		ValidatorSrcAddress: validatorSrcAddr,
		ValidatorDstAddress: validatorDstAddr,
		Amount:              amount,
	}
}

func (msg *RedelegateTx) ValidateBasic() types.Error {
	//
	if msg.DelegatorAddress.Empty() {
		return types.ErrEmptyDelegatorAddr("empty delegator address")
	}
	if msg.ValidatorSrcAddress.Empty() {
		return types.ErrEmptyValidatorAddr("empty validator address")
	}
	if msg.ValidatorDstAddress.Empty() {
		return types.ErrEmptyValidatorAddr("empty validator address")
	}
	if !msg.Amount.Amount.IsPositive() {
		return types.ErrBadSharesAmount("bad shares amount")
	}
	return nil
}

func (msg *RedelegateTx) GetSignBytes() []byte {
	tmsg := *msg
	tmsg.Signature = nil
	signBytes := tmsg.Bytes()
	return signBytes
}
func (msg *RedelegateTx) SetSignature(sig []byte) {
	msg.Signature = sig
}
func (msg *RedelegateTx) Bytes() []byte {
	bytes, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return bytes
}
func (msg *RedelegateTx) SetPubKey(pub []byte) {
	msg.PubKey = pub
}
func (msg *RedelegateTx) Route() string {return RouteKey}
func (msg *RedelegateTx) GetGas() uint64 { return msg.Gas}

func (msg *RedelegateTx) GetNonce() uint64 { return msg.Nonce}
func (msg *RedelegateTx) GetFromAddress() types.AccAddress { return msg.From}