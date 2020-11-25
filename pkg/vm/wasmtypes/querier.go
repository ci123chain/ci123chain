package types

import (
	"encoding/hex"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/vm/moduletypes/utils"
	"strings"
)

type ContractInfoParams struct {
	ContractAddress   sdk.AccAddress `json:"contract_address"`
}

func NewQueryContractInfoParams(contractAddress sdk.AccAddress) ContractInfoParams {
	params := ContractInfoParams{ContractAddress:contractAddress}
	return params
}

type CodeInfoParams struct {
	Hash   []byte   `json:"hash"`
}

func NewQueryCodeInfoParams(hashStr string) CodeInfoParams {
	hash, _ := hex.DecodeString(strings.ToLower(hashStr))
	params := CodeInfoParams{Hash: hash}
	return params
}

type ContractStateParam struct {
	ContractAddress  sdk.AccAddress  `json:"contract_address"`
	InvokerAddress	 sdk.AccAddress	 `json:"invoker_address"`
	QueryMessage     utils.CallData         `json:"query_message"`
}

func NewContractStateParam(addr, invokerAddr sdk.AccAddress,msg utils.CallData) ContractStateParam {
	param := ContractStateParam{
		ContractAddress: addr,
		InvokerAddress:  invokerAddr,
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

type ContractExistParams struct {
	WasmCodeHash       []byte  `json:"wasm_code_hash"`
}

func NewContractExistParams(hash []byte) ContractExistParams {
	param := ContractExistParams{
		WasmCodeHash: hash,
	}
	return param
}