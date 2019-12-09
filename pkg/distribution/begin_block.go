package distribution

import (
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
	k "github.com/tanhuiya/ci123chain/pkg/distribution/keeper"
	abci "github.com/tendermint/tendermint/abci/types"
)


func BeginBlock(ctx types.Context, req abci.RequestBeginBlock, distr k.DistrKeeper) {

	if ctx.BlockHeight() > 1 {
		height := ctx.BlockHeight()
		fee := distr.FeeCollectionKeeper.GetCollectedFees(ctx)
		//分配完奖励金之后清空奖金池
		distr.FeeCollectionKeeper.ClearCollectedFees(ctx)
		proposerAddress := distr.GetPreviousProposer(ctx)

		rewards := distr.GetProposerCurrentRewards(ctx, proposerAddress, height - 1)
		rewards = rewards.SafeAdd(fee)
		distr.SetProposerCurrentRewards(ctx, proposerAddress, rewards, height)

	}

	proposerAddr := types.AccAddr(req.Header.ProposerAddress)
	//存储proposeAddress
	distr.SetPreviousProposer(ctx, proposerAddr)
}
