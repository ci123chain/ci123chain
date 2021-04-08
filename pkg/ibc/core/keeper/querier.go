package keeper

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case types.QueryClientState:
			return clientState(ctx, req, k)
		case types.QueryClientStates:
			return clientStates(ctx, req, k)
		case types.QueryClientConsensusState:
			return consensusState(ctx, req, k)

		case types.QueryConnections:
			return connections(ctx, req, k)
		case types.QueryChannels:
			return channels(ctx, req, k)

		case types.QueryPacketCommitment:
			return packetCommitment(ctx, req, k)
		case types.QueryPacketCommitments:
			return packetCommitments(ctx, req, k)
		case types.UnreceivedPacketd:
			return unreceivedPackets(ctx, req, k)
		default:
			return nil, sdkerrors.ErrUnknownRequest
		}
	}
}


// ClientState implements the IBC QueryServer interface
func clientState(c sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	resp, err := keeper.ClientKeeper.ClientState(c, req)
	return resp, err
}

// ClientState implements the IBC QueryServer interface
func clientStates(c sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	resp, err := keeper.ClientKeeper.ClientStates(c, req)
	return resp, err
}

// ClientState implements the IBC QueryServer interface
func consensusState(c sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	resp, err := keeper.ClientKeeper.ConsensusState(c, req)
	return resp, err
}

// ClientStates implements the IBC QueryServer interface
func connections(c sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	resp, err := keeper.ConnectionKeeper.Connections(c, req)
	return resp, err
}

// ClientStates implements the IBC QueryServer interface
func channels(c sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	resp, err := keeper.ChannelKeeper.Channels(c, req)
	return resp, err
}

func packetCommitment(c sdk.Context, req abci.RequestQuery, keeper Keeper)  ([]byte, error) {
	resp, err := keeper.ChannelKeeper.PacketCommitment(c, req)
	return resp, err
}

func packetCommitments(c sdk.Context, req abci.RequestQuery, keeper Keeper)  ([]byte, error) {
	resp, err := keeper.ChannelKeeper.PacketCommitments(c, req)
	return resp, err
}

// UnreceivedPackets implements the IBC QueryServer interface
func unreceivedPackets(c sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	return keeper.ChannelKeeper.UnreceivedPackets(c, req)
}
