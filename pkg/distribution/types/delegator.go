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