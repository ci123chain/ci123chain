package types

import clienttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/clients/types"

// QueryServer is the server API for Query service.
//types QueryServer interface {
//	// Connection queries an IBC connection end.
//	Connection(sdk.Context, *QueryConnectionRequest) (*QueryConnectionResponse, error)
//	// Connections queries all the IBC connections of a chain.
//	Connections(sdk.Context, *QueryConnectionsRequest) (*QueryConnectionsResponse, error)
//	// ClientConnections queries the connection paths associated with a client
//	// state.
//	ClientConnections(sdk.Context, *QueryClientConnectionsRequest) (*QueryClientConnectionsResponse, error)
//	// ConnectionClientState queries the client state associated with the
//	// connection.
//	ConnectionClientState(sdk.Context, *QueryConnectionClientStateRequest) (*QueryConnectionClientStateResponse, error)
//	// ConnectionConsensusState queries the consensus state associated with the
//	// connection.
//	ConnectionConsensusState(sdk.Context, *QueryConnectionConsensusStateRequest) (*QueryConnectionConsensusStateResponse, error)
//}



// QueryConnectionResponse is the response type for the Query/Connection RPC
// method. Besides the connection end, it includes a proof and the height from
// which the proof was retrieved.
type QueryConnectionResponse struct {
	// connection associated with the request identifier
	Connection *ConnectionEnd `protobuf:"bytes,1,opt,name=connection,proto3" json:"connection,omitempty"`
	// merkle proof of existence
	Proof []byte `protobuf:"bytes,2,opt,name=proof,proto3" json:"proof,omitempty"`
	// height at which the proof was retrieved
	ProofHeight clienttypes.Height `protobuf:"bytes,3,opt,name=proof_height,json=proofHeight,proto3" json:"proof_height"`
}

// NewQueryConnectionResponse creates a new QueryConnectionResponse instance
func NewQueryConnectionResponse(
	connection ConnectionEnd, proof []byte, height clienttypes.Height,
) *QueryConnectionResponse {
	return &QueryConnectionResponse{
		Connection:  &connection,
		Proof:       proof,
		ProofHeight: height,
	}
}