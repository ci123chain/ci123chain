package types

import (
	codectypes "github.com/ci123chain/ci123chain/pkg/abci/codec/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
)

var CommitmentInterfaceRegistry codectypes.InterfaceRegistry

func init() {
	CommitmentInterfaceRegistry = codectypes.NewInterfaceRegistry()

	RegisterInterfaces(CommitmentInterfaceRegistry)
}

// RegisterInterfaces registers the commitment interfaces to protobuf Any.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*exported.Proof)(nil),
		&MerkleProof{},
	)
}