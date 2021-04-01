package types

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