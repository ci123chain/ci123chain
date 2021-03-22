package types

import (
	"github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
)

// QueryClientStateRequest is the request type for the Query/ClientState RPC
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
