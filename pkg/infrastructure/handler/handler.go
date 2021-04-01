package handler

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/infrastructure/keeper"
	infrastructure "github.com/ci123chain/ci123chain/pkg/infrastructure/types"
)

func NewHandler(k keeper.InfrastructureKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case *infrastructure.MsgStoreContent:
			return HandleMsgStoreContent(ctx, k, *msg)
		default:
			errMsg := fmt.Sprintf("unrecognized supply message types: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}


func HandleMsgStoreContent(ctx sdk.Context, k keeper.InfrastructureKeeper, msg infrastructure.MsgStoreContent) sdk.Result {
	em := ctx.EventManager()
	//em.EmitEvents(sdk.Events{
	//	sdk.NewEvent(transfer.EventType,
	//		sdk.NewAttribute([]byte(sdk.AttributeKeyMethod), []byte(infrastructure.EventStoreContent)),
	//		sdk.NewAttribute([]byte(sdk.AttributeKeyModule), []byte(infrastructure.AttributeValueModule)),
	//		sdk.NewAttribute([]byte(sdk.AttributeKeySender), []byte(msg.FromAddress.String())),
	//	),
	//})

	k.SetContent(ctx, []byte(msg.Key), msg.Content)
	return sdk.Result{ Events: em.Events(), }
}