package ante

import (
	"fmt"
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/account"
	"github.com/tanhuiya/ci123chain/pkg/account/exported"
	"github.com/tanhuiya/ci123chain/pkg/auth"
	"github.com/tanhuiya/ci123chain/pkg/auth/types"
	fc "github.com/tanhuiya/ci123chain/pkg/fc"
	"github.com/tanhuiya/ci123chain/pkg/transaction"
)
const price uint64 = 1

func NewAnteHandler( authKeeper auth.AuthKeeper, ak account.AccountKeeper, fck fc.FcKeeper) sdk.AnteHandler {
	return func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, res sdk.Result, abort bool) {

		stdTx, ok := tx.(transaction.Transaction)
		if !ok {
			// Set a gas meter with limit 0 as to prevent an infinite gas meter attack
			// during runTx.
			newCtx = SetGasMeter(simulate, ctx, 0)
			return newCtx, transaction.ErrInvalidTx(types.DefaultCodespace, "tx must be StdTx").Result(), true
		}
		address := stdTx.GetFromAddress()
		acc := ak.GetAccount(ctx, address)
		if acc == nil {
			newCtx := ctx.WithGasMeter(sdk.NewGasMeter(0))
			return newCtx, transaction.ErrInvalidTx(types.DefaultCodespace, "Invalid account").Result(), true
		}
		/*accountSequence := acc.GetSequence()
		txNonce := stdTx.GetNonce()
		if txNonce != accountSequence {
			newCtx := ctx.WithGasMeter(sdk.NewGasMeter(0))
			return newCtx, transaction.ErrInvalidTx(types.DefaultCodespace, "Unexpected nonce ").Result(), true
		}*/

		params := authKeeper.GetParams(ctx)
		// Ensure that the provided fees meet a minimum threshold for the validator,
		// if this is a CheckTx. This is only for local mempool purposes, and thus
		// is only ran on check tx.
		if ctx.IsCheckTx() && !simulate {
			res := EnsureSufficientMempoolFees()
			if !res.IsOK() {
				return newCtx, res, true
			}
		}
		gas := stdTx.GetGas()//用户期望的gas值 g.limit
		newCtx = SetGasMeter(simulate, ctx, gas)//设置为GasMeter的gasLimit,成为用户可承受的gas上限.
		//pms.TxSizeCostPerByte*sdk.Gas(len(newCtx.TxBytes()))
		var sg uint64 = 1
		newCtx.GasMeter().ConsumeGas(sg, "txSize") //计算最终的gas值. g.consumed

		//计算fee
		gasPrice := 2*price
		fee := newCtx.GasMeter().GasConsumed() * gasPrice
		fmt.Println("============================================")
		fmt.Println(fee)
		getFee := sdk.Coin(fee)

		newCtx = SetGasMeter(simulate, ctx, gas)
		// AnteHandlers must have their own defer/recover in order for the BaseApp
		// to know how much gas was used! This is because the GasMeter is created in
		// the AnteHandler, but if it panics the context won't be set properly in
		// runTx's recover call.
		defer func() {
			if r := recover(); r != nil {
				switch rType := r.(type) {
				case sdk.ErrorOutOfGas:
					log := fmt.Sprintf(
						"out of gas in location: %v; gasWanted: %d, gasUsed: %d",
						rType.Descriptor, gas, newCtx.GasMeter().GasConsumed(),
					)
					res = sdk.ErrOutOfGas(log).Result()

					res.GasWanted = gas
					res.GasUsed = newCtx.GasMeter().GasConsumed()
					abort = true
				default:
					panic(r)
				}
			}
		}()
		if err := tx.ValidateBasic(); err != nil {
			return newCtx, types.ErrTxValidateBasic(types.DefaultCodespace, err).Result(), true
		}

		newCtx.GasMeter().ConsumeGas(params.TxSizeCostPerByte*sdk.Gas(len(newCtx.TxBytes())), "txSize")

		res = DeductFees(acc, getFee, ak, ctx)

		if !res.IsOK() {
			return newCtx, res, true
		}
		//存储奖励金

		fck.AddCollectedFees(newCtx, getFee)

		return newCtx, sdk.Result{GasWanted:gas,GasUsed:fee}, false
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


func DeductFees(acc exported.Account, fee sdk.Coin, ak account.AccountKeeper, ctx sdk.Context) (sdk.Result) {
	coin := acc.GetCoin()
	newCoins, ok := coin.SafeSub(fee)
	if !ok {
		return sdk.ErrInsufficientFunds(
			fmt.Sprintf("insufficient funds to pay for fees; %s < %s", coin, fee),
		).Result()
	}

	if err := acc.SetCoin(newCoins); err != nil {
		return account.ErrSetAccount(types.DefaultCodespace, err).Result()
	}
	ak.SetAccount(ctx, acc)

	return sdk.Result{}
}