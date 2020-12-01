package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
)


const (
	QueryValidatorOutstandingRewards = "validator_outstanding_rewards"
	QueryCommunityPool  =  "community_pool"
	QueryWithdrawAddress  =  "withdraw_address"
	QueryValidatorCommission         = "validator_commission"
	QueryDelegationRewards           = "delegation_rewards"
	QueryAccountInfo  =  "account_info"
)


// params for query 'custom/distr/validator_outstanding_rewards'
type QueryValidatorOutstandingRewardsParams struct {
	ValidatorAddress sdk.AccAddress `json:"validator_address" yaml:"validator_address"`
}


// creates a new instance of QueryValidatorOutstandingRewardsParams
func NewQueryValidatorOutstandingRewardsParams(validatorAddr sdk.AccAddress) QueryValidatorOutstandingRewardsParams {
	return QueryValidatorOutstandingRewardsParams{
		ValidatorAddress: validatorAddr,
	}
}

// params for query 'custom/distr/withdraw_addr'
type QueryDelegatorWithdrawAddrParams struct {
	DelegatorAddress sdk.AccAddress `json:"delegator_address" yaml:"delegator_address"`
}

// NewQueryDelegatorWithdrawAddrParams creates a new instance of QueryDelegatorWithdrawAddrParams.
func NewQueryDelegatorWithdrawAddrParams(delegatorAddr sdk.AccAddress) QueryDelegatorWithdrawAddrParams {
	return QueryDelegatorWithdrawAddrParams{DelegatorAddress: delegatorAddr}
}


// params for query 'custom/distr/delegator_total_rewards' and 'custom/distr/delegator_validators'
type QueryDelegatorParams struct {
	DelegatorAddress sdk.AccAddress `json:"delegator_address" yaml:"delegator_address"`
}

// creates a new instance of QueryDelegationRewardsParams
func NewQueryDelegatorParams(delegatorAddr sdk.AccAddress) QueryDelegatorParams {
	return QueryDelegatorParams{
		DelegatorAddress: delegatorAddr,
	}
}

// params for query 'custom/distr/validator_commission'
type QueryValidatorCommissionParams struct {
	ValidatorAddress sdk.AccAddress `json:"validator_address" yaml:"validator_address"`
}

// creates a new instance of QueryValidatorCommissionParams
func NewQueryValidatorCommissionParams(validatorAddr sdk.AccAddress) QueryValidatorCommissionParams {
	return QueryValidatorCommissionParams{
		ValidatorAddress: validatorAddr,
	}
}

// params for query 'custom/distr/delegation_rewards'
type QueryDelegationRewardsParams struct {
	DelegatorAddress sdk.AccAddress `json:"delegator_address" yaml:"delegator_address"`
	ValidatorAddress sdk.AccAddress `json:"validator_address" yaml:"validator_address"`
}

// creates a new instance of QueryDelegationRewardsParams
func NewQueryDelegationRewardsParams(delegatorAddr sdk.AccAddress, validatorAddr sdk.AccAddress) QueryDelegationRewardsParams {
	return QueryDelegationRewardsParams{
		DelegatorAddress: delegatorAddr,
		ValidatorAddress: validatorAddr,
	}
}

type QueryDelegatorBalanceParams struct {
	AccountAddress   sdk.AccAddress    `json:"account_address" yaml:"account_address"`
}

func NewQueryDelegatorBalanceParams(accountAddr sdk.AccAddress) QueryDelegatorBalanceParams {
	return QueryDelegatorBalanceParams{AccountAddress:accountAddr}
}