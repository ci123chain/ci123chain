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
		feeCollectedInt := distr.AccountKeeper.GetBalance(ctx, feeCollector.GetAddress())
		//fmt.Printf("feeCollectedInt = %s\n", feeCollectedInt.String())
		feeCollected := sdk.NewDecCoinFromCoin(feeCollectedInt)
		//fmt.Printf("feeCollected = %s\n", feeCollected.String())

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
		//fmt.Printf("proposer multiplier = %s\n", proposerMultiplier.String())
		proposerReward := feeCollected.MulDecTruncate(proposerMultiplier)

		// pay previous proposer
		remainning := feeCollected
		//拿到validator
		proposerAddress := distr.GetPreviousProposerAddr(ctx)
		//fmt.Printf("proposerAddress = %s\n", sdk.ToAccAddress(proposerAddress).String())
		//fmt.Printf("proposer reward = %s\n", proposerReward.Amount.String())
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
			//fmt.Printf("voting validator = %s\n", sdk.ToAccAddress(vote.Validator.Address).String())

			powerFraction := sdk.NewDec(vote.Validator.Power).QuoTruncate(sdk.NewDec(previousTotalPower))
			reward := feeCollected.MulDecTruncate(voteMultiplier).MulDecTruncate(powerFraction)
			//fmt.Printf("voteMultiplier = %s\n", voteMultiplier.String())
			//fmt.Printf("feeCollected = %s\n", feeCollected.String())
			//fmt.Printf("powerFraction = %s\n", powerFraction.String())
			//fmt.Printf("validator reward = %s\n", reward.Amount.String())
			distr.AllocateTokensToValidator(ctx, validator, reward)
			remainning = remainning.Sub(reward)
		}

		/*fee := distr.FeeCollectionKeeper.GetCollectedFees(ctx)
		//分配完奖励金之后清空奖金池
		distr.FeeCollectionKeeper.ClearCollectedFees(ctx)
		//分给validators
		SaveValidatorsInfo(ctx, req, distr, height, fee)


		rewards := distr.GetProposerCurrentRewards(ctx, proposerAddress, height - ModuleHeight)
		rewards = rewards.Add(fee)
		distr.SetProposerCurrentRewards(ctx, proposerAddress, rewards, height)

		//store记录的数据维持在100个块内
		h := height - block
		if h >= ModuleHeight {
			DeleteHistoricalRewards(ctx, distr, h)
		}*/
		// allocate community funding
		feePool.CommunityPool = feePool.CommunityPool.Add(remainning)
		//fmt.Printf("fee pool community = %s\n", remainning.Amount.String())
		distr.SetFeePool(ctx, feePool)
	}
	proposerAddr := sdk.AccAddr(req.Header.ProposerAddress)
	//存储proposeAddress
	distr.SetPreviousProposerAddr(ctx, proposerAddr)
}
/*
func SaveValidatorsInfo(ctx sdk.Context, req abci.RequestBeginBlock, distr k.DistrKeeper,height int64, fee sdk.Coin) {
	votes := req.LastCommitInfo.Votes
	length := len(votes)
	//fmt.Println(length)
	var validatorAddresses lastCommitValidatorsAddr
	if length == 1 {
		addr := votes[0].Validator.Address
		valAddress := sdk.AccAddr(addr)
		//发放奖金给validators
		valRewards := distr.GetValidatorCurrentRewards(ctx, valAddress, height - ModuleHeight)
		valRewards = valRewards.Add(fee)
		distr.SetValidatorCurrentRewards(ctx, valAddress, valRewards, height)
	}else {
		for i := 0; i < length; i++ {
			addr := votes[i].Validator.Address
			address := fmt.Sprintf("%X", addr)
			valAddress := sdk.AccAddr(addr)
			//发放奖金给validators
			valRewards := distr.GetValidatorCurrentRewards(ctx, valAddress, height - ModuleHeight)
			valRewards = valRewards.Add(fee)
			distr.SetValidatorCurrentRewards(ctx, valAddress, valRewards)
			validatorAddresses.Address = append(validatorAddresses.Address, address)
		}
	}
	//存储上个高度的验证者集合的信息
	b, err := json.Marshal(validatorAddresses)
	if err != nil {
		panic(err)
	}
	distr.SetValidatorsInfo(ctx, b, height - ModuleHeight)
}
*/
/*func DeleteHistoricalRewards(ctx sdk.Context, distr k.DistrKeeper, height int64) {
	var valCommitAddresses lastCommitValidatorsAddr
	bz := distr.GetValidatorsInfo(ctx, height)
	if len(bz) < 1 {
		return
	}
	err := json.Unmarshal(bz, &valCommitAddresses)
	if err != nil {
		panic(err)
	}
	length := len(valCommitAddresses.Address)
	for k := 0; k < length; k ++ {
		key := getKey(valCommitAddresses.Address[k], height)
		distr.DeleteValidatorOldRewardsRecord(ctx, key)
	}

	distr.DeleteValidatorsInfo(ctx, height)
}*/