package keeper

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/staking/exported"
)


// Validator gets the Validator interface for a particular address
func (k StakingKeeper) Validator(ctx sdk.Context, address sdk.AccAddress) exported.ValidatorI {
	val, found := k.GetValidator(ctx, address)
	if !found {
		return nil
	}
	return val
}

// Delegation get the delegation interface for a particular set of delegator and validator addresses
func (k StakingKeeper) Delegation(ctx sdk.Context, addrDel sdk.AccAddress, addrVal sdk.AccAddress) exported.DelegationI {
	bond, ok := k.GetDelegation(ctx, addrDel, addrVal)
	if !ok {
		return nil
	}

	return bond
}