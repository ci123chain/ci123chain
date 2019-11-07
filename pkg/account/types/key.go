package types

import "github.com/tanhuiya/ci123chain/pkg/abci/types"

const ModuleName  = "accounts"

var (
	// AddressStoreKeyPrefix prefix for account-by-address store
	AddressStoreKeyPrefix = []byte{0x01}

	GlobalAccountNumberKey = []byte("globalAccountNumber")
)

// AddressStoreKey turn an address to types used to get it from the account store
func AddressStoreKey(addr types.AccAddress) []byte {
	return append(AddressStoreKeyPrefix, addr.Bytes()...)
}
