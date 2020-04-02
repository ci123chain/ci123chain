package staking

import (
	abci "github.com/tendermint/tendermint/abci/types"
	k "github.com/tanhuiya/ci123chain/pkg/staking/keeper"
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
)

func BeginBlock(ctx sdk.Context, k k.StakingKeeper) {
	//
	//k.TrackHistoricalInfo(ctx)

}

func EndBlock(ctx sdk.Context, k k.StakingKeeper) []abci.ValidatorUpdate {
	//
	return k.BlockValidatorUpdates(ctx)
}
