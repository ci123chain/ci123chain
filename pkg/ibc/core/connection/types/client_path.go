package types

// ClientPaths define all the connection paths for a client state.
type ClientPaths struct {
	// list of connection paths
	Paths []string `protobuf:"bytes,1,rep,name=paths,proto3" json:"paths,omitempty"`
}
