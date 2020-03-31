package types

import (
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
)

type ContractInfoParams struct {
	ContractAddress   sdk.AccAddress `json:"contract_address"`
}

func NewQueryContractInfoParams(contractAddress sdk.AccAddress) ContractInfoParams {

	params := ContractInfoParams{ContractAddress:contractAddress}
	return params
}

type CodeInfoParams struct {
	ID   uint64   `json:"id"`
}

func NewQueryCodeInfoParams(id uint64) CodeInfoParams {

	params := CodeInfoParams{ID:id}
	return params
}


type ContractState struct {
	Result   string    `json:"result"`
}