package keeper

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	staking "github.com/ci123chain/ci123chain/pkg/staking/types"
)

func (k DistrKeeper) AllocateTokensToValidator(ctx sdk.Context, val staking.Validator, tokens sdk.DecCoin) {
	//
	commission := tokens.MulDec(val.GetCommission())
	shared := tokens.Sub(commission)

	//update current commission
	currentCommission := k.GetValidatorAccumulatedCommission(ctx, val.GetOperator())
	currentCommission.Commission = currentCommission.Commission.Add(commission)
	k.SetValidatorAccumulatedCommission(ctx, val.GetOperator(), currentCommission)

	//update current reward
	currentReward := k.GetValidatorCurrentRewards(ctx, val.GetOperator())
	if currentReward.Rewards.Amount.IsZero() {
		currentReward.Rewards = shared
	}else {
		currentReward.Rewards = currentReward.Rewards.Add(shared)
	}
	k.SetValidatorCurrentRewards(ctx, val.GetOperator(), currentReward)

	// update outstanding rewards
	outstanding := k.GetValidatorOutstandingRewards(ctx, val.GetOperator())
	if outstanding.Rewards.Amount.IsZero(){
		outstanding.Rewards = tokens
	}else {
		outstanding.Rewards = outstanding.Rewards.Add(tokens)
	}
	k.SetValidatorOutstandingRewards(ctx, val.GetOperator(), outstanding)

}
