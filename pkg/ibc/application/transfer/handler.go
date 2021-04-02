package transfer

import (
	"encoding/json"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/application/transfer/types"
)

// NewHandler returns sdk.Handler for IBC token transfer module messages
func NewHandler(k types.MsgServer) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgTransfer:
			res, err := k.Transfer(ctx, msg)
			if err != nil {
				return sdk.NewError("ibc", 500, err.Error()).Result()
			}
			res1, _ := json.Marshal(res)
			return sdk.Result{Data: res1}

		default:
			return sdk.NewError("ibc", 500, "unrecognized ICS-20 transfer message type: %T").Result()
		}
	}
}
