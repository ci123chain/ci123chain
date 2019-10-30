package handler

import (
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/account"
)

func NewAnteHandler(am account.AccountMapper) types.AnteHandler {
	return func(ctx types.Context, tx types.Tx, simulate bool) (newCtx types.Context, result types.Result, abort bool) {
		return ctx, types.Result{}, false
	}
}
