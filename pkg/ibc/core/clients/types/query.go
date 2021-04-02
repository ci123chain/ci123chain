package types

import (
	"github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
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