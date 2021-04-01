package types

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
)

var (
	ModuleName = "capability"
	// KeyIndex defines the key that stores the current globally unique capability
	// index.
	KeyIndex = []byte("index")

	// KeyPrefixIndexCapability defines a key prefix that stores index to capability
	// name mappings.
	KeyPrefixIndexCapability = []byte("capability_index")
)

// IndexToKey returns bytes to be used as a key for a given capability index.
func IndexToKey(index uint64) []byte {
	return sdk.Uint64ToBigEndian(index)
}

// IndexFromKey returns an index from a call to IndexToKey for a given capability
// index.
func IndexFromKey(key []byte) uint64 {
	return sdk.BigEndianToUint64(key)
}


// RevCapabilityKey returns a reverse lookup key for a given module and capability
// name.
func RevCapabilityKey(module, name string) []byte {
	return []byte(fmt.Sprintf("%s/rev/%s", module, name))
}

// FwdCapabilityKey returns a forward lookup key for a given module and capability
// reference.
func FwdCapabilityKey(module string, cap *Capability) []byte {
	return []byte(fmt.Sprintf("%s/fwd/%p", module, cap))
}