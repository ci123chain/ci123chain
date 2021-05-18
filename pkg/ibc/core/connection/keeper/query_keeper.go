package keeper

import (
	store2 "github.com/ci123chain/ci123chain/pkg/abci/store"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/abci/types/pagination"
	clienttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/clients/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/connection/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/host"
	abci "github.com/tendermint/tendermint/abci/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Connections implements the Query/Connections gRPC method
func (q Keeper) ConnectionsRest(ctx sdk.Context, r abci.RequestQuery) ([]byte, error) {

	req, err := types.UnmarshalConnectionRequest(types.IBCConnectionCodec, r.Data)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	connections := []*types.IdentifiedConnection{}
	store := store2.NewPrefixStore(ctx.KVStore(q.storeKey), []byte(host.KeyConnectionPrefix))

	pageRes, err := pagination.Paginate(store, req.Pagination, func(key, value []byte) error {
		var result types.ConnectionEnd
		if err := q.cdc.UnmarshalBinaryBare(value, &result); err != nil {
			return err
		}

		connectionID, err := host.ParseConnectionPath(string(key))
		if err != nil {
			return err
		}

		identifiedConnection := types.NewIdentifiedConnection(connectionID, result)
		connections = append(connections, &identifiedConnection)
		return nil
	})

	if err != nil {
		return nil, err
	}

	res := types.QueryConnectionsResponse{
		Connections: connections,
		Pagination:  pageRes,
		Height:      clienttypes.GetSelfHeight(ctx),
	}
	return types.MustMarshalConnectionResponse(types.IBCConnectionCodec, res), nil
}