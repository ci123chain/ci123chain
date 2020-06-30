package keeper

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/distribution/types"
	"github.com/ci123chain/ci123chain/pkg/supply/exported"
)

func (k DistrKeeper) GetDistributionAccount(ctx sdk.Context) exported.ModuleAccountI {
	return k.SupplyKeeper.GetModuleAccount(ctx, types.ModuleName)
}

// get outstanding rewards
func (k DistrKeeper) GetValidatorOutstandingRewardsCoins(ctx sdk.Context, val sdk.AccAddress) sdk.DecCoin {
	return k.GetValidatorOutstandingRewards(ctx, val).Rewards
}


func (k DistrKeeper) GetFeePoolCommunity(ctx sdk.Context) sdk.DecCoin {
	return k.GetFeePool(ctx).CommunityPool
}