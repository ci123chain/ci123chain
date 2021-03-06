package keeper

import (
	"bytes"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/couchdb"
	"github.com/ci123chain/ci123chain/pkg/staking/types"
	gogotypes "github.com/gogo/protobuf/types"
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

	//fmt.Printf("validatorUpdates is %v\n", validatorUpdates)

	k.UnbondAllMatureValidatorQueue(ctx)

	matureUnbonds := k.DequeueAllMatureUBDQueue(ctx, ctx.BlockHeader().Time)
	for _, dvPair := range matureUnbonds {
		balances, err := k.CompleteUnbondingWithAmount(ctx, dvPair.DelegatorAddress, dvPair.ValidatorAddress)
		if err != nil {
			continue
		}
		em := ctx.EventManager()
		em.EmitEvent(
			sdk.NewEvent(
				types.EventTypeCompleteUnbonding,
				sdk.NewAttribute([]byte(sdk.AttributeKeyAmount), []byte(balances.String())),
				sdk.NewAttribute([]byte(types.AttributeKeyDstValidator), []byte(dvPair.ValidatorAddress.String())),
				sdk.NewAttribute([]byte(types.AttributeKeyDelegator), []byte(dvPair.DelegatorAddress.String())),
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

		em := ctx.EventManager()
		em.EmitEvent(
			sdk.NewEvent(
				types.EventTypeCompleteRedelegation,
				sdk.NewAttribute([]byte(sdk.AttributeKeyAmount), []byte(balances.String())),
				sdk.NewAttribute([]byte(types.AttributeKeyDelegator), []byte(dvvTriplet.DelegatorAddress.String())),
				sdk.NewAttribute([]byte(types.AttributeKeySrcValidator), []byte(dvvTriplet.ValidatorSrcAddress.String())),
				sdk.NewAttribute([]byte(types.AttributeKeyDstValidator), []byte(dvvTriplet.ValidatorDstAddress.String())),
			),
			)
	}
	return validatorUpdates
}

// Apply and return accumulated updates to the bonded validator set. Also,
// * Updates the active valset as keyed by LastValidatorPowerKey.
// * Updates the total power as keyed by LastTotalPowerKey.
// * Updates validator status' according to updated powers.
// * Updates the fee pool bonded vs not-bonded tokens.
// * Updates relevant indices.
// It gets called once after genesis, another time maybe after genesis transactions,
// then once at every EndBlock.
//
// CONTRACT: Only validators with non-zero power or zero-power that were bonded
// at the previous block height or were removed from the validator set entirely
// are returned to Tendermint.

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
		//fmt.Printf("validator.power = %d\n", validator.PotentialConsensusPower())

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
		copy(valAddrBytes[:], valAddr.Bytes()[:])
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
		err := k.notBondedTokensToBonded(ctx, amtFromNotBondedToBonded.Sub(amtFromBondedToNotBonded))
		if err != nil {
			panic(err)
		}
	case amtFromNotBondedToBonded.LT(amtFromBondedToNotBonded):
		err := k.bondedTokensToNotBonded(ctx, amtFromBondedToNotBonded.Sub(amtFromNotBondedToBonded))
		if err != nil {
			panic(err)
		}
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

	validator.BondedHeight = ctx.BlockHeight()

	// save the now bonded validator record to the two referenced stores
	err := k.SetValidator(ctx, validator)
	if err != nil {
		panic(err)
	}
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
	err := k.SetValidator(ctx, validator)
	if err != nil {
		panic(err)
	}
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
		realKey := iterator.Key()
		_, ok := iterator.(*couchdb.CouchIterator)
		if ok {
			realKey = sdk.GetRealKey(iterator.Key())
		}
		var valAddr [sdk.AddrLen]byte
		copy(valAddr[:], realKey[1:])
		powerBytes := iterator.Value()
		last[valAddr] = make([]byte, len(powerBytes))
		copy(last[valAddr], powerBytes)
	}
	return last
}

// perform all the store operations for when a validator status becomes unbonded
func (k StakingKeeper) completeUnbondingValidator(ctx sdk.Context, validator types.Validator) types.Validator {
	validator = validator.UpdateStatus(sdk.Unbonded)
	err := k.SetValidator(ctx, validator)
	if err != nil {
		panic(err)
	}
	return validator
}

// switches a validator from unbonding state to unbonded state
func (k StakingKeeper) unbondingToUnbonded(ctx sdk.Context, validator types.Validator) types.Validator {
	if !validator.IsUnbonding() {
		panic(fmt.Sprintf("bad state transition unbondingToBonded, validator: %v\n", validator))
	}
	return k.completeUnbondingValidator(ctx, validator)
}