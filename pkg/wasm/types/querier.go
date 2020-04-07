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

type ContractStateParam struct {
	ContractAddress  sdk.AccAddress  `json:"contract_address"`
	QueryMessage     []byte          `json:"query_message"`
}

func NewContractStateParam(addr sdk.AccAddress, msg []byte) ContractStateParam {
	param := ContractStateParam{
		ContractAddress: addr,
		QueryMessage:    msg,
	}
	return param
}

type ContractState struct {
	Result   string    `json:"result"`
}

type ContractListParams struct {
	AccountAddress       sdk.AccAddress   `json:"account_address"`
}

func NewContractListParams(accountAddress sdk.AccAddress) ContractListParams {
	params := ContractListParams{
		AccountAddress:  accountAddress,
	}
	return params
}

type ContractListResponse struct {
	ContractAddressList  []string   `json:"contract_address_list"`
}

func NewContractListResponse(contractList []string) ContractListResponse {
	return ContractListResponse{ContractAddressList:contractList}
}