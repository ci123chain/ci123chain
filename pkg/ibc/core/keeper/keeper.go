package keeper

import (
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	clientkeeper "github.com/ci123chain/ci123chain/pkg/ibc/core/clients/keeper"
	clienttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/clients/types"
	paramtypes "github.com/ci123chain/ci123chain/pkg/params/subspace"
)

type Keeper struct {
	cdc *codec.Codec
	ClientKeeper     clientkeeper.Keeper
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramSpace paramtypes.Subspace,
	stakingKeeper clienttypes.StakingKeeper,) *Keeper {
	clientKeeper := clientkeeper.NewKeeper(cdc, key, paramSpace, stakingKeeper)
	//connectionKeeper := connectionkeeper.NewKeeper(cdc, key, clientKeeper)
	//portKeeper := portkeeper.NewKeeper(scopedKeeper)
	//channelKeeper := channelkeeper.NewKeeper(cdc, key, clientKeeper, connectionKeeper, portKeeper, scopedKeeper)

	return &Keeper{
		cdc:              cdc,
		ClientKeeper:     clientKeeper,
		//ConnectionKeeper: connectionKeeper,
		//ChannelKeeper:    channelKeeper,
		//PortKeeper:       portKeeper,
	}
}

func (k Keeper) CreateClient(ctx sdk.Context, msg *clienttypes.MsgCreateClient) (*clienttypes.MsgCreateClientResponse, sdk.Error) {
	clientID, err := k.ClientKeeper.CreateClient(ctx, msg.ClientState, msg.ConsensusState)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			clienttypes.EventTypeCreateClient,
			sdk.NewAttributeString(clienttypes.AttributeKeyClientID, clientID),
			sdk.NewAttributeString(clienttypes.AttributeKeyClientType, msg.ClientState.ClientType()),
			sdk.NewAttributeString(clienttypes.AttributeKeyConsensusHeight, msg.ClientState.GetLatestHeight().String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttributeString(sdk.AttributeKeyModule, clienttypes.AttributeValueCategory),
		),
	})
	return &clienttypes.MsgCreateClientResponse{}, nil
}





















