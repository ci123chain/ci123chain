package types

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/util"
)

type MsgEditValidator struct {
	FromAddress		  sdk.AccAddress	`json:"from_address"`
	Signature 		  []byte   			`json:"signature"`
	PubKey			  []byte			`json:"pub_key"`

	Description       Description      	`json:"description"`
	ValidatorAddress  sdk.AccAddress   	`json:"validator_address"`
	CommissionRate    *sdk.Dec          `json:"commission_rate"`
	MinSelfDelegation *sdk.Int          `json:"min_self_delegation"`
}

func NewMsgEditValidator(from sdk.AccAddress, desc Description, commissionRate *sdk.Dec,
	minSelfDelegation *sdk.Int) *MsgEditValidator {
	return &MsgEditValidator{
		FromAddress: 		from,
		Description:    	desc,
		ValidatorAddress:  	from,
		CommissionRate:    	commissionRate,
		MinSelfDelegation:	minSelfDelegation,
	}
}

func (tx *MsgEditValidator) ValidateBasic() sdk.Error {
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
	if !tx.ValidatorAddress.Equals(tx.FromAddress) {
		return ErrInvalidAddress(DefaultCodespace, fmt.Sprintf("expected %s, got %s", tx.FromAddress, tx.ValidatorAddress))
	}
	return nil
}

func (tx *MsgEditValidator) GetSignBytes() []byte {
	tmsg := *tx
	tmsg.Signature = nil
	return util.TxHash(tmsg.Bytes())
}

func (tx *MsgEditValidator) SetSignature(sig []byte) {
	tx.Signature = sig
}

func (tx *MsgEditValidator) Bytes() []byte {
	bytes, err := StakingCodec.MarshalBinaryLengthPrefixed(tx)
	if err != nil {
		panic(err)
	}

	return bytes
}

func (tx *MsgEditValidator) SetPubKey(pub []byte) {
	tx.PubKey = pub
}

func (tx *MsgEditValidator) Route() string { return RouteKey }

func (tx *MsgEditValidator) MsgType() string { return "edit-validator" }

func (tx *MsgEditValidator) GetFromAddress() sdk.AccAddress { return tx.FromAddress }

func (tx *MsgEditValidator) GetSignature() []byte {
	return tx.Signature
}