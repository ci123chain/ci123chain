package keeper

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/staking/types"
)

// Return all validators that a delegator is bonded to. If maxRetrieve is supplied, the respective amount will be returned.
func (k StakingKeeper) GetDelegatorValidators(
	ctx sdk.Context, delegatorAddr sdk.AccAddress, maxRetrieve uint32,
) ([]types.Validator, sdk.Error) {
	var delegation types.Delegation

	validators := make([]types.Validator, maxRetrieve)

	delegatorPrefixKey := types.GetDelegationsKey(delegatorAddr)
	iterator := k.cdb.Iterator(sdk.NewPrefixedKey([]byte(k.storeKey.Name()), delegatorPrefixKey), sdk.NewPrefixedKey([]byte(k.storeKey.Name()), sdk.PrefixEndBytes(delegatorPrefixKey)))
	defer iterator.Close()

	i := 0
	for ; iterator.Valid() && i < int(maxRetrieve); iterator.Next() {
		types.StakingCodec.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &delegation)

		validator, found := k.GetValidator(ctx, delegation.ValidatorAddress)
		if !found {
			return nil, sdk.ErrNoValidatorFound("no validator found")
		}

		validators[i] = validator
		i++
	}

	return validators[:i], nil // trim
}

// return a validator that a delegator is bonded to
func (k StakingKeeper) GetDelegatorValidator(
	ctx sdk.Context, delegatorAddr sdk.AccAddress, validatorAddr sdk.AccAddress,
) (validator types.Validator, err sdk.Error) {

	delegation, found := k.GetDelegation(ctx, delegatorAddr, validatorAddr)
	if !found {
		return validator, sdk.ErrNoDelegation("no delegation")
	}

	validator, found = k.GetValidator(ctx, delegation.ValidatorAddress)
	if !found {
		return validator, sdk.ErrNoValidatorFound("no validator found")
	}

	return validator, nil
}


// return all redelegations for a delegator
func (k StakingKeeper) GetAllRedelegations(
	ctx sdk.Context, delegator sdk.AccAddress, srcValAddress, dstValAddress sdk.AccAddress,
) []types.Redelegation {
	delegatorPrefixKey := types.GetREDsKey(delegator)
	iterator := k.cdb.Iterator(sdk.NewPrefixedKey([]byte(k.storeKey.Name()), delegatorPrefixKey), sdk.NewPrefixedKey([]byte(k.storeKey.Name()), sdk.PrefixEndBytes(delegatorPrefixKey)))

	defer iterator.Close()

	srcValFilter := !(srcValAddress.Empty())
	dstValFilter := !(dstValAddress.Empty())

	redelegations := []types.Redelegation{}
	var redelegation types.Redelegation

	for ; iterator.Valid(); iterator.Next() {
		types.StakingCodec.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &redelegation)
		if srcValFilter && !(srcValAddress.Equals(redelegation.ValidatorSrcAddress)) {
			continue
		}
		if dstValFilter && !(dstValAddress.Equals(redelegation.ValidatorDstAddress)) {
			continue
		}

		redelegations = append(redelegations, redelegation)
	}

	return redelegations
}

// return all delegations for a delegator
func (k StakingKeeper) GetAllDelegatorDelegations(ctx sdk.Context, delegator sdk.AccAddress) []types.Delegation {
	delegations := make([]types.Delegation, 0)
	delegatorPrefixKey := types.GetDelegationsKey(delegator)
	iterator := k.cdb.Iterator(sdk.NewPrefixedKey([]byte(k.storeKey.Name()), delegatorPrefixKey), sdk.NewPrefixedKey([]byte(k.storeKey.Name()), sdk.PrefixEndBytes(delegatorPrefixKey)))

	defer iterator.Close()

	i := 0
	for ; iterator.Valid(); iterator.Next() {
		var delegation types.Delegation
		types.StakingCodec.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &delegation)
		delegations = append(delegations, delegation)
		i++
	}

	return delegations
}