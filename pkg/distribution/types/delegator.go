package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
)

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