package pre_staking

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/pre_staking/keeper"
	"github.com/ci123chain/ci123chain/pkg/pre_staking/types"
)

const CheckRecordWindow = 20000

func BeginBlock() {}

func EndBlock(ctx sdk.Context, k keeper.PreStakingKeeper) {
	UpdateDeadlineRecord(ctx, k)
}


func UpdateDeadlineRecord(ctx sdk.Context, ps keeper.PreStakingKeeper) {
	iterator := ps.Iter(ctx)
	prune :=  ctx.BlockHeight() % CheckRecordWindow == 0

	for ; iterator.Valid(); iterator.Next() {
		v := iterator.Value()
		if v != nil {
			var sv types.StakingVault
			ps.Cdc.MustUnmarshalBinaryBare(iterator.Value(), &sv)
			if sv.EndTime.Before(ctx.BlockTime()) && !sv.Processed {
				amount, err := ps.RemoveDeadlineDelegationAndWithdraw(ctx, sv.Validator, sv.Delegator, sv.Amount)
				if err != nil {
					panic(err)
				}
				moduleAcc := ps.SupplyKeeper.GetModuleAccount(ctx, types.DefaultCodespace)
				err = ps.AccountKeeper.Transfer(ctx, moduleAcc.GetAddress(), sv.Delegator, sdk.NewCoins(sdk.NewChainCoin(amount)))
				if err != nil {
					panic(err)
				}

				ps.UpdateStakingRecordProcessed(ctx, iterator.Key())
			}

			if prune {
				ps.DeleteStakingVault(ctx, iterator.Key())
			}
		}
	}
}


