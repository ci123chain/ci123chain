package types

import (
	"bytes"
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
)

// GenesisAccount is a struct for account initialization used exclusively during genesis
type GenesisAccount struct {
	Address       	types.AccAddress `json:"address" yaml:"address"`
	Coin        	uint64         `json:"coin" yaml:"coin"`
	Sequence      	uint64         `json:"sequence_number" yaml:"sequence_number"`
	AccountNumber 	uint64         `json:"account_number" yaml:"account_number"`
}

// GenesisAccounts defines a set of genesis account
type GenesisAccounts []GenesisAccount

func NewGenesisAccountRaw(address types.AccAddress, coin uint64) GenesisAccount {
	return GenesisAccount{
		Address: address,
		Coin:    coin,
		Sequence:0,
		AccountNumber:0,
	}
}

func (ga GenesisAccount) Validate() error {
	return nil
}

func (gaccs GenesisAccounts) Contains(acc types.AccAddress) bool {
	for _, gacc := range gaccs {
		if bytes.Equal(gacc.Address.Bytes(), acc.Bytes()) {
			return true
		}
	}
	return false
}