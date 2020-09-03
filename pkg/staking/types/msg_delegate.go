package types

import (
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/transaction"
	"github.com/ci123chain/ci123chain/pkg/util"
)

type DelegateTx struct {
	transaction.CommonTx
	DelegatorAddress  types.AccAddress    `json:"delegator_address"`
	ValidatorAddress  types.AccAddress    `json:"validator_address"`
	Amount            types.Coin		  `json:"amount"`
}

func NewDelegateTx(from types.AccAddress, gas ,nonce uint64, delegatorAddr types.AccAddress, validatorAddr types.AccAddress,
	amount types.Coin) DelegateTx {

	return DelegateTx{
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

func (msg *DelegateTx) ValidateBasic() types.Error {
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
		return ErrCheckParams(DefaultCodespace, "amount should be positive")
	}
	if !msg.From.Equal(msg.DelegatorAddress) {
		return ErrInvalidAddress(DefaultCodespace, fmt.Sprintf("expected %s, got %s", msg.From.String(), msg.DelegatorAddress.String()))
	}
	return nil
}

func (msg *DelegateTx) GetSignBytes() []byte {
	tmsg := *msg
	tmsg.Signature = nil
	signBytes := tmsg.Bytes()
	return util.TxHash(signBytes)
}
func (msg *DelegateTx) SetSignature(sig []byte) {
	msg.Signature = sig
}
func (msg *DelegateTx) Bytes() []byte {
	bytes, err := StakingCodec.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}

	return bytes
}
func (msg *DelegateTx) SetPubKey(pub []byte) {
	msg.CommonTx.PubKey = pub
}
func (msg *DelegateTx) Route() string {return RouteKey}
func (msg *DelegateTx) GetGas() uint64 { return msg.Gas}

func (msg *DelegateTx) GetNonce() uint64 { return msg.Nonce}
func (msg *DelegateTx) GetFromAddress() types.AccAddress { return msg.From}