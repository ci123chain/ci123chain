package channel


import (
	"github.com/ci123chain/ci123chain/pkg/ibc/core/channel/types"
	"github.com/gogo/protobuf/grpc"
)

// RegisterQueryService registers the gRPC query service for IBC client.
func RegisterQueryService(server grpc.Server, queryServer types.QueryServer) {
	types.RegisterQueryServer(server, queryServer)
}