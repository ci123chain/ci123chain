package keeper


import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/distribution/types"
	"github.com/ci123chain/ci123chain/pkg/staking/exported"
)

func (k DistrKeeper) initializeValidator(ctx sdk.Context, val exported.ValidatorI) {

	k.SetValidatorHistoricalRewards(ctx, val.GetOperator(), 0, types.NewValidatorHistoricalRewards(sdk.NewEmptyDecCoin(), 1))
	k.SetValidatorCurrentRewards(ctx, val.GetOperator(), types.NewValidatorCurrentRewards(sdk.NewEmptyDecCoin(), 1))

	k.SetValidatorAccumulatedCommission(ctx, val.GetOperator(), types.InitialValidatorAccumulatedCommission())

	k.SetValidatorOutstandingRewards(ctx, val.GetOperator(), types.ValidatorOutstandingRewards{Rewards:sdk.NewEmptyDecCoin()})
}

// increment validator period, returning the period just ended
func (k DistrKeeper) incrementValidatorPeriod(ctx sdk.Context, val exported.ValidatorI) uint64 {
	// fetch current rewards
	rewards := k.GetValidatorCurrentRewards(ctx, val.GetOperator())

	// calculate current ratio
	var current sdk.DecCoin
	if val.GetTokens().IsZero() {
		// can't calculate ratio for zero-token validators
		// ergo we instead add to the community pool
		feePool := k.GetFeePool(ctx)
		outstanding := k.GetValidatorOutstandingRewards(ctx, val.GetOperator())
		feePool.CommunityPool = feePool.CommunityPool.Add(rewards.Rewards)
		outstanding.Rewards = outstanding.Rewards.Sub(rewards.Rewards)
		k.SetFeePool(ctx, feePool)
		k.SetValidatorOutstandingRewards(ctx, val.GetOperator(), outstanding)

		current = sdk.NewEmptyDecCoin()
	}else {
		// note: necessary to truncate so we don't allow withdrawing more rewards than owed
		current = rewards.Rewards.QuoDecTruncate(val.GetTokens().ToDec())
	}

	// fetch historical rewards for last period
	historical := k.GetValidatorHistoricalRewards(ctx, val.GetOperator(), rewards.Period-1).CumulativeRewardRatio

	// decrement reference count
	k.decrementReferenceCount(ctx, val.GetOperator(), rewards.Period-1)
	// set new historical rewards with reference count of 1
	k.SetValidatorHistoricalRewards(ctx, val.GetOperator(), rewards.Period, types.NewValidatorHistoricalRewards(historical.Add(current), 1))

	// set current rewards, incrementing period by 1
	k.SetValidatorCurrentRewards(ctx, val.GetOperator(), types.NewValidatorCurrentRewards(sdk.NewEmptyDecCoin(), rewards.Period+1))

	return rewards.Period
}


// increment the reference count for a historical rewards value
func (k DistrKeeper) incrementReferenceCount(ctx sdk.Context, valAddr sdk.AccAddress, period uint64) {
	historical := k.GetValidatorHistoricalRewards(ctx, valAddr, period)
	if historical.ReferenceCount > 2 {
		panic("reference count should never exceed 2")
	}
	historical.ReferenceCount++
	k.SetValidatorHistoricalRewards(ctx, valAddr, period, historical)
}


// decrement the reference count for a historical rewards value, and delete if zero references remain
func (k DistrKeeper) decrementReferenceCount(ctx sdk.Context, valAddr sdk.AccAddress, period uint64) {
	historical := k.GetValidatorHistoricalRewards(ctx, valAddr, period)
	if historical.ReferenceCount == 0 {
		panic("cannot set negative reference count")
	}
	historical.ReferenceCount--
	if historical.ReferenceCount == 0 {
		k.DeleteValidatorHistoricalReward(ctx, valAddr, period)
	} else {
		k.SetValidatorHistoricalRewards(ctx, valAddr, period, historical)
	}
}