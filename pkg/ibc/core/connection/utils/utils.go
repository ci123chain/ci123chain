package utils

import (
	"github.com/ci123chain/ci123chain/pkg/client/context"
	clienttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/clients/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/connection/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/host"
	ibcclient "github.com/ci123chain/ci123chain/pkg/ibc/core/client"
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
