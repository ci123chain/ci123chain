package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
)

type QueryBondsParams struct {
	DelegatorAddr sdk.AccAddress    `json:"delegator_addr"`
	ValidatorAddr sdk.AccAddress     `json:"validator_addr"`
}

func NewQueryBondsParams(delegatorAddr sdk.AccAddress, validatorAddr sdk.AccAddress) QueryBondsParams {
	return QueryBondsParams{
		DelegatorAddr: delegatorAddr,
		ValidatorAddr: validatorAddr,
	}
}

// QueryValidatorsParams defines the params for the following queries:
// - 'custom/staking/validators'
type QueryValidatorsParams struct {
	Page, Limit int
	Status      string
}

func NewQueryValidatorsParams(page, limit int, status string) QueryValidatorsParams {
	return QueryValidatorsParams{page, limit, status}
}

type QueryValidatorParams struct {
	ValidatorAddr sdk.AccAddress  `json:"validator_addr"`
}

func NewQueryValidatorParams(validatorAddr sdk.AccAddress) QueryValidatorParams {
	return QueryValidatorParams{
		ValidatorAddr: validatorAddr,
	}
}

type QueryDelegatorParams struct {
	DelegatorAddr sdk.AccAddress
}

func NewQueryDelegatorParams(delegatorAddr sdk.AccAddress) QueryDelegatorParams {
	return QueryDelegatorParams{
		DelegatorAddr: delegatorAddr,
	}
}

// defines the params for the following queries:
// - 'custom/staking/redelegation'
type QueryRedelegationParams struct {
	DelegatorAddr    sdk.AccAddress
	SrcValidatorAddr sdk.AccAddress
	DstValidatorAddr sdk.AccAddress
}

func NewQueryRedelegationParams(delegatorAddr sdk.AccAddress,
	srcValidatorAddr, dstValidatorAddr sdk.AccAddress) QueryRedelegationParams {

	return QueryRedelegationParams{
		DelegatorAddr:    delegatorAddr,
		SrcValidatorAddr: srcValidatorAddr,
		DstValidatorAddr: dstValidatorAddr,
	}
}

type QueryOperatorAddressesParams struct {
	ConsAddresses   []sdk.AccAddr     `json:"cons_addresses"`
}

func NewQueryOperatorAddressParams(operatorAddress... sdk.AccAddr) QueryOperatorAddressesParams{
	var addresses = make([]sdk.AccAddr, 0)
	for _, v := range operatorAddress{
		addresses = append(addresses, v)
	}
	return QueryOperatorAddressesParams{ConsAddresses:addresses}
}

type ValidatorOperatorAddressResponse struct {
	ConsAddress        string    `json:"cons_address"`
	OperatorAddress    string    `json:"operator_address"`
	Found              bool     `json:"found"`
}

func NewValidatorOperatorAddressResponse(consAddr, addr string, found bool) ValidatorOperatorAddressResponse {
	return ValidatorOperatorAddressResponse{ConsAddress: consAddr, OperatorAddress:addr, Found: found}
}