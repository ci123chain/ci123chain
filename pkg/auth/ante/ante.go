package ante

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/account"
	"github.com/ci123chain/ci123chain/pkg/account/exported"
	types2 "github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/auth"
	"github.com/ci123chain/ci123chain/pkg/auth/types"
	"github.com/ci123chain/ci123chain/pkg/cryptosuit"
	"math/big"

	"github.com/ci123chain/ci123chain/pkg/supply"
	"github.com/ci123chain/ci123chain/pkg/transaction"
)
const (
	Price uint64 = 1
	ChainID int64 = 999
)

//const unit = 1000
func NewAnteHandler( authKeeper auth.AuthKeeper, ak account.AccountKeeper, sk supply.Keeper) sdk.AnteHandler {
	return func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, res sdk.Result, abort bool) {
		if simulate {
			return
		}
		var signer sdk.AccAddress
		//check sign
		if etx, ok := tx.(*types2.MsgEthereumTx); ok {
			from, err := etx.VerifySig(big.NewInt(ChainID))
			if err != nil {
				return newCtx, transaction.ErrInvalidTx(types.DefaultCodespace, "tx signature invalid").Result(), true
			}
			signer = sdk.AccAddress{from}
		} else {
			eth := cryptosuit.NewETHSignIdentity()
			valid, err := eth.Verifier(tx.GetSignBytes(), tx.GetSignature(), nil, tx.GetFromAddress().Bytes())
			if !valid || err != nil {
				return newCtx, transaction.ErrInvalidTx(types.DefaultCodespace, "tx signature invalid").Result(), true
			}
			signer = tx.GetFromAddress()
		}

		acc := ak.GetAccount(ctx, signer)
		if acc == nil {
			newCtx := ctx.WithGasMeter(sdk.NewGasMeter(0))
			return newCtx, transaction.ErrInvalidTx(types.DefaultCodespace, "Invalid account").Result(), true
		}
		accountSequence := acc.GetSequence()
		txNonce := tx.GetNonce()
		if txNonce != accountSequence {
			newCtx := ctx.WithGasMeter(sdk.NewGasMeter(0))
			return newCtx, transaction.ErrInvalidTx(types.DefaultCodespace, "Unexpected nonce ").Result(), true
		}

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
		gas := tx.GetGas()//用户期望的gas值 g.limit
		//检查是否足够支付gas limit, 并预先扣除
		if acc.GetCoin().Amount.LT(sdk.NewIntFromBigInt(big.NewInt(int64(gas)))) {
			return newCtx, sdk.ErrInsufficientCoins("Can't pay enough gasLimit").Result(),true
		}
		DeductFees(acc,sdk.NewUInt64Coin(gas),ak,ctx)
		newCtx = SetGasMeter(simulate, ctx, gas)//设置为GasMeter的gasLimit,成为用户可承受的gas上限.
		//pms.TxSizeCostPerByte*sdk.Gas(len(newCtx.TxBytes()))
		
		// AnteHandlers must have their own defer/recover in order for the BaseApp
		// to know how much gas was used! This is because the GasMeter is created in
		// the AnteHandler, but if it panics the context won't be set properly in
		// runTx's recover call.
		/*defer func() {
			if r := recover(); r != nil {
				switch rType := r.(types) {
				case sdk.ErrorOutOfGas:
					log := fmt.Sprintf(
						"out of gas in location: %v; gasWanted: %d, gasUsed: %d",
						rType.Descriptor, gas, newCtx.GasMeter().GasConsumed(),
					)
					res = sdk.ErrOutOfGas(log).Result()

					res.GasWanted = gas
					res.GasUsed = newCtx.GasMeter().GasConsumed()
					fmt.Println("-------- last ----------")
					fmt.Println(newCtx.GasMeter().GasConsumed())
					abort = true
				default:
					panic(r)
				}
			}
		}()*/
		//计算fee
		gasPrice := Price
		newCtx.GasMeter().ConsumeGas(params.TxSizeCostPerByte*sdk.Gas(len(newCtx.TxBytes())), "txSize")
		fee := newCtx.GasMeter().GasConsumed() * gasPrice
		getFee := sdk.NewUInt64Coin(fee)

		//存储奖励金到feeCollector Module账户
		feeCollectorModuleAccount := sk.GetModuleAccount(ctx, auth.FeeCollectorName)
		newFee := feeCollectorModuleAccount.GetCoin().Add(getFee)
		err := feeCollectorModuleAccount.SetCoin(newFee)

		if err != nil {
			fmt.Println("fee_collector module account set coin failed")
			panic(err)
		}
		ak.SetAccount(ctx, feeCollectorModuleAccount)
		//fck.AddCollectedFees(newCtx, getFee)

		//account sequence + 1
		nowSequence := accountSequence + 1
		_ = acc.SetSequence(nowSequence)
		ak.SetAccount(ctx, acc)
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


func DeductFees(acc exported.Account, fee sdk.Coin, ak account.AccountKeeper, ctx sdk.Context) sdk.Result {
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

func ReturnFees(acc exported.Account, restFee sdk.Coin, ak account.AccountKeeper, ctx sdk.Context) sdk.Result {
	coin := acc.GetCoin()
	newCoins:= coin.Add(restFee)

	if err := acc.SetCoin(newCoins); err != nil {
		return account.ErrSetAccount(types.DefaultCodespace, err).Result()
	}
	ak.SetAccount(ctx, acc)

	return sdk.Result{}
}