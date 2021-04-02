package core

import (
	"encoding/json"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	channeltypes "github.com/ci123chain/ci123chain/pkg/ibc/core/channel/types"
	clienttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/clients/types"
	connectiontypes "github.com/ci123chain/ci123chain/pkg/ibc/core/connection/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/keeper"
	"github.com/pkg/errors"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		var err error
		var res interface{}
		switch msg := msg.(type) {
		case *clienttypes.MsgCreateClient:
			res, err = k.CreateClient(ctx, msg)
			break
		case *clienttypes.MsgUpdateClient:
			res, err = k.UpdateClient(ctx, msg)
			break
			// IBC connection msgs
		case *connectiontypes.MsgConnectionOpenInit:
			res, err = k.ConnectionOpenInit(ctx, msg)
			break
		case *connectiontypes.MsgConnectionOpenTry:
			res, err = k.ConnectionOpenTry(ctx, msg)
			break
		case *connectiontypes.MsgConnectionOpenAck:
			res, err = k.ConnectionOpenAck(ctx, msg)
			break
		case *connectiontypes.MsgConnectionOpenConfirm:
			res, err = k.ConnectionOpenConfirm(ctx, msg)
			break
		// IBC channel msgs
		case *channeltypes.MsgChannelOpenInit:
			res, err = k.ChannelOpenInit(ctx, msg)
			break
		case *channeltypes.MsgChannelOpenTry:
			res, err = k.ChannelOpenTry(ctx, msg)
			break
		case *channeltypes.MsgChannelOpenAck:
			res, err = k.ChannelOpenAck(ctx, msg)
			break
		case *channeltypes.MsgChannelOpenConfirm:
			res, err = k.ChannelOpenConfirm(ctx, msg)
			break
		default:
			err = errors.Errorf("unrecognized ICS-20 transfer message type: %T")
		}
		if err != nil {
			return sdk.NewError("ibc", 501, err.Error()).Result()
		}
		res1, _ := json.Marshal(res)
		return sdk.Result{Data: res1}
	}
}