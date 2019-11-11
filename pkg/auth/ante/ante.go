package ante

import (
	"fmt"
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/auth"
	"github.com/tanhuiya/ci123chain/pkg/transaction"
)

func NewAnteHandler( authKeeper auth.AuthKeeper) types.AnteHandler {
	return func(ctx types.Context, tx types.Tx, simulate bool) (newCtx types.Context, res types.Result, abort bool) {

		return newCtx, types.Result{}, false

		commonTx, ok := tx.(*transaction.TransferTx)
		if !ok {
			newCtx := ctx.WithGasMeter(types.NewGasMeter(0))
			return newCtx, types.ErrInternal("undefined transaction Type ").Result(), true
		}

		params := authKeeper.GetParams(ctx)

		if ctx.IsCheckTx() && !simulate {
			// 检查 mempool 的gas limit
		}
		newCtx = SetGasMeter(simulate, ctx, commonTx.Common.Gas)

		defer func() {
			if r := recover(); r != nil {
				switch rType := r.(type) {
				case types.ErrorOutOfGas:
					log := fmt.Sprintf("out of gas in location: %v; gasWanted: %d, gasUsed: %d",
						rType.Descriptor, commonTx.Common.Gas, newCtx.GasMeter().GasConsumed(),
						)
					res = types.ErrOutOfGas(log).Result()
					res.GasWanted = commonTx.Common.Gas
					res.GasUsed = newCtx.GasMeter().GasConsumed()
					abort = true
				default:
					panic(r)
				}
			}
		}()


		if err := tx.ValidateBasic(); err != nil {
			return newCtx, err.Result(), true
		}

		newCtx.GasMeter().ConsumeGas( uint64(params.TxSizeCostPerByte) * types.Gas(len(newCtx.TxBytes())), "txsize")

		return newCtx, types.Result{GasWanted: commonTx.Common.Gas}, false
	}
}

func SetGasMeter(simulate bool, ctx types.Context, gaslimit uint64) types.Context {
	if simulate || ctx.BlockHeight() == 0 {
		return ctx.WithGasMeter(types.NewInfiniteGasMeter())
	}
	return ctx.WithGasMeter(types.NewGasMeter(gaslimit))
}
