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
			errMsg := fmt.Sprintf("unrecognized supply message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}


func HandleMsgStoreContent(ctx sdk.Context, k keeper.InfrastructureKeeper, msg infrastructure.MsgStoreContent) sdk.Result {

	k.SetContent(ctx, []byte(msg.Key), msg.Content)
	return sdk.Result{}
}