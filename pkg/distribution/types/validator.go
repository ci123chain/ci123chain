package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
)

type ValidatorAccumulatedCommission struct {
	Commission sdk.DecCoin   `json:"commission"`
}


type ValidatorCurrentRewards struct {
	Rewards   sdk.DecCoin     `json:"rewards"`
	Period    uint64          `json:"period"`
}

type ValidatorOutstandingRewards struct {
	Rewards   sdk.DecCoin     `json:"rewards"`
}

type ValidatorHistoricalRewards struct {
	CumulativeRewardRatio     sdk.DecCoin     `json:"cumulative_reward_ratio"`
	ReferenceCount            uint32		  `json:"reference_count"`
}

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