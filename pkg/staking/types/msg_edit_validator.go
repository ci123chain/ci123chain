package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
)

type MsgEditValidator struct {
	FromAddress		  sdk.AccAddress	`json:"from_address"`
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

func (tx *MsgEditValidator) ValidateBasic() error {
	if tx.ValidatorAddress.Empty() {
		return ErrInvalidParam
	}
	if tx.MinSelfDelegation != nil && tx.MinSelfDelegation.IsPositive() {
		return ErrInvalidParam
	}
	if tx.CommissionRate != nil {
		if tx.CommissionRate.GT(sdk.OneDec()) || tx.CommissionRate.IsNegative() {
			return ErrInvalidParam
		}
	}
	if !tx.ValidatorAddress.Equals(tx.FromAddress) {
		return ErrInvalidParam
	}
	return nil
}

func (tx *MsgEditValidator) Bytes() []byte {
	bytes, err := StakingCodec.MarshalBinaryLengthPrefixed(tx)
	if err != nil {
		panic(err)
	}

	return bytes
}

func (tx *MsgEditValidator) Route() string { return RouteKey }

func (tx *MsgEditValidator) MsgType() string { return "edit-validator" }

func (tx *MsgEditValidator) GetFromAddress() sdk.AccAddress { return tx.FromAddress }
