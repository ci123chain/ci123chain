package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

//type Delegation struct {
//	StartTime     time.Time        `json:"start_time"`
//	StorageTime   time.Duration        `json:"storage_time"`
//	EndTime       time.Time        `json:"end_time"`
//	Amount        sdk.Coin         `json:"amount"`
//	Validator     sdk.AccAddress   `json:"validator"`
//}
//
//type DelegationRecord struct {
//	Delegator     sdk.AccAddress   `json:"delegator"`
//	//Validator     sdk.AccAddress   `json:"validator"`
//	Delegations   map[string]Delegation     `json:"delegations"`
//	PrestakingAmount  sdk.Int           `json:"prestaking_amount"`
//}

type InitPrestaking struct {
	Delegator     sdk.AccAddress     `json:"delegator"`
	Staking       VaultRecord		 `json:"staking"`
}


type InitStakingRecords struct {
	Delegator      sdk.AccAddress     `json:"delegator"`
	Validator      sdk.AccAddress     `json:"validator"`
	Records        []StakingRecord   `json:"records"`
}

type DelegationRecord struct {
	PrestakingRecord      []InitPrestaking      `json:"prestaking_record"`
	DelStakingRecords     []InitStakingRecords  `json:"del_staking_records"`
}


type GenesisState struct {
	DaoAddress   string               `json:"dao_deployed"`
	Records      DelegationRecord     `json:"records"`
}

func NewGenesisState(records DelegationRecord, addr string) GenesisState {
	return GenesisState{
		Records: records,
		DaoAddress: addr,
	}
}

func DefaultGenesisState(_ []tmtypes.GenesisValidator) GenesisState {
	return NewGenesisState(DelegationRecord{
		nil, nil,
	}, "")
}