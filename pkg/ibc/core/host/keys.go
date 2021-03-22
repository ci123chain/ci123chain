package host

import (
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
)

const (
	ModuleName = "ibc"

	RouterKey string = ModuleName
)

// KVStore key prefixes for IBC
var (
	KeyClientStorePrefix = []byte("clients")
)

const (
	KeyClientState             = "clientState"
	KeyConsensusStatePrefix    = "consensusStates"
)
// ClientStateKey returns a store key under which a particular client state is stored
// in a client prefixed store
func ClientStateKey() []byte {
	return []byte(KeyClientState)
}

// ConsensusStatePath returns the suffix store key for the consensus state at a
// particular height stored in a client prefixed store.
func ConsensusStatePath(height exported.Height) string {
	return fmt.Sprintf("%s/%s", KeyConsensusStatePrefix, height)
}

// ConsensusStateKey returns the store key for a the consensus state of a particular
// client stored in a client prefixed store.
func ConsensusStateKey(height exported.Height) []byte {
	return []byte(ConsensusStatePath(height))
}

func FullClientStatePath(clientID string) string {
	return FullClientPath(clientID, KeyClientState)
}
// FullClientPath returns the full path of a specific client path in the format:
// "clients/{clientID}/{path}" as a string.
func FullClientPath(clientID string, path string) string {
	return fmt.Sprintf("%s/%s/%s", KeyClientStorePrefix, clientID, path)
}

func FullConsensusStatePath(clientID string, height exported.Height) string {
	return FullClientPath(clientID, ConsensusStatePath(height))
}