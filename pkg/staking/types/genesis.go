package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
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
	Address sdk.AccAddress   `json:"address"`
	Power   int64             `json:"power"`
}

// NewGenesisState creates a new GenesisState instanc e
func NewGenesisState(params Params, validators []Validator, delegations []Delegation) GenesisState {
	return GenesisState{
		Params:      params,
		Validators:  validators,
		Delegations: delegations,
	}
}

func DefaultValidators(validators []tmtypes.GenesisValidator, accAddresses []string) []Validator {
	var genesisValidators []Validator
	var genesisValidator Validator
	if len(validators) == 0 {
		return genesisValidators
	}else {
		for i := range validators {
			genesisValidator.OperatorAddress = sdk.HexToAddress(accAddresses[i])
			//genesisValidator.Address = validators[i].PubKey.Address()
			genesisValidator.ConsensusKey = validators[i].PubKey
			genesisValidator.Status = 1
			genesisValidator.Tokens = sdk.NewInt(400)
			genesisValidator.DelegatorShares = sdk.NewDecFromInt(genesisValidator.Tokens)
			genesisValidator.MinSelfDelegation = sdk.NewInt(400)
			genesisValidator.Commission = NewCommission(sdk.NewDecWithPrec(1, 2), sdk.NewDecWithPrec(4, 2), sdk.NewDecWithPrec(5, 1))
			genesisValidators = append(genesisValidators, genesisValidator)
		}
	}
	return genesisValidators
}

func MinselfDelegation(validators []Validator) (delegations []Delegation){
	if len(validators) == 0 {
		return delegations
	}else {
		for i := range validators {
			delegation := Delegation{
				DelegatorAddress: validators[i].OperatorAddress,
				ValidatorAddress: validators[i].OperatorAddress,
				Shares:           sdk.NewDecFromInt(validators[i].Tokens),
			}
			delegations = append(delegations, delegation)
		}
		return delegations
	}
}

// DefaultGenesisState gets the raw genesis raw message for testing
func DefaultGenesisState(validators []tmtypes.GenesisValidator, accAddresses []string) GenesisState {
	genesisValidators := DefaultValidators(validators, accAddresses)
	return NewGenesisState(DefaultParams(), genesisValidators, MinselfDelegation(genesisValidators))
}
