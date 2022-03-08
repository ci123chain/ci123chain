package pre_staking

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/pre_staking/keeper"
	"github.com/ci123chain/ci123chain/pkg/pre_staking/types"
	"math/big"
)

const (
	CheckRecordWindow = 20000
	baseMonth         = 720
)

func BeginBlock() {}

func EndBlock(ctx sdk.Context, k keeper.PreStakingKeeper) {
	UpdateDeadlineRecord(ctx, k)
}

func UpdateDeadlineRecord(ctx sdk.Context, ps keeper.PreStakingKeeper) {
	iterator := ps.Iter(ctx)
	prune := ctx.BlockHeight()%CheckRecordWindow == 0

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
				var z = sv.Amount.Amount.BigInt()
				base := int64(sv.StorageTime.Hours()) / baseMonth
				burnTokens := z.Mul(z, big.NewInt(base))

				if tokenmanager := ps.GetTokenManager(ctx); len(tokenmanager) > 0 {
					err = ps.SupplyKeeper.BurnEVMCoin(ctx, types.ModuleName, sdk.HexToAddress(tokenmanager), sv.Delegator, burnTokens)
					if err != nil {
						ctx.Logger().Warn("Burn evm coin failed")
					}
				} else {
					ctx.Logger().Warn("StakingToken Address not set")
				}
				ps.UpdateStakingRecordProcessed(ctx, iterator.Key())
			}

			if prune && sv.Processed{
				ps.DeleteStakingVault(ctx, iterator.Key())
			}
		}
	}
}
