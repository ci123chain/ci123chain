package staking

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	k "github.com/ci123chain/ci123chain/pkg/staking/keeper"
	abci "github.com/tendermint/tendermint/abci/types"
)

func BeginBlock(ctx sdk.Context, k k.StakingKeeper) {
	//
	k.TrackHistoricalInfo(ctx)

}

func EndBlock(ctx sdk.Context, k k.StakingKeeper) []abci.ValidatorUpdate {
	//
	return k.BlockValidatorUpdates(ctx)
}
