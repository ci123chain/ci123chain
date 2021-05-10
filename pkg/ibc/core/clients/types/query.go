package types

import (
	codectypes "github.com/ci123chain/ci123chain/pkg/abci/codec/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
	"sort"
)

var (
	_ codectypes.UnpackInterfacesMessage = QueryClientStateResponse{}
	_ codectypes.UnpackInterfacesMessage = QueryClientStatesResponse{}
	_ codectypes.UnpackInterfacesMessage = QueryConsensusStateResponse{}
	_ codectypes.UnpackInterfacesMessage = QueryConsensusStatesResponse{}
)

// IdentifiedClientStates defines a slice of ClientConsensusStates that supports the sort interface
type IdentifiedClientStates []IdentifiedClientState

// Len implements sort.Interface
func (ics IdentifiedClientStates) Len() int { return len(ics) }

// Less implements sort.Interface
func (ics IdentifiedClientStates) Less(i, j int) bool { return ics[i].ClientId < ics[j].ClientId }

// Swap implements sort.Interface
func (ics IdentifiedClientStates) Swap(i, j int) { ics[i], ics[j] = ics[j], ics[i] }

// Sort is a helper function to sort the set of IdentifiedClientStates in place
func (ics IdentifiedClientStates) Sort() IdentifiedClientStates {
	sort.Sort(ics)
	return ics
}

// NewQueryClientStateResponse creates a new QueryClientStateResponse instance.
func NewQueryClientStateResponse(
	clientState *codectypes.Any, proof []byte, height Height,
) *QueryClientStateResponse {
	return &QueryClientStateResponse{
		ClientState: clientState,
		Proof:       proof,
		ProofHeight: height,
	}
}

// UnpackInterfaces implements UnpackInterfacesMesssage.UnpackInterfaces
func (qcsr QueryClientStateResponse) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	return unpacker.UnpackAny(qcsr.ClientState, new(exported.ClientState))
}

// NewQueryConsensusStateResponse creates a new QueryConsensusStateResponse instance.
func NewQueryConsensusStateResponse(
	consensusState  *codectypes.Any, proof []byte, height Height,
) *QueryConsensusStateResponse {
	return &QueryConsensusStateResponse{
		ConsensusState: consensusState,
		Proof:          proof,
		ProofHeight:    height,
	}
}

// UnpackInterfaces implements UnpackInterfacesMesssage.UnpackInterfaces
func (qcsr QueryConsensusStateResponse) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	return unpacker.UnpackAny(qcsr.ConsensusState, new(exported.ConsensusState))
}

// UnpackInterfaces implements UnpackInterfacesMesssage.UnpackInterfaces
func (qcsr QueryClientStatesResponse) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	//for _, cs := range qcsr.ClientStates {
	//	if err := cs.UnpackInterfaces(unpacker); err != nil {
	//		return err
	//	}
	//}
	return nil
}

// UnpackInterfaces implements UnpackInterfacesMesssage.UnpackInterfaces
func (qcsr QueryConsensusStatesResponse) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	//for _, cs := range qcsr.ConsensusStates {
	//	if err := cs.UnpackInterfaces(unpacker); err != nil {
	//		return err
	//	}
	//}
	return nil
}
