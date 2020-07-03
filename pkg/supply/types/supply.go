package types

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/supply/exported"
	yaml "gopkg.in/yaml.v2"
)


// Implements Delegation interface
var _ exported.SupplyI = Supply{}


type Supply struct {
	Total  sdk.Coin    `json:"total"`   // total supply of tokens registered on the chain
}


// SetTotal sets the total supply.
func (supply Supply) SetTotal(total sdk.Coin) exported.SupplyI {
	supply.Total = total
	return supply
}

// GetTotal returns the supply total.
func (supply Supply) GetTotal() sdk.Coin {
	return supply.Total
}

// NewSupply creates a new Supply instance
func NewSupply(total sdk.Coin) exported.SupplyI {
	return Supply{total}
}

// DefaultSupply creates an empty Supply
func DefaultSupply() exported.SupplyI {
	return NewSupply(sdk.NewCoin(sdk.NewInt(100000000)))
}

// Inflate adds coins to the total supply
func (supply Supply) Inflate(amount sdk.Coin) exported.SupplyI {
	supply.Total = supply.Total.Add(amount)
	return supply
}

// Deflate subtracts coins from the total supply
func (supply Supply) Deflate(amount sdk.Coin) exported.SupplyI {
	supply.Total = supply.Total.Sub(amount)
	return supply
}

// String returns a human readable string representation of a supplier.
func (supply Supply) String() string {
	b, err := yaml.Marshal(supply)
	if err != nil {
		panic(err)
	}
	return string(b)
}

// ValidateBasic validates the Supply coins and returns error if invalid
func (supply Supply) ValidateBasic() error {
	if !supply.Total.IsValid() {
		return fmt.Errorf("invalid total supply: %s", supply.Total.String())
	}
	return nil
}