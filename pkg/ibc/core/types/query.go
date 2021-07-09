package types

import (
	"github.com/ci123chain/ci123chain/pkg/ibc/core/channel"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/connection"
	"github.com/gogo/protobuf/grpc"

	channeltypes "github.com/ci123chain/ci123chain/pkg/ibc/core/channel/types"
	client "github.com/ci123chain/ci123chain/pkg/ibc/core/clients"
	clienttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/clients/types"
	connectiontypes "github.com/ci123chain/ci123chain/pkg/ibc/core/connection/types"
)

// QueryServer defines the IBC interfaces that the gRPC query server must implement
type QueryServer interface {
	clienttypes.QueryServer
	connectiontypes.QueryServer
	channeltypes.QueryServer
}

// RegisterQueryService registers each individual IBC submodule query service
func RegisterQueryService(server grpc.Server, queryService QueryServer) {
	client.RegisterQueryService(server, queryService)
	connection.RegisterQueryService(server, queryService)
	channel.RegisterQueryService(server, queryService)
}
