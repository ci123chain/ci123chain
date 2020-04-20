package keeper

import (
	"bytes"
	"fmt"
	gogotypes "github.com/gogo/protobuf/types"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/staking/types"
	"time"
)

type cachedValidator struct {
	val   types.Validator
	marshalled  string
}

func newCachedValidator(val types.Validator, marshalled string) cachedValidator {
	return cachedValidator{
		val:        val,
		marshalled: marshalled,
	}
}

func (k StakingKeeper) GetValidator(ctx sdk.Context, addr sdk.AccAddress) (validator types.Validator, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.GetValidatorKey(addr))
	if value == nil {
		return validator, false
	}

	err := types.StakingCodec.UnmarshalBinaryLengthPrefixed(value, &validator)
	if err != nil {
		return validator, false
	}
	return validator, true
}

func (k StakingKeeper) mustGetValidator(ctx sdk.Context, addr sdk.AccAddress) types.Validator {
	validator, found := k.GetValidator(ctx, addr)
	if !found {
		panic(fmt.Sprintf("validator record not found for address: %X\n", addr))
	}
	return validator
}

func (k StakingKeeper) GetValidatorByConsAddr(ctx sdk.Context, consAddr sdk.AccAddress) (validator types.Validator, found bool) {
	store := ctx.KVStore(k.storeKey)
	opAddr := store.Get(types.GetValidatorByConsAddrKey(consAddr))
	if opAddr == nil {
		return validator, false
	}
	op := sdk.ToAccAddress(opAddr)
	return k.GetValidator(ctx, op)
}

func (k StakingKeeper) SetValidator(ctx sdk.Context, validator types.Validator) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := types.StakingCodec.MarshalBinaryLengthPrefixed(validator)
	if err != nil {
		return err
	}
	store.Set(types.GetValidatorKey(validator.OperatorAddress), bz)
	return nil
}

func (k StakingKeeper) SetValidatorByConsAddr(ctx sdk.Context, validator types.Validator) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetValidatorByConsAddrKey(validator.GetConsAddress()), validator.OperatorAddress.Bytes())
}

func (k StakingKeeper) SetValidatorByPowerIndex(ctx sdk.Context, validator types.Validator) {
	if validator.Jailed {
		return
	}
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetValidatorsByPowerIndexKey(validator), validator.OperatorAddress.Bytes())
}

func (k StakingKeeper) SetNewValidatorByPowerIndex(ctx sdk.Context, validator types.Validator) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetValidatorsByPowerIndexKey(validator), validator.OperatorAddress.Bytes())
}

func (k StakingKeeper) DeleteValidatorByPowerIndex(ctx sdk.Context, validator types.Validator) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetValidatorsByPowerIndexKey(validator))
}

func (k StakingKeeper) AddValidatorTokensAndShares(ctx sdk.Context, validator types.Validator, tokensToAdd sdk.Int) (valOut types.Validator,
	addedShares sdk.Dec) {
	k.DeleteValidatorByPowerIndex(ctx, validator)
	validator, addedShares = validator.AddTokensFromDel(tokensToAdd)
	_ = k.SetValidator(ctx, validator)
	k.SetValidatorByPowerIndex(ctx, validator)
	return validator, addedShares
}

func (k StakingKeeper) RemoveValidator(ctx sdk.Context, address sdk.AccAddress) {
	// first retrieve the old validator record
	validator, found := k.GetValidator(ctx, address)
	if !found {
		return
	}

	if !validator.IsUnbonded() {
		panic("cannot call RemoveValidator on bonded or unbonding validators")
	}
	if validator.Tokens.IsPositive() {
		panic("attempting to remove a validator which still contains tokens")
	}
	if validator.Tokens.IsPositive() {
		panic("validator being removed should never have positive tokens")
	}

	valConsAddr := validator.GetConsAddress()

	// delete the old validator record
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetValidatorKey(address))
	store.Delete(types.GetValidatorByConsAddrKey(valConsAddr))
	store.Delete(types.GetValidatorsByPowerIndexKey(validator))

	// call hooks
	k.AfterValidatorRemoved(ctx, valConsAddr, validator.OperatorAddress)
}

// Update the tokens of an existing validator, update the validators power index key
func (k StakingKeeper) RemoveValidatorTokensAndShares(ctx sdk.Context, validator types.Validator,
	sharesToRemove sdk.Dec) (valOut types.Validator, removedTokens sdk.Int) {

	k.DeleteValidatorByPowerIndex(ctx, validator)
	validator, removedTokens = validator.RemoveDelShares(sharesToRemove)
	_ = k.SetValidator(ctx, validator)
	k.SetValidatorByPowerIndex(ctx, validator)
	return validator, removedTokens
}

// Delete the last validator power.
func (k StakingKeeper) DeleteLastValidatorPower(ctx sdk.Context, operator sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetLastValidatorPowerKey(operator))
}

// gets a specific validator queue timeslice. A timeslice is a slice of ValAddresses corresponding to unbonding validators
// that expire at a certain time.
func (k StakingKeeper) GetValidatorQueueTimeSlice(ctx sdk.Context, timestamp time.Time) []sdk.AccAddress {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetValidatorQueueTimeKey(timestamp))
	if bz == nil {
		return nil
	}

	var va []sdk.AccAddress
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &va)
	return va
}


// Insert an validator address to the appropriate timeslice in the validator queue
func (k StakingKeeper) InsertValidatorQueue(ctx sdk.Context, val types.Validator) {
	timeSlice := k.GetValidatorQueueTimeSlice(ctx, val.UnbondingTime)
	timeSlice = append(timeSlice, val.OperatorAddress)
	k.SetValidatorQueueTimeSlice(ctx, val.UnbondingTime, timeSlice)
}

// Sets a specific validator queue timeslice.
func (k StakingKeeper) SetValidatorQueueTimeSlice(ctx sdk.Context, timestamp time.Time, keys []sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(&keys)
	store.Set(types.GetValidatorQueueTimeKey(timestamp), bz)
}


// Delete a validator address from the validator queue
func (k StakingKeeper) DeleteValidatorQueue(ctx sdk.Context, val types.Validator) {
	timeSlice := k.GetValidatorQueueTimeSlice(ctx, val.UnbondingTime)
	var newTimeSlice []sdk.AccAddress
	for _, addr := range timeSlice {
		if !bytes.Equal(addr.Bytes(), val.OperatorAddress.Bytes()) {
			newTimeSlice = append(newTimeSlice, addr)
		}
	}

	if len(newTimeSlice) == 0 {
		k.DeleteValidatorQueueTimeSlice(ctx, val.UnbondingTime)
	} else {
		k.SetValidatorQueueTimeSlice(ctx, val.UnbondingTime, newTimeSlice)
	}
}

// Deletes a specific validator queue timeslice.
func (k StakingKeeper) DeleteValidatorQueueTimeSlice(ctx sdk.Context, timestamp time.Time) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetValidatorQueueTimeKey(timestamp))
}

// Set the last validator power.
func (k StakingKeeper) SetLastValidatorPower(ctx sdk.Context, operator sdk.AccAddress, power int64) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(&gogotypes.Int64Value{Value: power})
	store.Set(types.GetLastValidatorPowerKey(operator), bz)
}

// returns an iterator for the current validator power store
func (k StakingKeeper) ValidatorsPowerStoreIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStoreReversePrefixIterator(store, types.ValidatorsByPowerIndexKey)
}

// returns an iterator for the consensus validators in the last block
func (k StakingKeeper) LastValidatorsIterator(ctx sdk.Context) (iterator sdk.Iterator) {
	store := ctx.KVStore(k.storeKey)
	iterator = sdk.KVStorePrefixIterator(store, types.LastValidatorPowerKey)
	return iterator
}

// Returns all the validator queue timeslices from time 0 until endTime
func (k StakingKeeper) ValidatorQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(types.ValidatorQueueKey, sdk.InclusiveEndBytes(types.GetValidatorQueueTimeKey(endTime)))
}

// Unbonds all the unbonding validators that have finished their unbonding period
func (k StakingKeeper) UnbondAllMatureValidatorQueue(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	validatorTimesliceIterator := k.ValidatorQueueIterator(ctx, ctx.BlockHeader().Time)
	defer validatorTimesliceIterator.Close()

	for ; validatorTimesliceIterator.Valid(); validatorTimesliceIterator.Next() {
		var timeslice []sdk.AccAddress
		k.cdc.MustUnmarshalBinaryLengthPrefixed(validatorTimesliceIterator.Value(), &timeslice)

		for _, valAddr := range timeslice {
			val, found := k.GetValidator(ctx, valAddr)
			if !found {
				panic("validator in the unbonding queue was not found")
			}

			if !val.IsUnbonding() {
				panic("unexpected validator in unbonding queue; status was not unbonding")
			}

			val = k.unbondingToUnbonded(ctx, val)
			if val.GetDelegatorShares().IsZero() {
				k.RemoveValidator(ctx, val.OperatorAddress)
			}
		}

		store.Delete(validatorTimesliceIterator.Key())
	}
}

func (k StakingKeeper) GetLastValidators(ctx sdk.Context) (validators []types.Validator) {
	store := ctx.KVStore(k.storeKey)

	// add the actual validator power sorted store
	maxValidators := k.MaxValidators(ctx)
	validators = make([]types.Validator, maxValidators)

	iterator := sdk.KVStorePrefixIterator(store, types.LastValidatorPowerKey)
	defer iterator.Close()

	i := 0
	for ; iterator.Valid(); iterator.Next() {

		// sanity check
		if i >= int(maxValidators) {
			panic("more validators than maxValidators found")
		}
		address := types.AddressFromLastValidatorPowerKey(iterator.Key())
		validator := k.mustGetValidator(ctx, sdk.ToAccAddress(address))

		validators[i] = validator
		i++
	}
	return validators[:i] // trim
}

// get the set of all validators with no limits, used during genesis dump
func (k StakingKeeper) GetAllValidators(ctx sdk.Context) (validators []types.Validator) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ValidatorsKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var validator types.Validator
		err := types.StakingCodec.UnmarshalBinaryLengthPrefixed(iterator.Value(), &validator)
		if err != nil {
			panic(err)
		}
		validators = append(validators, validator)
	}

	return validators
}