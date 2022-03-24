package _defer

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/account"
	types2 "github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/auth"
	"github.com/ci123chain/ci123chain/pkg/supply"
	"github.com/ci123chain/ci123chain/pkg/transaction"
	"github.com/ci123chain/ci123chain/pkg/util"
	"math/big"
)

//const unit = 1000
func NewDeferHandler( ak account.AccountKeeper, sk supply.Keeper) sdk.DeferHandler {
	return func(ctx sdk.Context, tx sdk.Tx, out bool, simulate bool) (res sdk.Result) {
		if out || simulate{
			return
		}

		//返还剩余gas
		var gasUsed uint64
		var signer sdk.AccAddress
		stdTx := tx.(transaction.Transaction)
		gasWanted := stdTx.GetGas()

		defer func() {
			if r := recover(); r != nil {
				res.GasUsed = gasWanted
				return
			}
		}()

		if etx, ok := tx.(*types2.MsgEthereumTx); ok {
			from, err := etx.VerifySig(big.NewInt(util.CHAINID))
			if err != nil {
				panic(err)
			}
			signer = sdk.AccAddress{from}
		} else {
			signer = tx.GetFromAddress()
		}

		acc := ak.GetAccount(ctx, signer)
		if acc == nil {
			return
		}

		gasUsed = ctx.GasMeter().GasConsumed()
		restgas := gasWanted - gasUsed
		if restgas == 0 {
			return 
		}
		calculateGas := sdk.CalculateGas(restgas)
		restFee := sdk.NewChainCoin(calculateGas)
		err := sk.SendCoinsFromModuleToAccount(ctx, auth.FeeCollectorName, signer, sdk.NewCoins(restFee))
		if err != nil {
			return
		}
		res.GasUsed = gasUsed
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