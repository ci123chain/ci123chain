package keeper

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/staking/exported"
	types2 "github.com/ci123chain/ci123chain/pkg/staking/types"
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


// iterate through the active validator set and perform the provided function
func (k StakingKeeper) IterateLastValidators(ctx sdk.Context, fn func(index int64, validator types2.Validator) (stop bool)) {
	iterator := k.LastValidatorsIterator(ctx)
	defer iterator.Close()

	i := int64(0)

	for ; iterator.Valid(); iterator.Next() {
		address :=  sdk.ToAccAddress(iterator.Key())

		validator, found := k.GetValidator(ctx, address)
		if !found {
			panic(fmt.Sprintf("validator record not found for address: %v\n", address))
		}

		stop := fn(i, validator) // XXX is this safe will the validator unexposed fields be able to get written to?
		if stop {
			break
		}
		i++
	}
}