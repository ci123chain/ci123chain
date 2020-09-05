package distribution

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/account"
	"github.com/ci123chain/ci123chain/pkg/distribution/keeper"
	"github.com/ci123chain/ci123chain/pkg/distribution/types"
	"github.com/ci123chain/ci123chain/pkg/supply"
)


func InitGenesis(ctx sdk.Context, ak account.AccountKeeper, sk supply.Keeper, k keeper.DistrKeeper, data types.GenesisState) {
	//
	var moduleHoldings = sdk.NewDecCoin(sdk.DefaultCoinDenom, sdk.NewInt(0))
	k.SetFeePool(ctx, data.FeePool)
	k.SetParams(ctx, data.Params)

	for _, dwi := range data.DelegatorWithdrawInfos {
		k.SetDelegatorWithdrawAddr(ctx, dwi.DelegatorAddress, dwi.WithdrawAddress)
	}
	k.SetPreviousProposerAddr(ctx, data.PreviousProposer)
	if data.OutstandingRewards != nil {
		for _, rew := range data.OutstandingRewards {
			k.SetValidatorOutstandingRewards(ctx, rew.ValidatorAddress, types.ValidatorOutstandingRewards{Rewards: rew.OutstandingRewards})
			moduleHoldings = moduleHoldings.Add(rew.OutstandingRewards)
		}
	}
	if data.ValidatorAccumulatedCommissions != nil {
		for _, acc := range data.ValidatorAccumulatedCommissions {
			k.SetValidatorAccumulatedCommission(ctx, acc.ValidatorAddress, acc.Accumulated)
		}
	}

	if data.ValidatorHistoricalRewards != nil {
		for _, his := range data.ValidatorHistoricalRewards {
			k.SetValidatorHistoricalRewards(ctx, his.ValidatorAddress, his.Period, his.Rewards)
		}
	}
	if data.ValidatorCurrentRewards != nil {
		for _, cur := range data.ValidatorCurrentRewards {
			fmt.Println("distr-genesis")
			k.SetValidatorCurrentRewards(ctx, cur.ValidatorAddress, cur.Rewards)
		}
	}
	if data.DelegatorStartingInfos != nil {
		for _, del := range data.DelegatorStartingInfos {
			k.SetDelegatorStartingInfo(ctx, del.ValidatorAddress, del.DelegatorAddress, del.StartingInfo)
		}
	}
	if data.ValidatorSlashEvents != nil {
		for _, evt := range data.ValidatorSlashEvents {
			k.SetValidatorSlashEvent(ctx, evt.ValidatorAddress, evt.Height, evt.Period, evt.Event)
		}
	}
/*
	moduleHoldings = moduleHoldings.Add(data.FeePool.CommunityPool)
	moduleHoldingsInt, _ := moduleHoldings.TruncateDecimal()

	// check if the module account exists
	moduleAcc := k.GetDistributionAccount(ctx)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	if ak.GetAllBalances(ctx, moduleAcc.GetAddress()).IsZero() {
		if err := ak.SetBalances(ctx, moduleAcc.GetAddress(), sdk.NewCoins(moduleHoldingsInt)); err != nil {
			panic(err)
		}

		sk.SetModuleAccount(ctx, moduleAcc)
	}
	*/

}