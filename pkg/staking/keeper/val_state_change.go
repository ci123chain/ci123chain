package keeper

import (
	"bytes"
	"fmt"
	gogotypes "github.com/gogo/protobuf/types"
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/staking/types"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	"sort"
)

func (k StakingKeeper) jailValidator(ctx sdk.Context, validator types.Validator) {
	if validator.Jailed {
		panic(fmt.Sprintf("cannot jail already jailed validator, validator: %v\n", validator))
	}

	validator.Jailed = true
	_ = k.SetValidator(ctx, validator)
	k.DeleteValidatorByPowerIndex(ctx, validator)
}

func (k StakingKeeper) BlockValidatorUpdates(ctx sdk.Context) []abcitypes.ValidatorUpdate {

	validatorUpdates := k.ApplyAndReturnValidatorSetUpdates(ctx)

	k.UnbondAllMatureValidatorQueue(ctx)

	matureUnbonds := k.DequeueAllMatureUBDQueue(ctx, ctx.BlockHeader().Time)
	for _, dvPair := range matureUnbonds {
		balances, err := k.CompleteUnbondingWithAmount(ctx, dvPair.DelegatorAddress, dvPair.ValidatorAddress)
		if err != nil {
			continue
		}
		em := sdk.NewEventManager()
		em.EmitEvent(
			sdk.NewEvent(
				types.EventTypeCompleteUnbonding,
				sdk.NewAttribute(sdk.AttributeKeyAmount, balances.String()),
				sdk.NewAttribute(types.AttributeKeyDstValidator, dvPair.ValidatorAddress.String()),
				sdk.NewAttribute(types.AttributeKeyDelegator, dvPair.DelegatorAddress.String()),
				),
			)
	}

	matureRedelegations := k.DequeueAllMatureRedelegationQueue(ctx, ctx.BlockHeader().Time)
	for _, dvvTriplet := range matureRedelegations {
		balances, err := k.CompleteRedelegationWithAmount(
			ctx,
			dvvTriplet.DelegatorAddress,
			dvvTriplet.ValidatorSrcAddress,
			dvvTriplet.ValidatorDstAddress,
		)
		if err != nil {
			continue
		}

		em := sdk.NewEventManager()
		em.EmitEvent(
			sdk.NewEvent(
				types.EventTypeCompleteRedelegation,
				sdk.NewAttribute(sdk.AttributeKeyAmount, balances.String()),
				sdk.NewAttribute(types.AttributeKeyDelegator, dvvTriplet.DelegatorAddress.String()),
				sdk.NewAttribute(types.AttributeKeySrcValidator, dvvTriplet.ValidatorSrcAddress.String()),
				sdk.NewAttribute(types.AttributeKeyDstValidator, dvvTriplet.ValidatorDstAddress.String()),
			),
			)
	}
	return validatorUpdates
}

func (k StakingKeeper) ApplyAndReturnValidatorSetUpdates(ctx sdk.Context) (updates []abcitypes.ValidatorUpdate) {
	maxValidators := k.GetParams(ctx).MaxValidators
	totalPower := sdk.ZeroInt()
	amtFromBondedToNotBonded, amtFromNotBondedToBonded := sdk.ZeroInt(), sdk.ZeroInt()

	// Retrieve the last validator set.
	// The persistent set is updated later in this function.
	// (see LastValidatorPowerKey).
	last := k.getLastValidatorsByAddr(ctx)

	// Iterate over validators, highest power to lowest.
	iterator := k.ValidatorsPowerStoreIterator(ctx)
	defer iterator.Close()
	for count := 0; iterator.Valid() && count < int(maxValidators); iterator.Next() {

		// everything that is iterated in this loop is becoming or already a
		// part of the bonded validator set

		valAddr := sdk.ToAccAddress(iterator.Value())
		validator := k.mustGetValidator(ctx, valAddr)

		if validator.Jailed {
			panic("should never retrieve a jailed validator from the power store")
		}

		// if we get to a zero-power validator (which we don't bond),
		// there are no more possible bonded validators
		if validator.PotentialConsensusPower() == 0 {
			break
		}

		// apply the appropriate state change if necessary
		switch {
		case validator.IsUnbonded():
			validator = k.unbondedToBonded(ctx, validator)
			amtFromNotBondedToBonded = amtFromNotBondedToBonded.Add(validator.GetTokens())
		case validator.IsUnbonding():
			validator = k.unbondingToBonded(ctx, validator)
			amtFromNotBondedToBonded = amtFromNotBondedToBonded.Add(validator.GetTokens())
		case validator.IsBonded():
			// no state change
		default:
			panic("unexpected validator status")
		}

		// fetch the old power bytes
		var valAddrBytes [sdk.AddrLen]byte
		copy(valAddrBytes[:], valAddr.Bytes())
		oldPowerBytes, found := last[valAddrBytes]

		newPower := validator.ConsensusPower()
		newPowerBytes := k.cdc.MustMarshalBinaryLengthPrefixed(&gogotypes.Int64Value{Value: newPower})

		// update the validator set if power has changed
		if !found || !bytes.Equal(oldPowerBytes, newPowerBytes) {
			updates = append(updates, validator.ABCIValidatorUpdate())
			k.SetLastValidatorPower(ctx, valAddr, newPower)
		}

		delete(last, valAddrBytes)

		count++
		totalPower = totalPower.Add(sdk.NewInt(newPower))
	}

	noLongerBonded := sortNoLongerBonded(last)
	for _, valAddrBytes := range noLongerBonded {

		validator := k.mustGetValidator(ctx, sdk.ToAccAddress(valAddrBytes))
		validator = k.bondedToUnbonding(ctx, validator)
		amtFromBondedToNotBonded = amtFromBondedToNotBonded.Add(validator.GetTokens())
		k.DeleteLastValidatorPower(ctx, validator.GetOperator())
		updates = append(updates, validator.ABCIValidatorUpdateZero())
	}

	// Update the pools based on the recent updates in the validator set:
	// - The tokens from the non-bonded candidates that enter the new validator set need to be transferred
	// to the Bonded pool.
	// - The tokens from the bonded validators that are being kicked out from the validator set
	// need to be transferred to the NotBonded pool.
	switch {
	// Compare and subtract the respective amounts to only perform one transfer.
	// This is done in order to avoid doing multiple updates inside each iterator/loop.
	case amtFromNotBondedToBonded.GT(amtFromBondedToNotBonded):
		k.notBondedTokensToBonded(ctx, amtFromNotBondedToBonded.Sub(amtFromBondedToNotBonded))
	case amtFromNotBondedToBonded.LT(amtFromBondedToNotBonded):
		k.bondedTokensToNotBonded(ctx, amtFromBondedToNotBonded.Sub(amtFromNotBondedToBonded))
	default:
		// equal amounts of tokens; no update required
	}

	// set total power on lookup index if there are any updates
	if len(updates) > 0 {
		k.SetLastTotalPower(ctx, totalPower)
	}

	return updates
}

// perform all the store operations for when a validator status becomes bonded
func (k StakingKeeper) bondValidator(ctx sdk.Context, validator types.Validator) types.Validator {
	// delete the validator by power index, as the key will change
	k.DeleteValidatorByPowerIndex(ctx, validator)

	validator = validator.UpdateStatus(sdk.Bonded)

	// save the now bonded validator record to the two referenced stores
	k.SetValidator(ctx, validator)
	k.SetValidatorByPowerIndex(ctx, validator)

	// delete from queue if present
	k.DeleteValidatorQueue(ctx, validator)

	// trigger hook
	k.AfterValidatorBonded(ctx, validator.GetConsAddress(), validator.OperatorAddress)

	return validator
}

// perform all the store operations for when a validator begins unbonding
func (k StakingKeeper) beginUnbondingValidator(ctx sdk.Context, validator types.Validator) types.Validator {
	params := k.GetParams(ctx)

	// delete the validator by power index, as the key will change
	k.DeleteValidatorByPowerIndex(ctx, validator)

	// sanity check
	if validator.Status != sdk.Bonded {
		panic(fmt.Sprintf("should not already be unbonded or unbonding, validator: %v\n", validator))
	}

	validator = validator.UpdateStatus(sdk.Unbonding)

	// set the unbonding completion time and completion height appropriately
	validator.UnbondingTime = ctx.BlockHeader().Time.Add(params.UnbondingTime)
	validator.UnbondingHeight = ctx.BlockHeader().Height

	// save the now unbonded validator record and power index
	k.SetValidator(ctx, validator)
	k.SetValidatorByPowerIndex(ctx, validator)

	// Adds to unbonding validator queue
	k.InsertValidatorQueue(ctx, validator)

	// trigger hook
	k.AfterValidatorBeginUnbonding(ctx, validator.GetConsAddress(), validator.OperatorAddress)

	return validator
}


func (k StakingKeeper) bondedToUnbonding(ctx sdk.Context, validator types.Validator) types.Validator {
	if !validator.IsBonded() {
		panic(fmt.Sprintf("bad state transition bondedToUnbonding, validator: %v\n", validator))
	}
	return k.beginUnbondingValidator(ctx, validator)
}

func (k StakingKeeper) unbondingToBonded(ctx sdk.Context, validator types.Validator) types.Validator {
	if !validator.IsUnbonding() {
		panic(fmt.Sprintf("bad state transition unbondingToBonded, validator: %v\n", validator))
	}
	return k.bondValidator(ctx, validator)
}

func (k StakingKeeper) unbondedToBonded(ctx sdk.Context, validator types.Validator) types.Validator {
	if !validator.IsUnbonded() {
		panic(fmt.Sprintf("bad state transition unbondedToBonded, validator: %v\n", validator))
	}
	return k.bondValidator(ctx, validator)
}

// map of operator addresses to serialized power
type validatorsByAddr map[[sdk.AddrLen]byte][]byte

// given a map of remaining validators to previous bonded power
// returns the list of validators to be unbonded, sorted by operator address
func sortNoLongerBonded(last validatorsByAddr) [][]byte {
	// sort the map keys for determinism
	noLongerBonded := make([][]byte, len(last))
	index := 0
	for valAddrBytes := range last {
		valAddr := make([]byte, sdk.AddrLen)
		copy(valAddr, valAddrBytes[:])
		noLongerBonded[index] = valAddr
		index++
	}
	// sorted by address - order doesn't matter
	sort.SliceStable(noLongerBonded, func(i, j int) bool {
		// -1 means strictly less than
		return bytes.Compare(noLongerBonded[i], noLongerBonded[j]) == -1
	})
	return noLongerBonded
}

// get the last validator set
func (k StakingKeeper) getLastValidatorsByAddr(ctx sdk.Context) validatorsByAddr {
	last := make(validatorsByAddr)
	iterator := k.LastValidatorsIterator(ctx)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var valAddr [sdk.AddrLen]byte
		// extract the validator address from the key (prefix is 1-byte)
		copy(valAddr[:], iterator.Key()[1:])
		powerBytes := iterator.Value()
		last[valAddr] = make([]byte, len(powerBytes))
		copy(last[valAddr], powerBytes)
	}
	return last
}

// perform all the store operations for when a validator status becomes unbonded
func (k StakingKeeper) completeUnbondingValidator(ctx sdk.Context, validator types.Validator) types.Validator {
	validator = validator.UpdateStatus(sdk.Unbonded)
	k.SetValidator(ctx, validator)
	return validator
}

// switches a validator from unbonding state to unbonded state
func (k StakingKeeper) unbondingToUnbonded(ctx sdk.Context, validator types.Validator) types.Validator {
	if !validator.IsUnbonding() {
		panic(fmt.Sprintf("bad state transition unbondingToBonded, validator: %v\n", validator))
	}
	return k.completeUnbondingValidator(ctx, validator)
}