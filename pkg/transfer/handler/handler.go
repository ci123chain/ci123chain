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
			errMsg := "Unrecognized Tx types: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handlerMsgTransfer(ctx sdk.Context, am keeper.AccountKeeper, msg *transfer.MsgTransfer) sdk.Result {
	if err := am.Transfer(ctx, msg.FromAddress, msg.To, msg.Amount); err != nil {
		return err.Result()
	}
	em := ctx.EventManager()
	//em.EmitEvents(sdk.Events{
	//	sdk.NewEvent(transfer.EventType,
	//		sdk.NewAttribute([]byte(sdk.AttributeKeyMethod), []byte(transfer.AttributeValueTransfer)),
	//		sdk.NewAttribute([]byte(sdk.AttributeKeyModule), []byte(transfer.AttributeValueCategory)),
	//		sdk.NewAttribute([]byte(sdk.AttributeKeySender), []byte(msg.FromAddress.String())),
	//		sdk.NewAttribute([]byte(sdk.AttributeKeyReceiver), []byte(msg.To.String())),
	//		sdk.NewAttribute([]byte(sdk.AttributeKeyAmount), []byte(msg.Amount.Amount.String())),
	//	),
	//})
	return sdk.Result{ Events: em.Events(), }
}