package keeper

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryClientState:
			return clientState(ctx, req, k)
		case types.QueryConsensusState:
			return consensusState(ctx, req, k)
		case types.QueryPacketCommitment:
			return packetCommitment(ctx, req, k)
		case types.QueryPacketCommitments:
			return packetCommitments(ctx, req, k)
		case types.UnreceivedPacketd:
			return unreceivedPackets(ctx, req, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown query endpoint")
		}
	}
}


// ClientState implements the IBC QueryServer interface
func clientState(c sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	resp, err := keeper.ClientKeeper.ClientState(c, req)
	return resp, err
}

// ClientState implements the IBC QueryServer interface
func consensusState(c sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	resp, err := keeper.ClientKeeper.ConsensusState(c, req)
	return resp, err
}

func packetCommitment(c sdk.Context, req abci.RequestQuery, keeper Keeper)  ([]byte, sdk.Error) {
	resp, err := keeper.ChannelKeeper.PacketCommitment(c, req)
	return resp, err
}

func packetCommitments(c sdk.Context, req abci.RequestQuery, keeper Keeper)  ([]byte, sdk.Error) {
	resp, err := keeper.ChannelKeeper.PacketCommitments(c, req)
	return resp, err
}

// UnreceivedPackets implements the IBC QueryServer interface
func unreceivedPackets(c sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	return keeper.ChannelKeeper.UnreceivedPackets(c, req)
}
