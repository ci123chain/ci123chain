package keeper

import (
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/staking/types"
	"time"
)

func (k StakingKeeper) UnbondingTime(ctx sdk.Context) (res time.Duration) {
	k.paramstore.Get(ctx, types.KeyUnbondingTime, &res)
	return
}

func (k StakingKeeper) MaxValidators(ctx sdk.Context) (res uint32) {
	k.paramstore.Get(ctx, types.KeyMaxValidators, &res)
	return
}

func (k StakingKeeper) MaxEntries(ctx sdk.Context) (res uint32) {
	k.paramstore.Get(ctx, types.KeyMaxEntries, &res)
	return
}

func (k StakingKeeper) HistoricalEntries(ctx sdk.Context) (res uint32) {
	k.paramstore.Get(ctx, types.KeyHistoricalEntries, &res)
	return
}

func (k StakingKeeper) BondDenom(ctx sdk.Context) (res string) {
	k.paramstore.Get(ctx, types.KeyBondDenom, &res)
	return
}

func (k StakingKeeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.UnbondingTime(ctx),
		k.MaxValidators(ctx),
		k.MaxEntries(ctx),
		k.HistoricalEntries(ctx),
		k.BondDenom(ctx),
		)
}

func (k StakingKeeper) SetParams(ctx sdk.Context, params types.Params) {
	//
	k.paramstore.SetParamSet(ctx, &params)
}