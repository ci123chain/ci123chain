package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"os"
)

const (
	RouteKey = "wasm"
	ModuleName = "wasm"
	StoreKey = "wasm"
	QueryContractInfo  = "queryContractInfo"
	QueryCodeInfo      = "queryCodeInfo"
	QueryContractState = "queryContractState"
	QueryContractList  = "queryContractList"
	QueryContractExist  = "queryContractExist"

	SystemContract = "system_contract"

	ModePerm os.FileMode = 0666
	SuffixName = ".wasm"
	FolderName = "wasm"

	InitFunctionName = "init"
	HandleFunctionName = "handle"
	QueryFunctionName = "query"
)

var (
	WasmerKey            		=  []byte("wasmer")

	CodeKeyPrefix        		=  []byte{0x01}
	ContractKeyPrefix    		=  []byte{0x02}
	ContractStorePrefix  		=  []byte{0x03}
	AccountContractListPrefix 	=  []byte{0x04}
)

func GetWasmerKey() []byte {
	return WasmerKey
}

func GetAccountContractListKey(accountAddr sdk.AccAddress) []byte {
	return append(AccountContractListPrefix, accountAddr.Bytes()...)
}

func GetCodeKey(codeHash []byte) []byte {
	return append(CodeKeyPrefix, codeHash...)
}

func GetContractAddressKey(addr sdk.AccAddress) []byte {
	return append(ContractKeyPrefix, addr.Bytes()...)
}

func GetContractStorePrefixKey(addr sdk.AccAddress) []byte {
	return append(ContractStorePrefix, addr.Bytes()...)
}