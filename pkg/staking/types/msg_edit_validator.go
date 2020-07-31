package types

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/transaction"
	"github.com/ci123chain/ci123chain/pkg/util"
)

type EditValidatorTx struct {
	transaction.CommonTx
	Description       Description      `json:"description"`
	ValidatorAddress  sdk.AccAddress   `json:"validator_address"`
	CommissionRate    *sdk.Dec          `json:"commission_rate"`
	MinSelfDelegation *sdk.Int          `json:"min_self_delegation"`
}


func NewEditValidatorTx(from sdk.AccAddress, gas ,nonce uint64, desc Description, commissionRate *sdk.Dec,
	minSelfDelegation *sdk.Int) EditValidatorTx {
	return EditValidatorTx{
		CommonTx:          transaction.CommonTx{
			From:      from,
			Nonce:     nonce,
			Gas:       gas,
		},
		Description:       desc,
		ValidatorAddress:  from,
		CommissionRate:    commissionRate,
		MinSelfDelegation: minSelfDelegation,
	}
}


func (tx *EditValidatorTx) ValidateBasic() sdk.Error {

	err := tx.VerifySignature(tx.GetSignBytes(), false)
	if err != nil {
		return ErrCheckParams(DefaultCodespace, err.Error())
	}
	if tx.ValidatorAddress.Empty() {
		return ErrInvalidAddress(DefaultCodespace, fmt.Sprintf("empty validator address"))
	}
	if tx.MinSelfDelegation != nil && tx.MinSelfDelegation.IsPositive() {
		return ErrCheckParams(DefaultCodespace, "invalid minSelfDelegation")
	}
	if tx.CommissionRate != nil {
		if tx.CommissionRate.GT(sdk.OneDec()) || tx.CommissionRate.IsNegative() {
			return ErrCheckParams(DefaultCodespace, "commission rate must be between 0 and 1 (inclusive)")
		}
	}
	if !tx.ValidatorAddress.Equals(tx.From) {
		return ErrInvalidAddress(DefaultCodespace, fmt.Sprintf("expected %s, got %s", tx.From, tx.ValidatorAddress))
	}
	return nil
}


func (tx *EditValidatorTx) GetSignBytes() []byte {
	tmsg := *tx
	tmsg.Signature = nil
	return util.TxHash(tmsg.Bytes())
}


func (tx *EditValidatorTx) SetSignature(sig []byte) {
	tx.Signature = sig
}


func (tx *EditValidatorTx) Bytes() []byte {
	bytes, err := StakingCodec.MarshalBinaryLengthPrefixed(tx)
	if err != nil {
		panic(err)
	}

	return bytes
}


func (tx *EditValidatorTx) SetPubKey(pub []byte) {
	tx.PubKey = pub
}

func (tx *EditValidatorTx) Route() string { return RouteKey }

func (tx *EditValidatorTx) GetGas() uint64 { return tx.Gas }

func (tx *EditValidatorTx) GetNonce() uint64 { return tx.Nonce }

func (tx *EditValidatorTx) GetFromAddress() sdk.AccAddress { return tx.From }