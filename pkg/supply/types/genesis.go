package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
)


type GenesisState struct {
	Supply     sdk.Coin    `json:"supply"`
}


// NewGenesisState creates a new genesis state.
func NewGenesisState(supply sdk.Coin) GenesisState {
	return GenesisState{supply}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(DefaultSupply().GetTotal())
}