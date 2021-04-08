package types

import (
	"github.com/ci123chain/ci123chain/pkg/abci/types/pagination"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
	"sort"
)

// QueryClientStateRequest is the request types for the Query/ClientState RPC
// method
type QueryClientStateRequest struct {
	// client state unique identifier
	ClientId string `json:"client_id,omitempty"`
}

type QueryClientStateResponse struct {
	// client state associated with the request identifier
	ClientState exported.ClientState `json:"client_state,omitempty"`
	// merkle proof of existence
	Proof []byte `json:"proof,omitempty"`
	// height at which the proof was retrieved
	ProofHeight Height `json:"proof_height"`
}


// QueryClientStatesRequest is the request type for the Query/ClientStates RPC
// method
type QueryClientStatesRequest struct {
	// pagination request
	Pagination *pagination.PageRequest `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

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

// QueryClientStatesResponse is the response type for the Query/ClientStates RPC
// method.
type QueryClientStatesResponse struct {
	// list of stored ClientStates of the chain.
	ClientStates IdentifiedClientStates `protobuf:"bytes,1,rep,name=client_states,json=clientStates,proto3,castrepeated=IdentifiedClientStates" json:"client_states"`
	// pagination response
	Pagination *pagination.PageResponse `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}


type QueryConsensusStateRequest struct {
	// client identifier
	ClientId string `protobuf:"bytes,1,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
	// consensus state revision number
	RevisionNumber uint64 `protobuf:"varint,2,opt,name=revision_number,json=revisionNumber,proto3" json:"revision_number,omitempty"`
	// consensus state revision height
	RevisionHeight uint64 `protobuf:"varint,3,opt,name=revision_height,json=revisionHeight,proto3" json:"revision_height,omitempty"`
	// latest_height overrrides the height field and queries the latest stored
	// ConsensusState
	LatestHeight bool `protobuf:"varint,4,opt,name=latest_height,json=latestHeight,proto3" json:"latest_height,omitempty"`
}

// QueryConsensusStateResponse is the response type for the Query/ConsensusState
// RPC method
type QueryConsensusStateResponse struct {
	// consensus state associated with the client identifier at the given height
	ConsensusState exported.ConsensusState `protobuf:"bytes,1,opt,name=consensus_state,json=consensusState,proto3" json:"consensus_state,omitempty"`
	// merkle proof of existence
	Proof []byte `protobuf:"bytes,2,opt,name=proof,proto3" json:"proof,omitempty"`
	// height at which the proof was retrieved
	ProofHeight Height `protobuf:"bytes,3,opt,name=proof_height,json=proofHeight,proto3" json:"proof_height"`
}

// NewQueryClientStateResponse creates a new QueryClientStateResponse instance.
func NewQueryClientStateResponse(
	clientState exported.ClientState, proof []byte, height Height,
) *QueryClientStateResponse {
	return &QueryClientStateResponse{
		ClientState: clientState,
		Proof:       proof,
		ProofHeight: height,
	}
}

// NewQueryConsensusStateResponse creates a new QueryConsensusStateResponse instance.
func NewQueryConsensusStateResponse(
	consensusState  exported.ConsensusState, proof []byte, height Height,
) *QueryConsensusStateResponse {
	return &QueryConsensusStateResponse{
		ConsensusState: consensusState,
		Proof:          proof,
		ProofHeight:    height,
	}
}