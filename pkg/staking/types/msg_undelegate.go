package types

import (
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/transaction"
	"github.com/ci123chain/ci123chain/pkg/util"
)

type UndelegateTx struct {
	transaction.CommonTx
	DelegatorAddress  types.AccAddress   `json:"delegator_address"`
	ValidatorAddress  types.AccAddress	 `json:"validator_address"`
	Amount            types.Coin		 `json:"amount"`
}

func NewUndelegateTx(from types.AccAddress, gas ,nonce uint64, delegatorAddr types.AccAddress, validatorAddr types.AccAddress,
	amount types.Coin) *UndelegateTx {
	//
	return &UndelegateTx{
		CommonTx:         transaction.CommonTx{
			From: from,
			Gas: 	gas,
			Nonce: nonce,
		},
		DelegatorAddress: delegatorAddr,
		ValidatorAddress: validatorAddr,
		Amount:           amount,
	}
}

func (msg *UndelegateTx) ValidateBasic() types.Error {
	//
	err := msg.VerifySignature(msg.GetSignBytes(), false)
	if err != nil {
		return ErrCheckParams(DefaultCodespace, err.Error())
	}
	if msg.DelegatorAddress.Empty() {
		return ErrInvalidAddress(DefaultCodespace, fmt.Sprintf("empty delegator address"))
	}
	if msg.ValidatorAddress.Empty() {
		return ErrInvalidAddress(DefaultCodespace, fmt.Sprintf("empty validator address"))
	}
	if !msg.Amount.Amount.IsPositive() {
		return types.ErrBadSharesAmount("bad shares amount")
	}
	if !msg.From.Equal(msg.DelegatorAddress) {
		return ErrInvalidAddress(DefaultCodespace, fmt.Sprintf("expected %s, got %s", msg.From.String(), msg.DelegatorAddress.String()))
	}
	return nil
}

func (msg *UndelegateTx) GetSignBytes() []byte {
	tmsg := *msg
	tmsg.Signature = nil
	return util.TxHash(tmsg.Bytes())
}
func (msg *UndelegateTx) SetSignature(sig []byte) {
	msg.Signature = sig
}
func (msg *UndelegateTx) Bytes() []byte {
	bytes, err := StakingCodec.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}

	return bytes
}
func (msg *UndelegateTx) SetPubKey(pub []byte) {
	msg.PubKey = pub
}
func (msg *UndelegateTx) Route() string {return RouteKey}
func (msg *UndelegateTx) GetGas() uint64 { return msg.Gas}

func (msg *UndelegateTx) GetNonce() uint64 { return msg.Nonce}
func (msg *UndelegateTx) GetFromAddress() types.AccAddress { return msg.From}