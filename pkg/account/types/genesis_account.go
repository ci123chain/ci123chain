package types

import (
	"bytes"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/account/exported"
)

// GenesisAccount is a struct for account initialization used exclusively during genesis
type GenesisAccount struct {
	Address       	types.AccAddress `json:"address" yaml:"address"`
	Coins        	types.Coins         `json:"coins" yaml:"coins"`
	Sequence      	uint64         `json:"sequence_number" yaml:"sequence_number"`
	AccountNumber 	uint64         `json:"account_number" yaml:"account_number"`
}

// GenesisAccounts defines a set of genesis account
type GenesisAccounts []GenesisAccount

func NewGenesisAccountRaw(address types.AccAddress, coin types.Coins) GenesisAccount {
	return GenesisAccount{
		Address: address,
		Coins:    coin,
		Sequence:0,
		AccountNumber:0,
	}
}

func (ga GenesisAccount) Validate() error {
	return nil
}

func (ga GenesisAccount) ToAccount() exported.Account {
	bacc := NewBaseAccount(ga.Address, ga.Coins, nil, ga.AccountNumber, ga.Sequence)
	return bacc
}

func (gaccs GenesisAccounts) Contains(acc types.AccAddress) bool {
	for _, gacc := range gaccs {
		if bytes.Equal(gacc.Address.Bytes(), acc.Bytes()) {
			return true
		}
	}
	return false
}