package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
)

type QueryPreStakingRecord struct {
	Delegator  sdk.AccAddress `json:"delegator"`
}

type QueryPreStakingResult struct {
	Amount  string   `json:"amount"`
}

type QueryStakingRecord struct {
	DelegatorAddr   sdk.AccAddress   `json:"delegator_addr"`
	ValidatorAddr   sdk.AccAddress   `json:"validator_addr"`
}