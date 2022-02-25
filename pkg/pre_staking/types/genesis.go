package types

import (
	tmtypes "github.com/tendermint/tendermint/types"
)

type DelegationRecord struct {
	StakingRecord      []StakingVault      `json:"prestaking_record"`
}

type GenesisState struct {
	StakingToken   	string             	`json:"staking_token"`
	Owner 			string 				`json:"owner"`
	Records      	DelegationRecord    `json:"records"`
}

func NewGenesisState(records DelegationRecord, addr string, owner string) GenesisState {
	return GenesisState{
		Records: records,
		StakingToken: addr,
		Owner: owner,
	}
}

func DefaultGenesisState(_ []tmtypes.GenesisValidator) GenesisState {
	return NewGenesisState(DelegationRecord{
		nil,
	}, "", "")
}