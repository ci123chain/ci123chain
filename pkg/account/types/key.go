package types

import "github.com/ci123chain/ci123chain/pkg/abci/types"

const ModuleName  = "accounts"
const RouteKey  = "accounts"
const QueryAccount  = "queryAccount"
var (
	// AddressStoreKeyPrefix prefix for account-by-address store
	AddressStoreKeyPrefix = []byte{0x01}

	GlobalAccountNumberKey = []byte("globalAccountNumber")

	BalancesPrefix = []byte("balances")
)

// AddressStoreKey turn an address to types used to get it from the account store
func AddressStoreKey(addr types.AccAddress) []byte {
	return append(AddressStoreKeyPrefix, addr.Bytes()...)
}
