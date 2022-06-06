package keeper

import (
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
	stake := val.TokensFromShares(del.GetShares())
	// calculate rewards for final period
	rewards = rewards.Add(k.calculateDelegationRewardsBetween(ctx, val, startingPeriod, endingPeriod, stake))
	return rewards
}

func (k DistrKeeper) withdrawDelegationRewards(ctx sdk.Context, val exported.ValidatorI, del exported.DelegationI) (sdk.Coin, error) {

	if !k.HasDelegatorStartingInfo(ctx, val.GetOperator(), del.GetDelegatorAddr()) {
		return sdk.Coin{}, types.ErrEmptyDelegationStartingInfo
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
