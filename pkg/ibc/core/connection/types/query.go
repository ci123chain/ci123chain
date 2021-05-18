package types

import (
	codectypes "github.com/ci123chain/ci123chain/pkg/abci/codec/types"
	clienttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/clients/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
)

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
//type QueryConnectionResponse struct {
//	// connection associated with the request identifier
//	Connection *ConnectionEnd `protobuf:"bytes,1,opt,name=connection,proto3" json:"connection,omitempty"`
//	// merkle proof of existence
//	Proof []byte `protobuf:"bytes,2,opt,name=proof,proto3" json:"proof,omitempty"`
//	// height at which the proof was retrieved
//	ProofHeight clienttypes.Height `protobuf:"bytes,3,opt,name=proof_height,json=proofHeight,proto3" json:"proof_height"`
//}

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

//// QueryClientConnectionsRequest is the request type for the
//// Query/ClientConnections RPC method
//type QueryClientConnectionsRequest struct {
//	// client identifier associated with a connection
//	ClientId string `protobuf:"bytes,1,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
//}
//
//// QueryClientConnectionsResponse is the response type for the
//// Query/ClientConnections RPC method
//type QueryClientConnectionsResponse struct {
//	// slice of all the connection paths associated with a client.
//	ConnectionPaths []string `protobuf:"bytes,1,rep,name=connection_paths,json=connectionPaths,proto3" json:"connection_paths,omitempty"`
//	// merkle proof of existence
//	Proof []byte `protobuf:"bytes,2,opt,name=proof,proto3" json:"proof,omitempty"`
//	// height at which the proof was generated
//	ProofHeight clienttypes.Height `protobuf:"bytes,3,opt,name=proof_height,json=proofHeight,proto3" json:"proof_height"`
//}
// NewQueryClientConnectionsRequest creates a new QueryClientConnectionsRequest instance
func NewQueryClientConnectionsRequest(clientID string) *QueryClientConnectionsRequest {
	return &QueryClientConnectionsRequest{
		ClientId: clientID,
	}
}

// NewQueryClientConnectionsResponse creates a new ConnectionPaths instance
func NewQueryClientConnectionsResponse(
	connectionPaths []string, proof []byte, height clienttypes.Height,
) *QueryClientConnectionsResponse {
	return &QueryClientConnectionsResponse{
		ConnectionPaths: connectionPaths,
		Proof:           proof,
		ProofHeight:     height,
	}
}


//// QueryConnectionsRequest is the request type for the Query/Connections RPC
//// method
//type QueryConnectionsRequest struct {
//	Pagination *pagination.PageRequest `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
//}
//
//// QueryConnectionsResponse is the response type for the Query/Connections RPC
//// method.
//type QueryConnectionsResponse struct {
//	// list of stored connections of the chain.
//	Connections []*IdentifiedConnection `protobuf:"bytes,1,rep,name=connections,proto3" json:"connections,omitempty"`
//	// pagination response
//	Pagination *pagination.PageResponse `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
//	// query block height
//	Height clienttypes.Height `protobuf:"bytes,3,opt,name=height,proto3" json:"height"`
//}

// NewQueryConnectionClientStateResponse creates a newQueryConnectionClientStateResponse instance
func NewQueryConnectionClientStateResponse(identifiedClientState clienttypes.IdentifiedClientState, proof []byte, height clienttypes.Height) *QueryConnectionClientStateResponse {
	return &QueryConnectionClientStateResponse{
		IdentifiedClientState: &identifiedClientState,
		Proof:                 proof,
		ProofHeight:           height,
	}
}

// NewQueryConnectionConsensusStateResponse creates a newQueryConnectionConsensusStateResponse instance
func NewQueryConnectionConsensusStateResponse(clientID string, anyConsensusState *codectypes.Any, consensusStateHeight exported.Height, proof []byte, height clienttypes.Height) *QueryConnectionConsensusStateResponse {
	return &QueryConnectionConsensusStateResponse{
		ConsensusState: anyConsensusState,
		ClientId:       clientID,
		Proof:          proof,
		ProofHeight:    height,
	}
}