package core

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	channeltypes "github.com/ci123chain/ci123chain/pkg/ibc/core/channel/types"
	clienttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/clients/types"
	connectiontypes "github.com/ci123chain/ci123chain/pkg/ibc/core/connection/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/keeper"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		var err error
		var res interface{}
		switch msg := msg.(type) {
		case *clienttypes.MsgCreateClient:
			res, err = k.CreateClient(ctx, msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *clienttypes.MsgUpdateClient:
			res, err = k.UpdateClient(ctx, msg)
			return sdk.WrapServiceResult(ctx, res, err)

			// IBC connection msgs
		case *connectiontypes.MsgConnectionOpenInit:
			res, err = k.ConnectionOpenInit(ctx, msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *connectiontypes.MsgConnectionOpenTry:
			res, err = k.ConnectionOpenTry(ctx, msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *connectiontypes.MsgConnectionOpenAck:
			res, err = k.ConnectionOpenAck(ctx, msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *connectiontypes.MsgConnectionOpenConfirm:
			res, err = k.ConnectionOpenConfirm(ctx, msg)
			return sdk.WrapServiceResult(ctx, res, err)

		// IBC channel msgs
		case *channeltypes.MsgChannelOpenInit:
			res, err = k.ChannelOpenInit(ctx, msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *channeltypes.MsgChannelOpenTry:
			res, err = k.ChannelOpenTry(ctx, msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *channeltypes.MsgChannelOpenAck:
			res, err = k.ChannelOpenAck(ctx, msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *channeltypes.MsgChannelOpenConfirm:
			res, err = k.ChannelOpenConfirm(ctx, msg)
			return sdk.WrapServiceResult(ctx, res, err)

			// IBC packet msgs get routed to the appropriate module callback
		case *channeltypes.MsgRecvPacket:
			res, err := k.RecvPacket(ctx, msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *channeltypes.MsgAcknowledgement:
			res, err := k.Acknowledgement(ctx, msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *channeltypes.MsgTimeout:
			res, err := k.Timeout(ctx, msg)
			return sdk.WrapServiceResult(ctx, res, err)

		//case *channeltypes.MsgTimeoutOnClose:
		//	res, err := k.TimeoutOnClose(ctx, msg)
		//	return sdk.WrapServiceResult(ctx, res, err)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized IBC message type: %T", msg)
		}
	}
}