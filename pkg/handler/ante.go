package handler

import (
	"CI123Chain/pkg/abci/types"
	"CI123Chain/pkg/account"
)

func NewAnteHandler(am account.AccountMapper) types.AnteHandler {
	return func(ctx types.Context, tx types.Tx, simulate bool) (newCtx types.Context, result types.Result, abort bool) {
		return ctx, types.Result{}, false
	}
}
