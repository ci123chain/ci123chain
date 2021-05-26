package keeper

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/staking/types"
	"time"
)

func (k StakingKeeper) UnbondingTime(ctx sdk.Context) (unbondingTime time.Duration) {
	k.paramstore.Get(ctx, types.KeyUnbondingTime, &unbondingTime)
	return
}

func (k StakingKeeper) MaxValidators(ctx sdk.Context) (maxValidators uint32) {
	k.paramstore.Get(ctx, types.KeyMaxValidators, &maxValidators)
	return
}

func (k StakingKeeper) MaxEntries(ctx sdk.Context) (maxEntries uint32) {
	k.paramstore.Get(ctx, types.KeyMaxEntries, &maxEntries)
	return
}

func (k StakingKeeper) HistoricalEntries(ctx sdk.Context) (HisEntries uint32) {
	k.paramstore.Get(ctx, types.KeyHistoricalEntries, &HisEntries)
	return
}

func (k StakingKeeper) BondDenom(ctx sdk.Context) (bondDenom string) {
	k.paramstore.Get(ctx, types.KeyBondDenom, &bondDenom)
	return
}

func (k StakingKeeper) GetParams(ctx sdk.Context) types.Params {
	unbondingTime := k.UnbondingTime(ctx)
	maxValidators := k.MaxValidators(ctx)
	maxEntries := k.MaxEntries(ctx)
	hisEntries := k.HistoricalEntries(ctx)
	bonddenom := k.BondDenom(ctx)
	return types.NewParams(
			unbondingTime,
			maxValidators,
			maxEntries,
			hisEntries,
			bonddenom,
		)
}

func (k StakingKeeper) SetParams(ctx sdk.Context, params types.Params) {
	//
	k.paramstore.SetParamSet(ctx, &params)
}