package types

import "github.com/ci123chain/ci123chain/pkg/abci/types"

const ModuleName  = "accounts"
const RouteKey  = "accounts"
const StoreKey  = "accounts"
const QueryAccount  = "queryAccount"
const QueryAccountNonce = "queryAccountNonce"
const QueryHistoryAccount = "historyAccount"
var (
	// AddressStoreKeyPrefix prefix for account-by-address store
	AddressStoreKeyPrefix = []byte{0x01}

	HeightUpdateKeyPrefix = []byte{0x02}

	HeightsUpdateKeyPrefix = []byte{0x03}

	GlobalAccountNumberKey = []byte("globalAccountNumber")

	BalancesPrefix = []byte("balances")

	EventTypeTransfer = "transfer"
	AttributeKeyRecipient = "recipient"
	AttributeKeySender    = "sender"
)

// AddressStoreKey turn an address to types used to get it from the account store
func AddressStoreKey(addr types.AccAddress) []byte {
	return append(AddressStoreKeyPrefix, addr.Bytes()...)
}


func HeightUpdateKey(addr types.AccAddress, height int64) []byte {
	a := append(HeightUpdateKeyPrefix, addr.Bytes()...)
	return append(a, []byte(string(height))...)
}

func HeightsUpdateKey(addr types.AccAddress) []byte {
	return append(HeightsUpdateKeyPrefix, addr.Bytes()...)
}