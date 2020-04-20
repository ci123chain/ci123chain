package distribution

import (
	"encoding/json"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	k "github.com/ci123chain/ci123chain/pkg/distribution/keeper"
	abci "github.com/tendermint/tendermint/abci/types"
)

func BeginBlock(ctx types.Context, req abci.RequestBeginBlock, distr k.DistrKeeper) {

	if ctx.BlockHeight() > ModuleHeight {
		height := ctx.BlockHeight()
		fee := distr.FeeCollectionKeeper.GetCollectedFees(ctx)
		//分配完奖励金之后清空奖金池
		distr.FeeCollectionKeeper.ClearCollectedFees(ctx)
		//分给validators
		SaveValidatorsInfo(ctx, req, distr, height, fee)
		proposerAddress := distr.GetPreviousProposer(ctx)

		rewards := distr.GetProposerCurrentRewards(ctx, proposerAddress, height - ModuleHeight)
		rewards = rewards.Add(fee)
		distr.SetProposerCurrentRewards(ctx, proposerAddress, rewards, height)

		//store记录的数据维持在100个块内
		h := height - block
		if h >= ModuleHeight {
			DeleteHistoricalRewards(ctx, distr, h)
		}
	}
	proposerAddr := types.AccAddr(req.Header.ProposerAddress)
	//存储proposeAddress
	distr.SetPreviousProposer(ctx, proposerAddr)
}

func SaveValidatorsInfo(ctx types.Context, req abci.RequestBeginBlock, distr k.DistrKeeper,height int64, fee types.Coin) {
	votes := req.LastCommitInfo.Votes
	length := len(votes)
	//fmt.Println(length)
	var validatorAddresses lastCommitValidatorsAddr
	if length == 1 {
		addr := votes[0].Validator.Address
		valAddress := types.AccAddr(addr)
		//发放奖金给validators
		valRewards := distr.GetValidatorCurrentRewards(ctx, valAddress, height - ModuleHeight)
		valRewards = valRewards.Add(fee)
		distr.SetValidatorCurrentRewards(ctx, valAddress, valRewards, height)
	}else {
		for i := 0; i < length; i++ {
			addr := votes[i].Validator.Address
			address := fmt.Sprintf("%X", addr)
			valAddress := types.AccAddr(addr)
			//发放奖金给validators
			valRewards := distr.GetValidatorCurrentRewards(ctx, valAddress, height - ModuleHeight)
			valRewards = valRewards.Add(fee)
			distr.SetValidatorCurrentRewards(ctx, valAddress, valRewards, height)
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

func DeleteHistoricalRewards(ctx types.Context, distr k.DistrKeeper, height int64) {
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
}