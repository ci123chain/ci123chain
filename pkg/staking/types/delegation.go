package types

import (
	"fmt"
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	yaml "gopkg.in/yaml.v2"
	"time"
)

type Delegation struct {
	DelegatorAddress  sdk.AccAddress
	ValidatorAddress  sdk.AccAddress
	Shares            sdk.Dec
}
// IsMature - is the current entry mature
func (e RedelegationEntry) IsMature(currentTime time.Time) bool {
	return !e.CompletionTime.After(currentTime)
}

// RemoveEntry - remove entry at index i to the unbonding delegation
func (red *Redelegation) RemoveEntry(i int64) {
	red.Entries = append(red.Entries[:i], red.Entries[i+1:]...)
}

func NewDelegation(delegatorAddr sdk.AccAddress, validatorAddr sdk.AccAddress, shares sdk.Dec) Delegation {
	return Delegation{
		DelegatorAddress:delegatorAddr,
		ValidatorAddress:validatorAddr,
		Shares:shares,
	}
}

// String returns a human readable string representation of a Delegation.
func (d Delegation) String() string {
	out, _ := yaml.Marshal(d)
	return string(out)
}
func (d Delegation) GetDelegatorAddr() sdk.AccAddress { return d.DelegatorAddress }
func (d Delegation) GetValidatorAddr() sdk.AccAddress { return d.ValidatorAddress }


type Delegations []Delegation

func (d Delegation) GetShares() sdk.Dec {return d.Shares}

type UnbondingDelegation struct {
	DelegatorAddress   sdk.AccAddress
	ValidatorAddress   sdk.AccAddress
	Entries            []UnbondingDelegationEntry
}

type UnbondingDelegationEntry struct {
	CreationHeight    int64
	CompletionTime    time.Time
	InitialBalance    sdk.Int
	Balance           sdk.Int
}

// IsMature - is the current entry mature
func (e UnbondingDelegationEntry) IsMature(currentTime time.Time) bool {
	return !e.CompletionTime.After(currentTime)
}

// RemoveEntry - remove entry at index i to the unbonding delegation
func (ubd *UnbondingDelegation) RemoveEntry(i int64) {
	ubd.Entries = append(ubd.Entries[:i], ubd.Entries[i+1:]...)
}

// NewUnbondingDelegation - create a new unbonding delegation object
func NewUnbondingDelegation(
	delegatorAddr sdk.AccAddress, validatorAddr sdk.AccAddress,
	creationHeight int64, minTime time.Time, balance sdk.Int,
) UnbondingDelegation {

	return UnbondingDelegation{
		DelegatorAddress: delegatorAddr,
		ValidatorAddress: validatorAddr,
		Entries: []UnbondingDelegationEntry{
			NewUnbondingDelegationEntry(creationHeight, minTime, balance),
		},
	}
}

func NewUnbondingDelegationEntry(creationHeight int64, completionTime time.Time, balance sdk.Int) UnbondingDelegationEntry {
	return UnbondingDelegationEntry{
		CreationHeight: creationHeight,
		CompletionTime: completionTime,
		InitialBalance: balance,
		Balance:        balance,
	}
}

// AddEntry - append entry to the unbonding delegation
func (ubd *UnbondingDelegation) AddEntry(creationHeight int64, minTime time.Time, balance sdk.Int) {
	entry := NewUnbondingDelegationEntry(creationHeight, minTime, balance)
	ubd.Entries = append(ubd.Entries, entry)
}

type DelegationResponse struct {
	Delegation
	Balance sdk.Coin `json:"balance" yaml:"balance"`
}

// NewDelegationResp creates a new DelegationResponse instance
func NewDelegationResp(
	delegatorAddr sdk.AccAddress, validatorAddr sdk.AccAddress, shares sdk.Dec, balance sdk.Coin,
) DelegationResponse {
	return DelegationResponse{
		Delegation: NewDelegation(delegatorAddr, validatorAddr, shares),
		Balance:    balance,
	}
}

// String implements the Stringer interface for DelegationResponse.
func (d DelegationResponse) String() string {
	return fmt.Sprintf("%s\n  Balance:   %s", d.Delegation.String(), d.Balance)
}

// DelegationResponses is a collection of DelegationResp
type DelegationResponses []DelegationResponse
