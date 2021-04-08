package utils

import (
	"github.com/ci123chain/ci123chain/pkg/abci/types/pagination"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/ibc"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/channel/types"
	ibcclient "github.com/ci123chain/ci123chain/pkg/ibc/core/client"
	clienttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/clients/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/host"
	coretypes "github.com/ci123chain/ci123chain/pkg/ibc/core/types"
	sdkerrors "github.com/pkg/errors"
)

// QueryChannel returns a channel end.
// If prove is true, it performs an ABCI store query in order to retrieve the merkle proof. Otherwise,
// it uses the gRPC query client.
func QueryChannel(
	clientCtx context.Context, portID, channelID string, prove bool,
) (*types.QueryChannelResponse, error) {
	return queryChannelABCI(clientCtx, portID, channelID, prove)
}

func queryChannelABCI(clientCtx context.Context, portID, channelID string, prove bool) (*types.QueryChannelResponse, error) {
	key := host.ChannelKey(portID, channelID)

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

	// check if channel exists
	if len(value) == 0 {
		return nil, sdkerrors.Wrapf(types.ErrChannelNotFound, "portID (%s), channelID (%s)", portID, channelID)
	}

	var channel types.Channel
	if err := clientCtx.Cdc.UnmarshalBinaryBare(value, &channel); err != nil {
		return nil, err
	}

	return types.NewQueryChannelResponse(channel, proofBz, proofHeight), nil
}


func QueryChannelsABCI(clientCtx context.Context, offset, limit uint64,) (*types.QueryChannelsResponse, error) {
	path := "/custom/" + ibc.ModuleName + "/" + coretypes.QueryChannels
	req := &types.QueryChannelsRequest{
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
		return nil, sdkerrors.Wrap(types.ErrChannelNotFound, "channels not found")
	}

	clientStatesResp := types.MustUnmarshalQueryChannelsResp(types.ChannelCdc, value)
	return &clientStatesResp, nil
}
