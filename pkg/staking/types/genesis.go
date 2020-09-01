package types

import (
	"encoding/hex"
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

func DefaultValidators(validators []tmtypes.GenesisValidator) []Validator {
	var genesisValidators []Validator
	var genesisValidator Validator
	if len(validators) == 0 {
		return genesisValidators
	}else {
		for i := range validators {
			genesisValidator.OperatorAddress = sdk.ToAccAddress(validators[i].PubKey.Address())
			genesisValidator.Address = validators[i].PubKey.Address()
			pubByte, _ := cdc.MarshalJSON(validators[i].PubKey)
			pubstr := hex.EncodeToString(pubByte)
			genesisValidator.ConsensusKey = pubstr
			genesisValidator.Status = 1
			//genesisValidator.DelegatorShares = sdk.NewDec(100)
			//genesisValidator.Tokens = sdk.NewInt(100)
			genesisValidators = append(genesisValidators, genesisValidator)
		}
	}
	return genesisValidators
}

// DefaultGenesisState gets the raw genesis raw message for testing
func DefaultGenesisState(validators []tmtypes.GenesisValidator) GenesisState {
	genesisValidators := DefaultValidators(validators)
	//return GenesisState{Params:DefaultParams(), Validators:genesisValidators}
	return NewGenesisState(DefaultParams(), genesisValidators, nil)
}
