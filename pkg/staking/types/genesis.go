package types

import (
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
)

type GenesisState struct {
	Params               Params                 `json:"params"`
	LastTotalPower       sdk.Int                `json:"last_total_power"`
	LastValidatorPowers  []LastValidatorPower   `json:"last_validator_powers"`
	Validators           Validators             `json:"validators"`
	Delegations          Delegations            `json:"delegations"`
	UnbondingDelegations []UnbondingDelegation  `json:"unbonding_delegations"`
	Redelegations        []Redelegation         `json:"redelegations"`
	Exported             bool                   `json:"exported"`

}


// LastValidatorPower required for validator set update logic
type LastValidatorPower struct {
	Address sdk.AccAddress
	Power   int64
}

// NewGenesisState creates a new GenesisState instanc e
func NewGenesisState(params Params, validators []Validator, delegations []Delegation) GenesisState {
	return GenesisState{
		Params:      params,
		Validators:  validators,
		Delegations: delegations,
	}
}

// DefaultGenesisState gets the raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{Params:DefaultParams()}
}
