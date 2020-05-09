package _defer

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/account"
	"github.com/ci123chain/ci123chain/pkg/auth/ante"
	"github.com/ci123chain/ci123chain/pkg/transaction"
)

const Price uint64 = 1
//const unit = 1000
func NewDeferHandler( ak account.AccountKeeper) sdk.DeferHandler {
	return func(ctx sdk.Context, tx sdk.Tx, out bool) (res sdk.Result) {
		//实际扣gas
		var gasused uint64
		stdTx, _ := tx.(transaction.Transaction)
		address := stdTx.GetFromAddress()
		acc := ak.GetAccount(ctx, address)

		if out {
			gasused = stdTx.GetGas()
		} else {
			gasused = ctx.GasMeter().GasConsumed()
		}
		fee := sdk.NewUInt64Coin(gasused)
		res = ante.DeductFees(acc, fee, ak, ctx)
		return
	}
}

func SetGasMeter(simulate bool, ctx sdk.Context, gaslimit uint64) sdk.Context {
	if simulate || ctx.BlockHeight() == 0 {
		return ctx.WithGasMeter(sdk.NewInfiniteGasMeter())
	}
	return ctx.WithGasMeter(sdk.NewGasMeter(gaslimit))
}

func EnsureSufficientMempoolFees() sdk.Result {
	//minGasPrices := ctx.MinGasPrices()
	//if !minGasPrices.IsZero() {
	//	requiredFees := make(sdk.Coins, len(minGasPrices))
	//
	//	// Determine the required fees by multiplying each required minimum gas
	//	// price by the gas limit, where fee = ceil(minGasPrice * gasLimit).
	//	glDec := sdk.NewDec(int64(stdFee.Gas))
	//	for i, gp := range minGasPrices {
	//		fee := gp.Amount.Mul(glDec)
	//		requiredFees[i] = sdk.NewCoin(gp.Denom, fee.Ceil().RoundInt())
	//	}
	//
	//	if !stdFee.Amount.IsAnyGTE(requiredFees) {
	//		return sdk.ErrInsufficientFee(
	//			fmt.Sprintf(
	//				"insufficient fees; got: %q required: %q", stdFee.Amount, requiredFees,
	//			),
	//		).Result()
	//	}
	//}

	return sdk.Result{}
}