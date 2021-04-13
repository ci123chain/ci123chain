package keeper

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/ci123chain/ci123chain/pkg/staking/types"
)

// Return all validators that a delegator is bonded to. If maxRetrieve is supplied, the respective amount will be returned.
func (k StakingKeeper) GetDelegatorValidators(
	ctx sdk.Context, delegatorAddr sdk.AccAddress, maxRetrieve uint32,
) ([]types.Validator, error) {
	var delegation types.Delegation

	validators := make([]types.Validator, maxRetrieve)

	prefix := types.GetDelegationsKey(delegatorAddr)
	store := ctx.KVStore(k.storeKey)
	iterator := store.RemoteIterator(prefix, sdk.PrefixEndBytes(prefix))
	if !iterator.Valid(){
		iterator.Close()
		store := ctx.KVStore(k.storeKey)
		iterator = sdk.KVStoreReversePrefixIterator(store, prefix)
	}

	defer iterator.Close()

	i := 0
	for ; iterator.Valid() && i < int(maxRetrieve); iterator.Next() {
		types.StakingCodec.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &delegation)

		validator, found := k.GetValidator(ctx, delegation.ValidatorAddress)
		if !found {
			return nil, sdkerrors.Wrap(sdkerrors.ErrResponse, "no validator found")
		}

		validators[i] = validator
		i++
	}

	return validators[:i], nil // trim
}

// return a validator that a delegator is bonded to
func (k StakingKeeper) GetDelegatorValidator(
	ctx sdk.Context, delegatorAddr sdk.AccAddress, validatorAddr sdk.AccAddress,
) (validator types.Validator, err error) {

	delegation, found := k.GetDelegation(ctx, delegatorAddr, validatorAddr)
	if !found {
		return validator, sdkerrors.Wrap(sdkerrors.ErrResponse, "no delegation found")
	}

	validator, found = k.GetValidator(ctx, delegation.ValidatorAddress)
	if !found {
		return validator, sdkerrors.Wrap(sdkerrors.ErrResponse, "no validator found")
	}

	return validator, nil
}


// return all redelegations for a delegator
func (k StakingKeeper) GetAllRedelegations(
	ctx sdk.Context, delegator sdk.AccAddress, srcValAddress, dstValAddress sdk.AccAddress,
) []types.Redelegation {
	prefix := types.GetREDsKey(delegator)
	store := ctx.KVStore(k.storeKey)
	iterator := store.RemoteIterator(prefix, sdk.PrefixEndBytes(prefix))
	if !iterator.Valid() {
		iterator.Close()
		store := ctx.KVStore(k.storeKey)
		iterator = sdk.KVStoreReversePrefixIterator(store, prefix)
	}

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
	prefix := types.GetDelegationsKey(delegator)
	store := ctx.KVStore(k.storeKey)
	iterator := store.RemoteIterator(prefix, sdk.PrefixEndBytes(prefix))
	if !iterator.Valid() {
		iterator.Close()
		store := ctx.KVStore(k.storeKey)
		iterator = sdk.KVStoreReversePrefixIterator(store, prefix)
	}
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