package distribution

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	k "github.com/ci123chain/ci123chain/pkg/distribution/keeper"
	abci "github.com/tendermint/tendermint/abci/types"
	"log"
)

func BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock, distr k.DistrKeeper) {

	var previousTotalPower, sumPreviousPrecommitPower int64
	for _, voteInfo := range req.LastCommitInfo.GetVotes() {
		previousTotalPower += voteInfo.Validator.Power
		if voteInfo.SignedLastBlock {
			sumPreviousPrecommitPower += voteInfo.Validator.Power
		}
	}

	if ctx.BlockHeight() > ModuleHeight {
		//height := ctx.BlockHeight()

		feeCollector := distr.SupplyKeeper.GetModuleAccount(ctx, distr.FeeCollectorName)
		feeCollectedInt := sdk.NewChainCoin(distr.AccountKeeper.GetBalance(ctx, feeCollector.GetAddress()).AmountOf(sdk.ChainCoinDenom))
		feeCollected := sdk.NewDecCoinFromCoin(feeCollectedInt)

		err := distr.SupplyKeeper.SendCoinsFromModuleToModule(ctx, distr.FeeCollectorName, ModuleName, feeCollectedInt)
		if err != nil {
			panic(err)
		}

		feePool := distr.GetFeePool(ctx)
		if previousTotalPower == 0 {
			feePool.CommunityPool = feePool.CommunityPool.Add(feeCollected)
			distr.SetFeePool(ctx, feePool)
			return
		}
		// calculate fraction votes
		previousFractionVotes := sdk.NewDec(sumPreviousPrecommitPower).Quo(sdk.NewDec(previousTotalPower))
		// calculate previous proposer reward
		baseProposerReward := distr.GetBaseProposerReward(ctx)
		bonuseProposerReward := distr.GetBonusProposerReward(ctx)

		proposerMultiplier := baseProposerReward.Add(bonuseProposerReward.MulTruncate(previousFractionVotes))
		proposerReward := feeCollected.MulDecTruncate(proposerMultiplier)

		// pay previous proposer
		remainning := feeCollected
		//拿到validator
		proposerAddress := distr.GetPreviousProposerAddr(ctx)
		proposerValidator, found := distr.StakingKeeper.GetValidatorByConsAddr(ctx, sdk.ToAccAddress(proposerAddress))
		if found {
			distr.AllocateTokensToValidator(ctx, proposerValidator, proposerReward)
			remainning = remainning.Sub(proposerReward)
		}else {
			log.Println(fmt.Sprintf(
				"WARNING: Attempt to allocate proposer rewards to unknown proposer %s. "+
					"This should happen only if the proposer unbonded completely within a single block, "+
					"which generally should not happen except in exceptional circumstances (or fuzz testing). "+
					"We recommend you investigate immediately.",
				sdk.ToAccAddress(proposerAddress).String()))
		}

		// calculate fraction allocated to validators
		communityTax := distr.GetCommunityTax(ctx)
		voteMultiplier := sdk.OneDec().Sub(proposerMultiplier).Sub(communityTax)

		previousVotes := req.LastCommitInfo.GetVotes()

		for _, vote := range previousVotes {
			//拿到validator.
			validator, _ := distr.StakingKeeper.GetValidatorByConsAddr(ctx, sdk.ToAccAddress(vote.Validator.Address))

			powerFraction := sdk.NewDec(vote.Validator.Power).QuoTruncate(sdk.NewDec(previousTotalPower))
			reward := feeCollected.MulDecTruncate(voteMultiplier).MulDecTruncate(powerFraction)
			distr.AllocateTokensToValidator(ctx, validator, reward)
			remainning = remainning.Sub(reward)
		}

		// allocate community funding
		feePool.CommunityPool = feePool.CommunityPool.Add(remainning)
		distr.SetFeePool(ctx, feePool)
	}
	proposerAddr := sdk.AccAddr(req.Header.ProposerAddress)
	//存储proposeAddress
	distr.SetPreviousProposerAddr(ctx, proposerAddr)
}