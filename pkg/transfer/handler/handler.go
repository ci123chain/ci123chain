package handler

import (
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/account/keeper"
	"github.com/ci123chain/ci123chain/pkg/transfer"
	"reflect"
)

func NewHandler(am keeper.AccountKeeper) types.Handler {
	return func(ctx types.Context, msg types.Msg) types.Result {
		switch msg := msg.(type) {
		case *transfer.MsgTransfer:
			return handlerMsgTransfer(ctx, am, msg)
		default:
			errMsg := "Unrecognized Tx type: " + reflect.TypeOf(msg).Name()
			return types.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handlerMsgTransfer(ctx types.Context, am keeper.AccountKeeper, msg *transfer.MsgTransfer) types.Result {
	if err := am.Transfer(ctx, msg.FromAddress, msg.To, msg.Amount); err != nil {
		return err.Result()
	}
	return types.Result{}
}