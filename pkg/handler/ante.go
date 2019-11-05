package handler

import (
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/account/keeper"
)

func NewAnteHandler(am keeper.AccountKeeper) types.AnteHandler {
	return func(ctx types.Context, tx types.Tx, simulate bool) (newCtx types.Context, result types.Result, abort bool) {

		//commonTx, ok := tx.(transaction.CommonTx)
		//if !ok {
		//	newCtx := ctx.WithGasMeter(types.NewGasMeter(0))
		//	return newCtx, types.ErrInternal("undefined transaction Type ").Result(), true
		//}
		//
		//if ctx.IsCheckTx() && !simulate {
		//	// 检查 mempool 的gas limit
		//}
		//newCtx = SetGasMeter(simulate, ctx, commonTx.Gas)
		//
		//if err := tx.ValidateBasic(); err != nil {
		//	return newCtx, err.Result(), true
		//}

		//newCtx.GasMeter().ConsumeGas()

		return newCtx, types.Result{}, false
	}
}

func SetGasMeter(simulate bool, ctx types.Context, gaslimit uint64) types.Context {
	if simulate || ctx.BlockHeight() == 0 {
		return ctx.WithGasMeter(types.NewInfiniteGasMeter())
	}
	return ctx.WithGasMeter(types.NewGasMeter(gaslimit))
}
