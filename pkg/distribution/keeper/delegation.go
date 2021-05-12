package keeper

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/distribution/types"
	"github.com/ci123chain/ci123chain/pkg/staking/exported"
)


// initialize starting info for a new delegation
func (k DistrKeeper) initializeDelegation(ctx sdk.Context, val sdk.AccAddress, del sdk.AccAddress) {
	// period has already been incremented - we want to store the period ended by this delegation action
	previousPeriod := k.GetValidatorCurrentRewards(ctx, val).Period - 1

	// increment reference count for the period we're going to track
	k.incrementReferenceCount(ctx, val, previousPeriod)

	validator := k.StakingKeeper.Validator(ctx, val)
	delegation := k.StakingKeeper.Delegation(ctx, del, val)

	// calculate delegation stake in tokens
	// we don't store directly, so multiply delegation shares * (tokens per share)
	// note: necessary to truncate so we don't allow withdrawing more rewards than owed
	stake := validator.TokensFromSharesTruncated(delegation.GetShares())
	k.SetDelegatorStartingInfo(ctx, val, del, types.NewDelegatorStartingInfo(previousPeriod, stake, uint64(ctx.BlockHeight())))
}

// calculate the rewards accrued by a delegation between two periods
func (k DistrKeeper) calculateDelegationRewardsBetween(ctx sdk.Context, val exported.ValidatorI,
	startingPeriod, endingPeriod uint64, stake sdk.Dec) (rewards sdk.DecCoin) {
	// sanity check
	if startingPeriod > endingPeriod {
		panic("startingPeriod cannot be greater than endingPeriod")
	}

	// sanity check
	if stake.IsNegative() {
		panic("stake should not be negative")
	}

	// return staking * (ending - starting)
	starting := k.GetValidatorHistoricalRewards(ctx, val.GetOperator(), startingPeriod)
	ending := k.GetValidatorHistoricalRewards(ctx, val.GetOperator(), endingPeriod)
	difference := ending.CumulativeRewardRatio.Sub(starting.CumulativeRewardRatio)
	if difference.IsNegative() {
		panic("negative rewards should not be possible")
	}
	// note: necessary to truncate so we don't allow withdrawing more rewards than owed
	rewards = difference.MulDecTruncate(stake)
	return
}

// calculate the total rewards accrued by a delegation
func (k DistrKeeper) calculateDelegationRewards(ctx sdk.Context, val exported.ValidatorI, del exported.DelegationI, endingPeriod uint64) sdk.DecCoin {

	var rewards = sdk.NewEmptyDecCoin()
	// fetch starting info for delegation
	startingInfo := k.GetDelegatorStartingInfo(ctx, del.GetValidatorAddr(), del.GetDelegatorAddr())
	if startingInfo.Height == uint64(ctx.BlockHeight()) {
		// started this height, no rewards yet
		return rewards
	}
	startingPeriod := startingInfo.PreviousPeriod
	stake := startingInfo.Stake

	//startingHeight := startingInfo.Height

	// A total stake sanity check; Recalculated final stake should be less than or
	// equal to current stake here. We cannot use Equals because stake is truncated
	// when multiplied by slash fractions (see above). We could only use equals if
	// we had arbitrary-precision rationals.
	currentStake := val.TokensFromShares(del.GetShares())

	if stake.GT(currentStake) {
		// Account for rounding inconsistencies between:
		//
		//     currentStake: calculated as in staking with a single computation
		//     stake:        calculated as an accumulation of stake
		//                   calculations across validator's distribution periods
		//
		// These inconsistencies are due to differing order of operations which
		// will inevitably have different accumulated rounding and may lead to
		// the smallest decimal place being one greater in stake than
		// currentStake. When we calculated slashing by period, even if we
		// round down for each slash fraction, it's possible due to how much is
		// being rounded that we slash less when slashing by period instead of
		// for when we slash without periods. In other words, the single slash,
		// and the slashing by period could both be rounding down but the
		// slashing by period is simply rounding down less, thus making stake >
		// currentStake
		//
		// A small amount of this error is tolerated and corrected for,
		// however any greater amount should be considered a breach in expected
		// behaviour.
		marginOfErr := sdk.SmallestDec().MulInt64(3)
		if stake.LTE(currentStake.Add(marginOfErr)) {
			stake = currentStake
		} else {
			panic(fmt.Sprintf("calculated final stake for delegator %s greater than current stake"+
				"\n\tfinal stake:\t%s"+
				"\n\tcurrent stake:\t%s",
				del.GetDelegatorAddr(), stake, currentStake))
		}
	}
	// calculate rewards for final period
	rewards = rewards.Add(k.calculateDelegationRewardsBetween(ctx, val, startingPeriod, endingPeriod, stake))
	return rewards
}

func (k DistrKeeper) withdrawDelegationRewards(ctx sdk.Context, val exported.ValidatorI, del exported.DelegationI) (sdk.Coin, error) {

	if !k.HasDelegatorStartingInfo(ctx, val.GetOperator(), del.GetDelegatorAddr()) {
		return sdk.Coin{}, types.ErrEmptyDelegationStartingInfo(fmt.Sprintf("%v has no delegator starting info with validtor %v", val.GetOperator().String(), del.GetDelegatorAddr().String()))
	}

	// end current period and calculate rewards
	endingPeriod := k.incrementValidatorPeriod(ctx, val)

	rewards := k.calculateDelegationRewards(ctx, val, del, endingPeriod)

	outstanding := k.GetValidatorOutstandingRewardsCoins(ctx, del.GetValidatorAddr())

	// truncate coins, return remainder to community pool
	coins, remainder := rewards.TruncateDecimal()

	// add coins to user account
	if !coins.IsZero() {
		withdrawAddr := k.GetDelegatorWithdrawAddr(ctx, del.GetDelegatorAddr())
		err := k.SupplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, withdrawAddr, sdk.NewCoins(coins))
		if err != nil {
			return sdk.NewEmptyCoin(), err
		}
	}

	// update the outstanding rewards and the community pool only if the
	// transaction was successful
	k.SetValidatorOutstandingRewards(ctx, del.GetValidatorAddr(), types.ValidatorOutstandingRewards{Rewards: outstanding.Sub(rewards)})
	feePool := k.GetFeePool(ctx)
	feePool.CommunityPool = feePool.CommunityPool.Add(remainder)
	k.SetFeePool(ctx, feePool)

	// decrement reference count of starting period
	startingInfo := k.GetDelegatorStartingInfo(ctx, del.GetValidatorAddr(), del.GetDelegatorAddr())
	startingPeriod := startingInfo.PreviousPeriod
	k.decrementReferenceCount(ctx, del.GetValidatorAddr(), startingPeriod)

	// remove delegator starting info
	k.DeleteDelegatorStartingInfo(ctx, del.GetValidatorAddr(), del.GetDelegatorAddr())

	return coins, nil
}
