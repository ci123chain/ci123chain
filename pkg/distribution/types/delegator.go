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

type DelegatorAccountInfo struct {
	Available      sdk.Coin   `json:"available"`  ///本账户余额
	Delegated      sdk.Coin   `json:"delegated"`  ///本账户所有的抵押数
	Unbonding      sdk.Coin   `json:"unbonding"`  ///当前处于解押时期中的资金总数
	Reward         sdk.Coin   `json:"reward"`     ///目前可提取的奖励总数（从上次提取奖励到现在）
	Commission     sdk.Coin   `json:"commission"` ///本验证者账户目前可提取的佣金数（从上次提取佣金到现在，不是验证者账户则为0）
}

func NewDelegatorAccountInfo(available, delegated, unbonding, reward, commission sdk.Coin) DelegatorAccountInfo {
	return DelegatorAccountInfo{
		Available:  available,
		Delegated:  delegated,
		Unbonding:  unbonding,
		Reward:     reward,
		Commission: commission,
	}
}