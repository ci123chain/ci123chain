package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
	"time"
)

type Delegation struct {
	StartTime     time.Time        `json:"start_time"`
	StorageTime   time.Time        `json:"storage_time"`
	EndTime       time.Time        `json:"end_time"`
	Amount        sdk.Coin         `json:"amount"`
	Validator     sdk.AccAddress   `json:"validator"`
}

type DelegationRecord struct {
	Delegator     sdk.AccAddress   `json:"delegator"`
	//Validator     sdk.AccAddress   `json:"validator"`
	Delegations   []Delegation     `json:"delegations"`
	PrestakingAmount  sdk.Int           `json:"prestaking_amount"`
}


type GenesisState struct {
	DaoDeployed  bool                   `json:"dao_deployed"`
	Records      []DelegationRecord     `json:"records"`
}

func NewGenesisState(records []DelegationRecord, deployed bool) GenesisState {
	return GenesisState{
		Records: records,
		DaoDeployed: deployed,
	}
}

func DefaultGenesisState(_ []tmtypes.GenesisValidator) GenesisState {
	return NewGenesisState(nil, false)
}