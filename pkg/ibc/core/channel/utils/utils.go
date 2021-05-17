package utils

import (
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
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
	cdc := codec.NewProtoCodec(clientCtx.InterfaceRegistry)
	if err := cdc.UnmarshalBinaryBare(value, &channel); err != nil {
		return nil, err
	}

	return types.NewQueryChannelResponse(channel, proofBz, proofHeight), nil
}




// QueryPacketCommitment returns a packet commitment.
// If prove is true, it performs an ABCI store query in order to retrieve the merkle proof. Otherwise,
// it uses the gRPC query client.
func QueryPacketCommitment(
	clientCtx context.Context, portID, channelID string,
	sequence uint64, prove bool,
) (*types.QueryPacketCommitmentResponse, error) {
	key := host.PacketCommitmentKey(portID, channelID, sequence)

	value, proofBz, proofHeight, err := ibcclient.QueryTendermintProof(clientCtx, key)
	if err != nil {
		return nil, err
	}

	// check if packet commitment exists
	if len(value) == 0 {
		return nil, sdkerrors.Wrapf(types.ErrPacketCommitmentNotFound, "portID (%s), channelID (%s), sequence (%d)", portID, channelID, sequence)
	}

	return types.NewQueryPacketCommitmentResponse(value, proofBz, proofHeight), nil
}




// QueryPacketAcknowledgement returns the data about a packet acknowledgement.
// If prove is true, it performs an ABCI store query in order to retrieve the merkle proof. Otherwise,
// it uses the gRPC query client
func QueryPacketAcknowledgement(clientCtx context.Context, portID, channelID string, sequence uint64, prove bool) (*types.QueryPacketAcknowledgementResponse, error) {
	key := host.PacketAcknowledgementKey(portID, channelID, sequence)

	value, proofBz, proofHeight, err := ibcclient.QueryTendermintProof(clientCtx, key)
	if err != nil {
		return nil, err
	}

	if len(value) == 0 {
		return nil, sdkerrors.Wrapf(types.ErrInvalidAcknowledgement, "portID (%s), channelID (%s), sequence (%d)", portID, channelID, sequence)
	}

	return types.NewQueryPacketAcknowledgementResponse(value, proofBz, proofHeight), nil
}



// QueryPacketReceipt returns data about a packet receipt.
// If prove is true, it performs an ABCI store query in order to retrieve the merkle proof. Otherwise,
// it uses the gRPC query client.
func QueryPacketReceipt(
	clientCtx context.Context, portID, channelID string,
	sequence uint64, prove bool,
) (*types.QueryPacketReceiptResponse, error) {
	key := host.PacketReceiptKey(portID, channelID, sequence)

	value, proofBz, proofHeight, err := ibcclient.QueryTendermintProof(clientCtx, key)
	if err != nil {
		return nil, err
	}

	return types.NewQueryPacketReceiptResponse(value != nil, proofBz, proofHeight), nil
}

// instead grpc request


func QueryChannels(clientCtx context.Context, offset, limit uint64,) (*types.QueryChannelsResponse, error) {
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


// QueryPacketAcknowledgement returns the data about a packet acknowledgement.
// If prove is true, it performs an ABCI store query in order to retrieve the merkle proof. Otherwise,
// it uses the gRPC query client
func QueryPacketCommitments(clientCtx context.Context, portID, channelID string, offset, limit uint64) (*types.QueryPacketCommitmentsResponse, error) {
	path := "/custom/" + ibc.ModuleName + "/" + coretypes.QueryPacketCommitments
	req := &types.QueryPacketCommitmentsRequest{
		PortId: portID,
		ChannelId: channelID,
		Pagination: &pagination.PageRequest{
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

	var resp types.QueryPacketCommitmentsResponse
	types.ChannelCdc.MustUnmarshalJSON(value, &resp)
	return &resp, nil
}

func QueryUnreceivedPackets(clientCtx context.Context, portID, channelID string, seqs []uint64) (*types.QueryUnreceivedPacketsResponse, error) {
	path := "/custom/" + ibc.ModuleName + "/" + coretypes.UnreceivedPacketd
	req := &types.QueryUnreceivedPacketsRequest{
		PortId: portID,
		ChannelId: channelID,
		PacketCommitmentSequences: seqs,
	}
	key := clientCtx.Cdc.MustMarshalJSON(req)
	value, _, err := ibcclient.QueryABCI(clientCtx, path, key, false)
	if err != nil {
		return nil, err
	}

	var resp types.QueryUnreceivedPacketsResponse
	types.ChannelCdc.MustUnmarshalJSON(value, &resp)
	return &resp, nil
}
