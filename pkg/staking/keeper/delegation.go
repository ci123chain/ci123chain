package keeper

import (
	"bytes"
	"errors"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/staking/types"
	"time"
)

func (k StakingKeeper) GetDelegation(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.AccAddress) (delegation types.Delegation,
	found bool) {
	//
	store := ctx.KVStore(k.storeKey)
	key := types.GetDelegationKey(delAddr, valAddr)
	value := store.Get(key)
	if value == nil {
		return delegation, false
	}
	types.StakingCodec.MustUnmarshalBinaryLengthPrefixed(value, &delegation)
	return delegation, true
}

func (k StakingKeeper) Delegate(ctx sdk.Context, delAddr sdk.AccAddress, bondAmt sdk.Int, tokenSrc sdk.BondStatus,
	validator types.Validator, subtractAccount bool) (newShares sdk.Dec, err error) {
	if validator.InvalidExRate() {
		return sdk.ZeroDec(), types.ErrDelegatorShareExRateInvalid
	}

	delegation, found := k.GetDelegation(ctx, delAddr, validator.OperatorAddress)
	if !found {
		delegation = types.NewDelegation(delAddr, validator.OperatorAddress, sdk.ZeroDec())
	}

	if found {
		k.BeforeDelegationSharesModified(ctx, delAddr, validator.OperatorAddress)
	}else {
		k.BeforeDelegationCreated(ctx, delAddr, validator.OperatorAddress)
	}

	if subtractAccount {
		if tokenSrc == sdk.Bonded {
			panic("delegation token source cannot be bonded")
		}

		var sendName string
		switch {
		case validator.IsBonded():
			sendName = types.BondedPoolName
		case validator.IsUnbonding(), validator.IsUnbonded():
			sendName = types.NotBondedPoolName
		default:
			panic("invalid validator status")
			//return newShares, types.ErrInvalidValidatorStatus
		}

		//coins := sdk.NewCoins(sdk.NewCoin(bondAmt))
		coin := sdk.NewChainCoin(bondAmt)
		err := k.SupplyKeeper.DelegateCoinsFromAccountToModule(ctx, delegation.DelegatorAddress, sendName, coin)
		if err != nil {
			return sdk.Dec{}, err
		}
	} else {

		switch {
		case tokenSrc == sdk.Bonded && validator.IsBonded():
			//
		case (tokenSrc == sdk.Unbonded || tokenSrc == sdk.Unbonding) && !validator.IsBonded():
			//
		case (tokenSrc == sdk.Unbonded || tokenSrc == sdk.Unbonding) && validator.IsBonded():
			err := k.notBondedTokensToBonded(ctx, bondAmt)
			if err != nil {
				return newShares, types.ErrBondedTokendFailed
			}
		case tokenSrc == sdk.Bonded && !validator.IsBonded():
			// transfer pools
			err := k.bondedTokensToNotBonded(ctx, bondAmt)
			if err != nil {
				return newShares, types.ErrBondedTokensToNoBondedFailed
			}
		default:
			//return newShares, types.ErrUnknowTokenSource
			panic("unknown token source bond status")
		}
	}

	validator, newShares = k.AddValidatorTokensAndShares(ctx, validator, bondAmt)
	//update delegation
	delegation.Shares = delegation.Shares.Add(newShares)
	k.SetDelegation(ctx, delegation)
	k.AfterDelegationModified(ctx, delegation.DelegatorAddress, delegation.ValidatorAddress)
	return newShares, nil
}

func (k StakingKeeper) SetDelegation(ctx sdk.Context, delegation types.Delegation) {
	store := ctx.KVStore(k.storeKey)
	b := types.StakingCodec.MustMarshalBinaryLengthPrefixed(delegation)
	store.Set(types.GetDelegationKey(delegation.DelegatorAddress, delegation.ValidatorAddress), b)
}

// remove a delegation
func (k StakingKeeper) RemoveDelegation(ctx sdk.Context, delegation types.Delegation) {
	k.BeforeDelegationRemoved(ctx, delegation.DelegatorAddress, delegation.ValidatorAddress)
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetDelegationKey(delegation.DelegatorAddress, delegation.ValidatorAddress))
}


func (k StakingKeeper) ValidateUnbondAmount(ctx sdk.Context, delAddr sdk.AccAddress,
	valAddr sdk.AccAddress, amt sdk.Int) (shares sdk.Dec, err error) {
	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		return shares, types.ErrNoExpectedValidator
	}
	shares, err = validator.SharesFromTokens(amt)
	if err != nil {
		return shares, err
	}

	del, found := k.GetDelegation(ctx, delAddr, valAddr)
	if !found {
		return shares, types.ErrNoDelegation

	}
	sharesTruncated, err := validator.SharesFromTokensTruncated(amt)
	if err != nil {
		return shares, err
	}
	delShares := del.GetShares()

	if sharesTruncated.GT(delShares) {
		return shares, types.ErrBadSharesAmount
	}

	if shares.GT(delShares) {
		shares = delShares
	}

	return shares, nil
}

func (k StakingKeeper) Redelegate(ctx sdk.Context, delAddr sdk.AccAddress, valSrcAddr, valDstAddr sdk.AccAddress,
	sharesAmount sdk.Dec) (completionTime time.Time, err error) {
	if bytes.Equal(valSrcAddr.Bytes(), valDstAddr.Bytes()) {
		return time.Time{}, types.ErrSelfRedelegation
	}

	dstValidator, found := k.GetValidator(ctx, valDstAddr)
	if !found {
		return time.Time{}, types.ErrBadRedelegationDst
	}
	srcValidator, found := k.GetValidator(ctx, valSrcAddr)
	if !found {
		return time.Time{}, types.ErrBadRedelegationDst
	}

	if k.HasReceivingRedelegation(ctx, delAddr, valSrcAddr) {
		return time.Time{}, types.ErrTransitiveRedelegation
	}

	if k.HasMaxRedelegationEntries(ctx, delAddr, valSrcAddr, valDstAddr) {
		return time.Time{}, types.ErrMaxRedelegationEntries
	}

	returnAmount, err := k.unbond(ctx, delAddr, valSrcAddr, sharesAmount)
	if err != nil {
		return time.Time{}, err
	}

	if returnAmount.IsZero() {
		return time.Time{}, types.ErrTinyRedelegationAmount
	}

	sharesCreated, err := k.Delegate(ctx, delAddr, returnAmount, srcValidator.GetStatus(), dstValidator, false)
	if err != nil {
		return time.Time{}, err
	}

	completionTime, height, completeNow := k.getBeginInfo(ctx, valSrcAddr)


	if completeNow {
		return completionTime, nil
	}

	red := k.SetRedelegationEntry(ctx, delAddr, valSrcAddr, valDstAddr,
		height, completionTime, returnAmount, sharesAmount, sharesCreated)

	k.InsertRedelegationQueue(ctx, red, completionTime)
	return completionTime, nil
}

func (k StakingKeeper) Undelegate(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.AccAddress, sharesAmount sdk.Dec) (time.Time, error) {
	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		return time.Time{}, types.ErrNoDelegatorForAddress
	}

	if k.HasMaxUnbondingDelegationEntries(ctx, delAddr, valAddr) {
		return time.Time{}, types.ErrMaxUnbondingDelegationEntries
	}

	returnAmount, err := k.unbond(ctx, delAddr, valAddr, sharesAmount)
	if err != nil {
		return time.Time{}, err
	}

	if validator.IsBonded() {
		err := k.bondedTokensToNotBonded(ctx, returnAmount)
		if err != nil {
			return time.Time{}, err
		}
	}
	completionTime := ctx.BlockHeader().Time.Add(k.UnbondingTime(ctx))
	ubd := k.SetUnbondingDelegationEntry(ctx, delAddr, valAddr, ctx.BlockHeight(), completionTime, returnAmount)
	k.InsertUBDQueue(ctx, ubd, completionTime)

	return completionTime, nil
}

func (k StakingKeeper) HasReceivingRedelegation(ctx sdk.Context, delAddr sdk.AccAddress, valDstAddr sdk.AccAddress) bool {
	prefix := types.GetREDsByDelToValDstIndexKey(delAddr, valDstAddr)
	store := ctx.KVStore(k.storeKey)
	iterator := store.RemoteIterator(prefix, sdk.PrefixEndBytes(prefix))
	if !iterator.Valid() {
		iterator.Close()
		store := ctx.KVStore(k.storeKey)
		iterator = sdk.KVStorePrefixIterator(store, prefix)
	}
	defer iterator.Close()

	return iterator.Valid()
}

// HasMaxRedelegationEntries - redelegation has maximum number of entries
func (k StakingKeeper) HasMaxRedelegationEntries(ctx sdk.Context,
	delegatorAddr sdk.AccAddress, validatorSrcAddr,
	validatorDstAddr sdk.AccAddress) bool {

	red, found := k.GetRedelegation(ctx, delegatorAddr, validatorSrcAddr, validatorDstAddr)
	if !found {
		return false
	}
	return len(red.Entries) >= int(k.MaxEntries(ctx))
}

// return a redelegation
func (k StakingKeeper) GetRedelegation(ctx sdk.Context,
	delAddr sdk.AccAddress, valSrcAddr, valDstAddr sdk.AccAddress) (red types.Redelegation, found bool) {

	store := ctx.KVStore(k.storeKey)
	key := types.GetREDKey(delAddr, valSrcAddr, valDstAddr)
	value := store.Get(key)
	if value == nil {
		return red, false
	}

	types.StakingCodec.MustUnmarshalBinaryLengthPrefixed(value, &red)
	return red, true
}

// HasMaxUnbondingDelegationEntries - check if unbonding delegation has maximum number of entries
func (k StakingKeeper) HasMaxUnbondingDelegationEntries(ctx sdk.Context,
	delegatorAddr sdk.AccAddress, validatorAddr sdk.AccAddress) bool {

	ubd, found := k.GetUnbondingDelegation(ctx, delegatorAddr, validatorAddr)
	if !found {
		return false
	}
	return len(ubd.Entries) >= int(k.MaxEntries(ctx))
}

// return a given amount of all the delegator unbonding-delegations
func (k StakingKeeper) GetUnbondingDelegations(ctx sdk.Context, delegator sdk.AccAddress,
	maxRetrieve uint16) (unbondingDelegations []types.UnbondingDelegation) {

	unbondingDelegations = make([]types.UnbondingDelegation, maxRetrieve)

	prefix := types.GetUBDsKey(delegator)
	store := ctx.KVStore(k.storeKey)
	iterator := store.RemoteIterator(prefix, sdk.PrefixEndBytes(prefix))
	if !iterator.Valid() {
		iterator.Close()
		store := ctx.KVStore(k.storeKey)
		iterator = sdk.KVStorePrefixIterator(store, prefix)
	}
	defer iterator.Close()

	i := 0
	for ; iterator.Valid() && i < int(maxRetrieve); iterator.Next() {
		var unbondingDelegation types.UnbondingDelegation
		types.StakingCodec.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &unbondingDelegation)
		unbondingDelegations[i] = unbondingDelegation
		i++
	}
	return unbondingDelegations[:i] // trim if the array length < maxRetrieve
}

// return a unbonding delegation
func (k StakingKeeper) GetUnbondingDelegation(
	ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.AccAddress,
) (ubd types.UnbondingDelegation, found bool) {

	store := ctx.KVStore(k.storeKey)
	key := types.GetUBDKey(delAddr, valAddr)
	value := store.Get(key)
	if value == nil {
		return ubd, false
	}

	types.StakingCodec.MustUnmarshalBinaryLengthPrefixed(value, &ubd)
	return ubd, true
}

// set the unbonding delegation and associated index
func (k StakingKeeper) SetUnbondingDelegation(ctx sdk.Context, ubd types.UnbondingDelegation) {
	store := ctx.KVStore(k.storeKey)
	bz := types.StakingCodec.MustMarshalBinaryLengthPrefixed(ubd)
	key := types.GetUBDKey(ubd.DelegatorAddress, ubd.ValidatorAddress)
	store.Set(key, bz)
	store.Set(types.GetUBDByValIndexKey(ubd.DelegatorAddress, ubd.ValidatorAddress), []byte{}) // index, store empty bytes
}

// remove the unbonding delegation object and associated index
func (k StakingKeeper) RemoveUnbondingDelegation(ctx sdk.Context, ubd types.UnbondingDelegation) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetUBDKey(ubd.DelegatorAddress, ubd.ValidatorAddress)
	store.Delete(key)
	store.Delete(types.GetUBDByValIndexKey(ubd.DelegatorAddress, ubd.ValidatorAddress))
}

func (k StakingKeeper) SetUnbondingDelegationEntry(ctx sdk.Context, delegatorAddr sdk.AccAddress, validatorAddr sdk.AccAddress,
	creationHeight int64, minTime time.Time, balance sdk.Int) types.UnbondingDelegation {
	ubd, found := k.GetUnbondingDelegation(ctx, delegatorAddr, validatorAddr)
	if found {
		ubd.AddEntry(creationHeight, minTime, balance)
	} else {
		ubd = types.NewUnbondingDelegation(delegatorAddr, validatorAddr, creationHeight, minTime, balance)
	}
	k.SetUnbondingDelegation(ctx, ubd)
	return ubd
}

func (k StakingKeeper) Unbond(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.AccAddress, shares sdk.Dec) (amount sdk.Int, err error) {
	return k.unbond(ctx, delAddr, valAddr, shares)
}

func (k StakingKeeper) unbond(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.AccAddress, shares sdk.Dec) (amount sdk.Int, err error) {
	delegation, found := k.GetDelegation(ctx, delAddr, valAddr)
	if !found {
		return amount, types.ErrNoDelegatorForAddress
	}

	k.BeforeDelegationSharesModified(ctx, delAddr, valAddr)

	if delegation.Shares.LT(shares) {
		return amount, types.ErrNotEnoughDelegationShares
	}

	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		return amount, types.ErrNoValidatorFound
	}

	delegation.Shares = delegation.Shares.Sub(shares)
	isValidatorOperator := delegation.DelegatorAddress.Equals(validator.OperatorAddress)

	if isValidatorOperator && !validator.Jailed && validator.
		TokensFromShares(delegation.Shares).TruncateInt().LT(validator.MinSelfDelegation) {
		k.jailValidator(ctx, validator)
		validator = k.mustGetValidator(ctx, validator.OperatorAddress)
	}
	if delegation.Shares.IsZero() {
		k.RemoveDelegation(ctx, delegation)
	} else {
		k.SetDelegation(ctx, delegation)
		k.AfterDelegationModified(ctx, delegation.DelegatorAddress, delegation.ValidatorAddress)
	}

	validator, amount = k.RemoveValidatorTokensAndShares(ctx, validator, shares)

	if validator.DelegatorShares.IsZero() && validator.IsUnbonded() {
		k.RemoveValidator(ctx, validator.OperatorAddress)
	}
	return amount, nil
}

func(k StakingKeeper) getBeginInfo(ctx sdk.Context, valSrcAddr sdk.AccAddress,
	) (completionTime time.Time, height int64, completeNow bool) {
	validator, found := k.GetValidator(ctx, valSrcAddr)

	switch {
	case !found || validator.IsBonded():
		completionTime = ctx.BlockHeader().Time.Add(k.UnbondingTime(ctx))
		height = ctx.BlockHeight()
		return completionTime, height, false
	case validator.IsUnbonded():
		return completionTime, height, true
	case validator.IsUnbonding():
		return validator.UnbondingTime, validator.UnbondingHeight, false
	default:
		panic(fmt.Sprintf("unknown validator status: %s", validator.Status))
	}
}

func (k StakingKeeper) SetRedelegationEntry(ctx sdk.Context, delegatorAddr sdk.AccAddress,
	validatorSrcAddr, validatorDstAddr sdk.AccAddress, creationHeight int64,
	minTime time.Time, balance sdk.Int, sharesSrc, sharesDst sdk.Dec) types.Redelegation {
	red, found := k.GetRedelegation(ctx, delegatorAddr, validatorSrcAddr, validatorDstAddr)
	if found {
		red.AddEntry(creationHeight, minTime, balance, sharesDst)
	} else {
		red = types.NewRedelegation(delegatorAddr, validatorSrcAddr,
			validatorDstAddr, creationHeight, minTime, balance, sharesDst)
	}
	k.SetRedelegation(ctx, red)
	return red
}

// set a redelegation and associated index
func (k StakingKeeper) SetRedelegation(ctx sdk.Context, red types.Redelegation) {
	store := ctx.KVStore(k.storeKey)
	bz := types.StakingCodec.MustMarshalBinaryLengthPrefixed(red)
	key := types.GetREDKey(red.DelegatorAddress, red.ValidatorSrcAddress, red.ValidatorDstAddress)
	store.Set(key, bz)
	store.Set(types.GetREDByValSrcIndexKey(red.DelegatorAddress, red.ValidatorSrcAddress, red.ValidatorDstAddress), []byte{})
	store.Set(types.GetREDByValDstIndexKey(red.DelegatorAddress, red.ValidatorSrcAddress, red.ValidatorDstAddress), []byte{})
}


// Gets a specific redelegation queue timeslice. A timeslice is a slice of DVVTriplets corresponding to redelegations
// that expire at a certain time.
func (k StakingKeeper) GetRedelegationQueueTimeSlice(ctx sdk.Context, timestamp time.Time) (dvvTriplets []types.DVVTriplet) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetRedelegationTimeKey(timestamp))
	if bz == nil {
		return []types.DVVTriplet{}
	}

	triplets := types.DVVTriplets{}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &triplets)
	return triplets.Triplets
}

// Sets a specific redelegation queue timeslice.
func (k StakingKeeper) SetRedelegationQueueTimeSlice(ctx sdk.Context, timestamp time.Time, keys []types.DVVTriplet) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(&types.DVVTriplets{Triplets: keys})
	store.Set(types.GetRedelegationTimeKey(timestamp), bz)
}

func (k StakingKeeper) InsertRedelegationQueue(ctx sdk.Context, red types.Redelegation,
	completionTime time.Time) {
	timeSlice := k.GetRedelegationQueueTimeSlice(ctx, completionTime)
	dvvTriplet := types.DVVTriplet{
		DelegatorAddress:    red.DelegatorAddress,
		ValidatorSrcAddress: red.ValidatorSrcAddress,
		ValidatorDstAddress: red.ValidatorDstAddress}

	if len(timeSlice) == 0 {
		k.SetRedelegationQueueTimeSlice(ctx, completionTime, []types.DVVTriplet{dvvTriplet})
	} else {
		timeSlice = append(timeSlice, dvvTriplet)
		k.SetRedelegationQueueTimeSlice(ctx, completionTime, timeSlice)
	}
}

// Sets a specific unbonding queue timeslice.
func (k StakingKeeper) SetUBDQueueTimeSlice(ctx sdk.Context, timestamp time.Time, keys []types.DVPair) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(&types.DVPairs{Pairs: keys})
	store.Set(types.GetUnbondingDelegationTimeKey(timestamp), bz)
}

func (k StakingKeeper) InsertUBDQueue(ctx sdk.Context, ubd types.UnbondingDelegation,
	completionTime time.Time) {

	timeSlice := k.GetUBDQueueTimeSlice(ctx, completionTime)
	dvPair := types.DVPair{DelegatorAddress: ubd.DelegatorAddress, ValidatorAddress: ubd.ValidatorAddress}
	if len(timeSlice) == 0 {
		k.SetUBDQueueTimeSlice(ctx, completionTime, []types.DVPair{dvPair})
	} else {
		timeSlice = append(timeSlice, dvPair)
		k.SetUBDQueueTimeSlice(ctx, completionTime, timeSlice)
	}
}

func (k StakingKeeper) GetUBDQueueTimeSlice(ctx sdk.Context, timestamp time.Time) (dvPairs []types.DVPair) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetUnbondingDelegationTimeKey(timestamp))
	if bz == nil {
		return []types.DVPair{}
	}

	pairs := types.DVPairs{}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &pairs)
	return pairs.Pairs
}

// Returns a concatenated list of all the timeslices inclusively previous to
// currTime, and deletes the timeslices from the queue
func (k StakingKeeper) DequeueAllMatureRedelegationQueue(ctx sdk.Context, currTime time.Time) (matureRedelegations []types.DVVTriplet) {
	store := ctx.KVStore(k.storeKey)

	// gets an iterator for all timeslices from time 0 until the current Blockheader time
	iterator := k.RedelegationQueueIterator(ctx, ctx.BlockHeader().Time)
	for ; iterator.Valid(); iterator.Next() {
		timeslice := types.DVVTriplets{}
		value := iterator.Value()
		k.cdc.MustUnmarshalBinaryLengthPrefixed(value, &timeslice)

		matureRedelegations = append(matureRedelegations, timeslice.Triplets...)

		realKey := iterator.Key()
		store.Delete(realKey)
	}

	return matureRedelegations
}

// Returns all the redelegation queue timeslices from time 0 until endTime
func (k StakingKeeper) RedelegationQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	prefix := types.GetRedelegationTimeKey(endTime)
	store := ctx.KVStore(k.storeKey)
	iterator := store.RemoteIterator(prefix, sdk.PrefixEndBytes(prefix))
	if !iterator.Valid() {
		iterator.Close()
		store := ctx.KVStore(k.storeKey)
		iterator = sdk.KVStorePrefixIterator(store, prefix)
	}
	return iterator
}

// CompleteRedelegationWithAmount completes the redelegations of all mature entries in the
// retrieved redelegation object and returns the total redelegation (initial)
// balance or an error upon failure.
func (k StakingKeeper) CompleteRedelegationWithAmount(
	ctx sdk.Context, delAddr sdk.AccAddress, valSrcAddr, valDstAddr sdk.AccAddress,
) (sdk.Coins, error) {

	red, found := k.GetRedelegation(ctx, delAddr, valSrcAddr, valDstAddr)
	if !found {
		return nil, errors.New("no ReDelegation exist")
	}

	//bondDenom := k.GetParams(ctx).BondDenom
	balances := sdk.NewCoins()
	ctxTime := ctx.BlockHeader().Time

	// loop through all the entries and complete mature redelegation entries
	for i := 0; i < len(red.Entries); i++ {
		entry := red.Entries[i]
		if entry.IsMature(ctxTime) {
			red.RemoveEntry(int64(i))
			i--

			if !entry.InitialBalance.IsZero() {
				balances = balances.Add(sdk.NewCoins(sdk.NewChainCoin(entry.InitialBalance)))
			}
		}
	}

	// set the redelegation or remove it if there are no more entries
	if len(red.Entries) == 0 {
		k.RemoveRedelegation(ctx, red)
	} else {
		k.SetRedelegation(ctx, red)
	}

	return balances, nil
}

// CompleteUnbondingWithAmount completes the unbonding of all mature entries in
// the retrieved unbonding delegation object and returns the total unbonding
// balance or an error upon failure.
func (k StakingKeeper) CompleteUnbondingWithAmount(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.AccAddress) (sdk.Coins, error) {
	ubd, found := k.GetUnbondingDelegation(ctx, delAddr, valAddr)
	if !found {
		return nil, errors.New("no UnBonding delegation")
	}

	//bondDenom := k.GetParams(ctx).BondDenom
	balances := sdk.NewCoins()
	ctxTime := ctx.BlockHeader().Time

	// loop through all the entries and complete unbonding mature entries
	for i := 0; i < len(ubd.Entries); i++ {
		entry := ubd.Entries[i]
		if entry.IsMature(ctxTime) {
			ubd.RemoveEntry(int64(i))
			i--

			// track undelegation only when remaining or truncated shares are non-zero
			if !entry.Balance.IsZero() {
				amt := sdk.NewChainCoin(entry.Balance)
				err := k.SupplyKeeper.UndelegateCoinsFromModuleToAccount(
					ctx, types.NotBondedPoolName, ubd.DelegatorAddress, amt,//sdk.NewCoin(amt),
				)
				if err != nil {
					return nil, err
				}

				balances = balances.Add(sdk.NewCoins(amt))
			}
		}
	}

	// set the unbonding delegation or remove it if there are no more entries
	if len(ubd.Entries) == 0 {
		k.RemoveUnbondingDelegation(ctx, ubd)
	} else {
		k.SetUnbondingDelegation(ctx, ubd)
	}

	return balances, nil
}

// Returns all the unbonding queue timeslices from time 0 until endTime
func (k StakingKeeper) UBDQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	prefix := types.GetUnbondingDelegationTimeKey(endTime)
	store := ctx.KVStore(k.storeKey)
	iterator := store.RemoteIterator(prefix, sdk.PrefixEndBytes(prefix))
	if !iterator.Valid() {
		iterator.Close()
		store := ctx.KVStore(k.storeKey)
		iterator = sdk.KVStoreReversePrefixIterator(store, prefix)
	}
	return iterator
}

// Returns a concatenated list of all the timeslices inclusively previous to
// currTime, and deletes the timeslices from the queue
func (k StakingKeeper) DequeueAllMatureUBDQueue(ctx sdk.Context, currTime time.Time) (matureUnbonds []types.DVPair) {
	store := ctx.KVStore(k.storeKey)

	// gets an iterator for all timeslices from time 0 until the current Blockheader time
	iterator := k.UBDQueueIterator(ctx, ctx.BlockHeader().Time)
	for ; iterator.Valid(); iterator.Next() {
		timeslice := types.DVPairs{}
		value := iterator.Value()
		k.cdc.MustUnmarshalBinaryLengthPrefixed(value, &timeslice)

		matureUnbonds = append(matureUnbonds, timeslice.Pairs...)

		realKey := iterator.Key()
		//_, ok := iterator.(*couchdb.CouchIterator)
		//if ok {
		//	realKey = sdk.GetRealKey(iterator.Key())
		//}
		store.Delete(realKey)
	}

	return matureUnbonds
}

// remove a redelegation object and associated index
func (k StakingKeeper) RemoveRedelegation(ctx sdk.Context, red types.Redelegation) {
	store := ctx.KVStore(k.storeKey)
	redKey := types.GetREDKey(red.DelegatorAddress, red.ValidatorSrcAddress, red.ValidatorDstAddress)
	store.Delete(redKey)
	store.Delete(types.GetREDByValSrcIndexKey(red.DelegatorAddress, red.ValidatorSrcAddress, red.ValidatorDstAddress))
	store.Delete(types.GetREDByValDstIndexKey(red.DelegatorAddress, red.ValidatorSrcAddress, red.ValidatorDstAddress))
}

// return all redelegations from a particular validator
func (k StakingKeeper) GetRedelegationsFromSrcValidator(ctx sdk.Context, valAddr sdk.AccAddress) (reds []types.Redelegation) {
	store := ctx.KVStore(k.storeKey)
	prefix := types.GetREDsFromValSrcIndexKey(valAddr)
	iterator := store.RemoteIterator(prefix, sdk.PrefixEndBytes(prefix))
	if !iterator.Valid() {
		iterator.Close()
		iterator = sdk.KVStorePrefixIterator(store, prefix)
	}
	defer iterator.Close()
	var red types.Redelegation

	for ; iterator.Valid(); iterator.Next() {
		realKey := iterator.Key()
		//_, ok := iterator.(*couchdb.CouchIterator)
		//if ok {
		//	realKey = sdk.GetRealKey(iterator.Key())
		//}
		key := types.GetREDKeyFromValSrcIndexKey(realKey)
		value := store.Get(key)
		types.StakingCodec.MustUnmarshalBinaryLengthPrefixed(value, &red)
		reds = append(reds, red)
	}
	return reds
}

// return all delegations to a specific validator. Useful for querier.
func (k StakingKeeper) GetValidatorDelegations(ctx sdk.Context, valAddr sdk.AccAddress) (delegations []types.Delegation) { //nolint:interfacer
	prefix := types.DelegationKey
	store := ctx.KVStore(k.storeKey)
	iterator := store.RemoteIterator(prefix, sdk.PrefixEndBytes(prefix))
	if !iterator.Valid() {
		iterator.Close()
		store := ctx.KVStore(k.storeKey)
		iterator = sdk.KVStorePrefixIterator(store, prefix)
	}
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var delegation types.Delegation
		//delegation := types.MustUnmarshalDelegation(k.cdc, iterator.Value())
		types.StakingCodec.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &delegation)
		if delegation.GetValidatorAddr().Equals(valAddr) {
			delegations = append(delegations, delegation)
		}
	}
	return delegations
}

// rearranges the ValIndexKey to get the UBDKey
func GetUBDKeyFromValIndexKey(indexKey []byte) []byte {
	addrs := indexKey[1:] // remove prefix bytes
	if len(addrs) != 2*sdk.AddrLen {
		panic("unexpected key length")
	}

	valAddr := addrs[:sdk.AddrLen]
	delAddr := addrs[sdk.AddrLen:]

	return types.GetUBDKey(sdk.ToAccAddress(delAddr), sdk.ToAccAddress(valAddr))
}

// return all unbonding delegations from a particular validator
func (k StakingKeeper) GetUnbondingDelegationsFromValidator(ctx sdk.Context, valAddr sdk.AccAddress) (ubds []types.UnbondingDelegation) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.GetUBDsByValIndexKey(valAddr))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		key := GetUBDKeyFromValIndexKey(iterator.Key())
		value := store.Get(key)
		ubd := types.MustUnmarshalUBD(k.cdc, value)
		ubds = append(ubds, ubd)
	}

	return ubds
}