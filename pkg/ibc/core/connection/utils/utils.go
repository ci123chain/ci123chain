package utils

import (
	"github.com/ci123chain/ci123chain/pkg/abci/types/pagination"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/ibc"
	clienttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/clients/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/connection/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/host"
	ibcclient "github.com/ci123chain/ci123chain/pkg/ibc/core/client"
	coretypes "github.com/ci123chain/ci123chain/pkg/ibc/core/types"
	sdkerrors "github.com/pkg/errors"
)

// QueryConnection returns a connection end.
// If prove is true, it performs an ABCI store query in order to retrieve the merkle proof. Otherwise,
// it uses the gRPC query client.
func QueryConnection(
	clientCtx context.Context, connectionID string, prove bool,
) (*types.QueryConnectionResponse, error) {
	return queryConnectionABCI(clientCtx, connectionID, prove)
}


func queryConnectionABCI(clientCtx context.Context, connectionID string, prove bool) (*types.QueryConnectionResponse, error) {
	key := host.ConnectionKey(connectionID)
	var value, proofBz []byte
	var err error
	var proofHeight clienttypes.Height
	if prove {
		value, proofBz, proofHeight, err = ibcclient.QueryTendermintProof(clientCtx, key)
	} else {
		value, proofHeight, err = ibcclient.QueryTendermint(clientCtx, key)
	}
	if err != nil {
		return nil, err
	}

	// check if connection exists
	if len(value) == 0 {
		return nil, sdkerrors.Wrap(types.ErrConnectionNotFound, connectionID)
	}

	var connection types.ConnectionEnd
	if err := clientCtx.Cdc.UnmarshalBinaryBare(value, &connection); err != nil {
		return nil, err
	}

	return types.NewQueryConnectionResponse(connection, proofBz, proofHeight), nil
}


// QueryClientStateABCI queries the store to get the light client state and a merkle proof.
func QueryConnectionsABCI(
	clientCtx context.Context, offset, limit uint64,
) (*types.QueryConnectionsResponse, error) {
	path := "/custom/" + ibc.ModuleName + "/" + coretypes.QueryConnections

	req := &types.QueryConnectionsRequest{
		Pagination: &pagination.PageRequest{
			Key:        []byte(""),
			Offset:     offset,
			Limit:      limit,
			CountTotal: true,
		},
	}
	key := clientCtx.Cdc.MustMarshalJSON(req)
	value, _, err := ibcclient.QueryABCI(clientCtx, path, key, false)
	if err != nil {
		return nil, err
	}

	// check if client exists
	if len(value) == 0 {
		return nil, sdkerrors.Wrap(types.ErrConnectionNotFound, "connections not found")
	}

	clientStatesResp := types.MustUnmarshalConnectionResponse(types.IBCConnectionCodec, value)
	return &clientStatesResp, nil
}
