package types

import (
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
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

	ModePerm os.FileMode = 0666
	SuffixName = ".wasm"
	FolderName = "wasm"

	InitFunctionName = "init"
	HandleFunctionName = "handle"
	QueryFunctionName = "query"
)

var (
	WasmerKey            =  []byte("wasmer")

	CodeKeyPrefix        =  []byte{0x01}
	ContractKeyPrefix    =  []byte{0x02}
	ContractStorePrefix  =  []byte{0x03}
)

func GetWasmerKey() []byte {
	return WasmerKey
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