package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
)


//PreviousPeriod：形成当前StartingInfo的时期；PreviousPeriod = CurrentRewards.Period - 1;
//Stake: 该委托者拥有的令牌数量;    stake = Delegator.Shares * validator.Tokens/validator.Shares
//Height: 形成当前StartingInfo的区块高度;
type DelegatorStartingInfo struct {
	PreviousPeriod            uint64         `json:"previous_period"`
	Stake                     sdk.Dec        `json:"stake"`
	Height                    uint64         `json:"height"`
}

// create a new DelegatorStartingInfo
func NewDelegatorStartingInfo(previousPeriod uint64, stake sdk.Dec, height uint64) DelegatorStartingInfo {
	return DelegatorStartingInfo{
		PreviousPeriod: previousPeriod,
		Stake:          stake,
		Height:         height,
	}
}


type RewardAccount struct {
	Amount sdk.Coin  `json:"amount"`
	Address string   `json:"address"`
}

type RewardsAccount struct {
	sdk.Coin
	Validator      []RewardAccount      `json:"validator"`
}

func NewRewardsAccount(amount sdk.Coin,  validators []RewardAccount) RewardsAccount {
	return RewardsAccount{
		Coin: amount,
		Validator: validators,
	}
}

type DelegatorAccountInfo struct {
	Available      sdk.Coin   `json:"available"`  ///本账户余额
	Delegated      sdk.Coin   `json:"delegated"`  ///本账户所有的抵押数
	Unbonding      sdk.Coin   `json:"unbonding"`  ///当前处于解押时期中的资金总数
	Reward         RewardsAccount   `json:"reward"`     ///目前可提取的奖励总数（从上次提取奖励到现在）
	Commission     sdk.Coin   `json:"commission"` ///本验证者账户目前可提取的佣金数（从上次提取佣金到现在，不是验证者账户则为0）
}

func NewDelegatorAccountInfo(available, delegated, unbonding, commission sdk.Coin, reward RewardsAccount) DelegatorAccountInfo {
	d :=  DelegatorAccountInfo{
		Available:  available,
		Delegated:  delegated,
		Unbonding:  unbonding,
		Commission: commission,
		Reward:     reward,
	}
	return d
}