package types


// GenesisState defines the ibc connection submodule's genesis state.
type GenesisState struct {
	Connections           []IdentifiedConnection `protobuf:"bytes,1,rep,name=connections,proto3" json:"connections"`
	ClientConnectionPaths []ConnectionPaths      `protobuf:"bytes,2,rep,name=client_connection_paths,json=clientConnectionPaths,proto3" json:"client_connection_paths" yaml:"client_connection_paths"`
	// the sequence for the next generated connection identifier
	NextConnectionSequence uint64 `protobuf:"varint,3,opt,name=next_connection_sequence,json=nextConnectionSequence,proto3" json:"next_connection_sequence,omitempty" yaml:"next_connection_sequence"`
}

// ConnectionPaths define all the connection paths for a given client state.
type ConnectionPaths struct {
	// client state unique identifier
	ClientId string `protobuf:"bytes,1,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty" yaml:"client_id"`
	// list of connection paths
	Paths []string `protobuf:"bytes,2,rep,name=paths,proto3" json:"paths,omitempty"`
}


// DefaultGenesisState returns the ibc connection submodule's default genesis state.
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Connections:            []IdentifiedConnection{},
		ClientConnectionPaths:  []ConnectionPaths{},
		NextConnectionSequence: 0,
	}
}
