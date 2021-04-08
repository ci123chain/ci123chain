package keeper

import (
	"github.com/ci123chain/ci123chain/pkg/abci/store"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/abci/types/pagination"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/channel/types"
	clienttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/clients/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/host"
	abci "github.com/tendermint/tendermint/abci/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
	"strings"
)

func (q Keeper) Channels(ctx sdk.Context, r abci.RequestQuery) ([]byte, error) {
	req, err := types.UnmarshalQueryChannelRequest(types.ChannelCdc, r.Data)
	if err != nil {
		return nil, err
	}
	channels := []*types.IdentifiedChannel{}
	channelStore := store.NewPrefixStore(ctx.KVStore(q.storeKey), []byte(host.KeyChannelEndPrefix))

	pageRes, err := pagination.Paginate(channelStore, req.Pagination, func(key, value []byte) error {
		var result types.Channel
		if err := q.cdc.UnmarshalBinaryBare(value, &result); err != nil {
			return err
		}

		portID, channelID, err := host.ParseChannelPath(string(key))
		if err != nil {
			return err
		}

		identifiedChannel := types.NewIdentifiedChannel(portID, channelID, result)
		channels = append(channels, &identifiedChannel)
		return nil
	})

	if err != nil {
		return nil, err
	}

	selfHeight := clienttypes.GetSelfHeight(ctx)

	resp := types.QueryChannelsResponse{
		Channels:   channels,
		Pagination: pageRes,
		Height:     selfHeight,
	}
	return types.MustMarshalQueryChannelsResp(types.ChannelCdc, resp), nil
}

// PacketCommitment implements the Query/PacketCommitment gRPC method
func (q Keeper) PacketCommitment(ctx sdk.Context, req abci.RequestQuery) ([]byte, error) {
	var reqPacketCommitment types.QueryPacketCommitmentRequest
	if err := types.ChannelCdc.UnmarshalJSON(req.Data, &reqPacketCommitment); err != nil {
		return nil, err
	}

	if err := validateRPCRequest(reqPacketCommitment.PortId, reqPacketCommitment.ChannelId); err != nil {
		return nil, err
	}

	if reqPacketCommitment.Sequence == 0 {
		return nil, status.Error(codes.InvalidArgument, "packet sequence cannot be 0")
	}

	commitmentBz := q.GetPacketCommitment(ctx, reqPacketCommitment.PortId, reqPacketCommitment.ChannelId, reqPacketCommitment.Sequence)
	if len(commitmentBz) == 0 {
		return nil, status.Error(codes.NotFound, "packet commitment hash not found")
	}

	selfHeight := clienttypes.GetSelfHeight(ctx)
    resp := types.NewQueryPacketCommitmentResponse(commitmentBz, nil, selfHeight)
	return types.ChannelCdc.MustMarshalJSON(resp), nil
}


// PacketCommitments implements the Query/PacketCommitments gRPC method
func (q Keeper) PacketCommitments(ctx sdk.Context, requst abci.RequestQuery) ([]byte, error) {
	var req types.QueryPacketCommitmentsRequest
	if err := types.ChannelCdc.UnmarshalJSON(requst.Data, &req); err != nil {
		return nil, err
	}

	if err := validateRPCRequest(req.PortId, req.ChannelId); err != nil {
		return nil, err
	}

	commitments := []*types.PacketState{}
	channelStore := store.NewPrefixStore(ctx.KVStore(q.storeKey), []byte(host.PacketCommitmentPrefixPath(req.PortId, req.ChannelId)))

	pageRes, err := pagination.Paginate(channelStore, req.Pagination, func(key, value []byte) error {
		keySplit := strings.Split(string(key), "/")

		sequence, err := strconv.ParseUint(keySplit[len(keySplit)-1], 10, 64)
		if err != nil {
			return err
		}

		commitment := types.NewPacketState(req.PortId, req.ChannelId, sequence, value)
		commitments = append(commitments, &commitment)
		return nil
	})

	if err != nil {
		return nil, err
	}

	selfHeight := clienttypes.GetSelfHeight(ctx)
	resp := types.QueryPacketCommitmentsResponse{
		Commitments: commitments,
		Pagination:  pageRes,
		Height:      selfHeight,
	}
	return types.ChannelCdc.MustMarshalJSON(resp), nil
}


// UnreceivedPackets implements the Query/UnreceivedPackets gRPC method. Given
// a list of counterparty packet commitments, the querier checks if the packet
// has already been received by checking if a receipt exists on this
// chain for the packet sequence. All packets that haven't been received yet
// are returned in the response
// Usage: To use this method correctly, first query all packet commitments on
// the sending chain using the Query/PacketCommitments gRPC method.
// Then input the returned sequences into the QueryUnreceivedPacketsRequest
// and send the request to this Query/UnreceivedPackets on the **receiving**
// chain. This gRPC method will then return the list of packet sequences that
// are yet to be received on the receiving chain.
//
// NOTE: The querier makes the assumption that the provided list of packet
// commitments is correct and will not function properly if the list
// is not up to date. Ideally the query height should equal the latest height
// on the counterparty's client which represents this chain.
func (q Keeper) UnreceivedPackets(ctx sdk.Context, requst abci.RequestQuery) ([]byte, error) {
	var req types.QueryUnreceivedPacketsRequest
	if err := types.ChannelCdc.UnmarshalJSON(requst.Data, &req); err != nil {
		return nil, err
	}

	if err := validateRPCRequest(req.PortId, req.ChannelId); err != nil {
		return nil, err
	}

	var unreceivedSequences = []uint64{}
	for i, seq := range req.PacketCommitmentSequences {
		if seq == 0 {
			return nil, status.Errorf(codes.InvalidArgument, "packet sequence %d cannot be 0", i)
		}

		// if packet receipt exists on the receiving chain, then packet has already been received
		if _, found := q.GetPacketReceipt(ctx, req.PortId, req.ChannelId, seq); !found {
			unreceivedSequences = append(unreceivedSequences, seq)
		}

	}

	selfHeight := clienttypes.GetSelfHeight(ctx)
	resp := types.QueryUnreceivedPacketsResponse{
		Sequences: unreceivedSequences,
		Height:    selfHeight,
	}
	return types.ChannelCdc.MustMarshalJSON(resp), nil
}


func validateRPCRequest(portID, channelID string) error {
	if err := host.PortIdentifierValidator(portID); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	if err := host.ChannelIdentifierValidator(channelID); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	return nil
}
