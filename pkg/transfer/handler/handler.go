package handler

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/account/keeper"
	"github.com/ci123chain/ci123chain/pkg/transfer"
	"reflect"
)

func NewHandler(am keeper.AccountKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case *transfer.MsgTransfer:
			return handlerMsgTransfer(ctx, am, msg)
		default:
			errMsg := "Unrecognized Tx type: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handlerMsgTransfer(ctx sdk.Context, am keeper.AccountKeeper, msg *transfer.MsgTransfer) sdk.Result {
	if err := am.Transfer(ctx, msg.FromAddress, msg.To, msg.Amount); err != nil {
		return err.Result()
	}
	em := ctx.EventManager()
	em.EmitEvents(sdk.Events{
		sdk.NewEvent(transfer.EventType,
			sdk.NewAttribute(sdk.AttributeKeyMethod, transfer.AttributeValueTransfer),
			sdk.NewAttribute(sdk.AttributeKeyModule, transfer.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.FromAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyReceiver, msg.To.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.Amount.String()),
		),
	})
	return sdk.Result{ Events: em.Events(), }
}