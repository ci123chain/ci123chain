package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
)

//验证者的佣金
type ValidatorAccumulatedCommission struct {
	Commission sdk.DecCoin   `json:"commission"`
}

//Period:   时期；
//Rewards:  奖金数量；累积值；每次的提取奖金操作会更新Period； 每次+1；
type ValidatorCurrentRewards struct {
	Rewards   sdk.DecCoin     `json:"rewards"`
	Period    uint64          `json:"period"`
}

//validator总的获得的奖金数；包括佣金，奖金
type ValidatorOutstandingRewards struct {
	Rewards   sdk.DecCoin     `json:"rewards"`
}

//以validator地址 和 period时期 作为key来存储；表示到某个时期为止，这个validator所获得的所有奖金；
// CumulativeRewardRatio: 累积的奖金数（表示的是每一个token所拥有的奖金数；）
type ValidatorHistoricalRewards struct {
	CumulativeRewardRatio     sdk.DecCoin     `json:"cumulative_reward_ratio"`
	ReferenceCount            uint32		  `json:"reference_count"`
}
//惩罚
type ValidatorSlashEvent struct {
	ValidatorPeriod           uint64         `json:"validator_period"`
	Fraction                  sdk.Dec        `json:"fraction"`
}

func NewValidatorHistoricalRewards(cumulativeRewardRatio sdk.DecCoin, referenceCount uint32) ValidatorHistoricalRewards{
	return ValidatorHistoricalRewards{
		CumulativeRewardRatio: cumulativeRewardRatio,
		ReferenceCount:        referenceCount,
	}
}

func NewValidatorCurrentRewards(rewards sdk.DecCoin, period uint64) ValidatorCurrentRewards{
	return ValidatorCurrentRewards{
		Rewards: rewards,
		Period:  period,
	}
}

// return the initial accumulated commission (zero)
func InitialValidatorAccumulatedCommission() ValidatorAccumulatedCommission {
	return ValidatorAccumulatedCommission{
		Commission:sdk.NewEmptyDecCoin(),
	}
}