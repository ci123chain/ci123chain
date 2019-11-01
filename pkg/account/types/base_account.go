package types

import (
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tendermint/tendermint/crypto"
)

type BaseAccount struct {
	Address       	types.AccAddress `json:"address" yaml:"address"`
	Coin        	types.Coin         `json:"coin" yaml:"coin"`
	Sequence      	uint64         `json:"sequence_number" yaml:"sequence_number"`
	AccountNumber 	uint64         `json:"account_number" yaml:"account_number"`
	PubKey 			crypto.PubKey  `json:"pub_key" yaml:"pub_key"`
}

// NewBaseAccount creates a new BaseAccount object
func NewBaseAccount(address types.AccAddress, coin types.Coin,
	pubKey crypto.PubKey, accountNumber uint64, sequence uint64) *BaseAccount {

	return &BaseAccount{
		Address:       address,
		Coin:          coin,
		PubKey:        pubKey,
		AccountNumber: accountNumber,
		Sequence:      sequence,
	}
}

// GetAddress - Implements sdk.Account.
func (acc BaseAccount) GetAddress() types.AccAddress {
	return acc.Address
}

// SetAddress - Implements sdk.Account.
func (acc *BaseAccount) SetAddress(addr types.AccAddress) error {
	if err := addr.Validate(); err != nil {
		return err
	}
	acc.Address = addr
	return nil
}


// GetPubKey - Implements sdk.Account.
func (acc BaseAccount) GetPubKey() crypto.PubKey {
	return acc.PubKey
}

// SetPubKey - Implements sdk.Account.
func (acc *BaseAccount) SetPubKey(pubKey crypto.PubKey) error {
	acc.PubKey = pubKey
	return nil
}

// GetCoins - Implements sdk.Account.
func (acc *BaseAccount) GetCoin() types.Coin {
	return acc.Coin
}

// SetCoins - Implements sdk.Account.
func (acc *BaseAccount) SetCoin(coin types.Coin) error {
	acc.Coin = coin
	return nil
}

// GetAccountNumber - Implements Account
func (acc *BaseAccount) GetAccountNumber() uint64 {
	return acc.AccountNumber
}

// SetAccountNumber - Implements Account
func (acc *BaseAccount) SetAccountNumber(accNumber uint64) error {
	acc.AccountNumber = accNumber
	return nil
}

// GetSequence - Implements sdk.Account.
func (acc *BaseAccount) GetSequence() uint64 {
	return acc.Sequence
}

// SetSequence - Implements sdk.Account.
func (acc *BaseAccount) SetSequence(seq uint64) error {
	acc.Sequence = seq
	return nil
}

