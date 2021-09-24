package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"time"
)

type StakingRecord struct {
	StorageTime  time.Time  `json:"storage_time"`
	UpdateTime   time.Time   `json:"update_time"`
	EndTime      time.Time  `json:"end_time"`
	Amount       sdk.Coin   `json:"amount"`
}


func NewStakingRecord(st, ut, et time.Time, amount sdk.Coin) StakingRecord {
	return StakingRecord{
		StorageTime: st,
		UpdateTime:  ut,
		EndTime:     et,
		Amount:      amount,
	}
}

type StakingRecords []StakingRecord

func (s StakingRecords) Len() int {
	return len(s)
}

func (s StakingRecords) Less(i, j int) bool {
	return s[i].EndTime.Before(s[j].EndTime)
}

func (s StakingRecords) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
