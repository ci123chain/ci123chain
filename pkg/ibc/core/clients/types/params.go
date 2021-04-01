package types

import (
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
	paramtypes "github.com/ci123chain/ci123chain/pkg/params/subspace"
	"strings"
)


var (
	// DefaultAllowedClients are "06-solomachine" and "07-tendermint"
	DefaultAllowedClients = []string{exported.Solomachine, exported.Tendermint}

	// KeyAllowedClients is store's key for AllowedClients Params
	KeyAllowedClients = []byte("AllowedClients")
)


type Params struct {
	// allowed_clients defines the list of allowed client state types.
	AllowedClients []string `json:"allowed_clients,omitempty" yaml:"allowed_clients"`
}

// ParamKeyTable types declaration for parameters
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new parameter configuration for the ibc transfer module
func NewParams(allowedClients ...string) Params {
	return Params{
		AllowedClients: allowedClients,
	}
}

// DefaultParams is the default parameter configuration for the ibc-transfer module
func DefaultParams() Params {
	return NewParams(DefaultAllowedClients...)
}

// Validate all ibc-transfer module parameters
func (p Params) Validate() error {
	return validateClients(p.AllowedClients)
}

// ParamSetPairs implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyAllowedClients, p.AllowedClients, validateClients),
	}
}

func validateClients(i interface{}) error {
	clients, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter types: %T", i)
	}

	for i, clientType := range clients {
		if strings.TrimSpace(clientType) == "" {
			return fmt.Errorf("client types %d cannot be blank", i)
		}
	}

	return nil
}

func (p Params) IsAllowedClient(clientType string) bool {
	for _, allowedClient := range p.AllowedClients {
		if allowedClient == clientType {
			return true
		}
	}
	return false
}


