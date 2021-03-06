package evmtypes

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	ethcmn "github.com/ethereum/go-ethereum/common"
)

const (
	// ModuleName string name of module
	ModuleName = "evm"
)

// KVStore key prefixes
var (
	KeyPrefixBlockHash   = []byte{0x01}
	KeyPrefixBloom       = []byte{0x02}
	KeyPrefixLogs        = []byte{0x03}
	KeyPrefixCode        = []byte{0x04}
	KeyPrefixStorage     = []byte{0x05}
	KeyPrefixChainConfig = []byte{0x06}
)

// BloomKey defines the store key for a block Bloom
func BloomKey(height int64) []byte {
	return sdk.Uint64ToBigEndian(uint64(height))
}

// AddressStoragePrefix returns a prefix to iterate over a given account storage.
func AddressStoragePrefix(address ethcmn.Address) []byte {
	return append(KeyPrefixStorage, address.Bytes()...)
}
