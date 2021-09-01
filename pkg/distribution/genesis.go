package distribution

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/account"
	"github.com/ci123chain/ci123chain/pkg/distribution/keeper"
	"github.com/ci123chain/ci123chain/pkg/distribution/types"
	"github.com/ci123chain/ci123chain/pkg/supply"
)


func InitGenesis(ctx sdk.Context, ak account.AccountKeeper, sk supply.Keeper, k keeper.DistrKeeper, data types.GenesisState) {
	//
	var moduleHoldings = sdk.NewDecCoin(sdk.ChainCoinDenom, sdk.NewInt(0))
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

func ExportGenesis(ctx sdk.Context, ak account.AccountKeeper, sk supply.Keeper, k keeper.DistrKeeper) GenesisState {

	params := k.GetParams(ctx)
	feePool := k.GetFeePool(ctx)

	dwi := make([]types.DelegatorWithdrawInfo, 0)
	k.IterateDelegatorWithdrawAddrs(ctx, func(del sdk.AccAddress, addr sdk.AccAddress) (stop bool) {
		dwi = append(dwi, types.DelegatorWithdrawInfo{
			DelegatorAddress: del,
			WithdrawAddress:  addr,
		})
		return false
	})

	pp := k.GetPreviousProposerConsAddr(ctx)
	outstanding := make([]types.ValidatorOutstandingRewardsRecord, 0)

	rs := make([]types.ValidatorOutstandingRewardsRecord, 0)
	k.IterateValidatorOutstandingRewards(ctx,
		func(addr sdk.AccAddress, rewards types.ValidatorOutstandingRewards) (stop bool) {
			outstanding = append(outstanding, types.ValidatorOutstandingRewardsRecord{
				ValidatorAddress:   addr,
				OutstandingRewards: rewards.Rewards,
			})
			return false
		},
	)

	acc := make([]types.ValidatorAccumulatedCommissionRecord, 0)
	k.IterateValidatorAccumulatedCommissions(ctx,
		func(addr sdk.AccAddress, commission types.ValidatorAccumulatedCommission) (stop bool) {
			acc = append(acc, types.ValidatorAccumulatedCommissionRecord{
				ValidatorAddress: addr,
				Accumulated:      commission,
			})
			return false
		},
	)

	his := make([]types.ValidatorHistoricalRewardsRecord, 0)
	k.IterateValidatorHistoricalRewards(ctx,
		func(val sdk.AccAddress, period uint64, rewards types.ValidatorHistoricalRewards) (stop bool) {
			his = append(his, types.ValidatorHistoricalRewardsRecord{
				ValidatorAddress: val,
				Period:           period,
				Rewards:          rewards,
			})
			return false
		},
	)

	cur := make([]types.ValidatorCurrentRewardsRecord, 0)
	k.IterateValidatorCurrentRewards(ctx,
		func(val sdk.AccAddress, rewards types.ValidatorCurrentRewards) (stop bool) {
			cur = append(cur, types.ValidatorCurrentRewardsRecord{
				ValidatorAddress: val,
				Rewards:          rewards,
			})
			return false
		},
	)

	dels := make([]types.DelegatorStartingInfoRecord, 0)
	k.IterateDelegatorStartingInfos(ctx,
		func(val sdk.AccAddress, del sdk.AccAddress, info types.DelegatorStartingInfo) (stop bool) {
			dels = append(dels, types.DelegatorStartingInfoRecord{
				ValidatorAddress: val,
				DelegatorAddress: del,
				StartingInfo:     info,
			})
			return false
		},
	)

	slashes := make([]types.ValidatorSlashEventRecord, 0)
	k.IterateValidatorSlashEvents(ctx,
		func(val sdk.AccAddress, height uint64, event types.ValidatorSlashEvent) (stop bool) {
			slashes = append(slashes, types.ValidatorSlashEventRecord{
				ValidatorAddress:    val,
				Height:              height,
				Period:              event.ValidatorPeriod,
				Event: event,
			})
			return false
		},
	)

	return types.NewGenesisState(params, feePool, dwi, pp.Bytes(), rs, acc, his, cur, dels, slashes)
}