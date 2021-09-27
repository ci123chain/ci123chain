package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"math/big"
	"time"
)

type StakingRecord struct {
	VaultID      *big.Int   `json:"vault_id"`
	EndTime      time.Time  `json:"end_time"`
	Amount       sdk.Coin   `json:"amount"`
}


func NewStakingRecord(id *big.Int,et time.Time, amount sdk.Coin) StakingRecord {
	return StakingRecord{
		VaultID:     id,
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
