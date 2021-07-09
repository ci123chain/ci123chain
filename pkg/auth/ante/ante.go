package ante

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/account"
	"github.com/ci123chain/ci123chain/pkg/account/exported"
	types2 "github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/auth"
	"github.com/ci123chain/ci123chain/pkg/cryptosuit"
	"math/big"

	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/ci123chain/ci123chain/pkg/supply"
	"github.com/ci123chain/ci123chain/pkg/util"
)
const (
	Price uint64 = 1
)

//const unit = 1000
func NewAnteHandler( authKeeper auth.AuthKeeper, ak account.AccountKeeper, sk supply.Keeper) sdk.AnteHandler {
	return func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, res sdk.Result, err error, abort bool) {
		if simulate {
			return
		}
		var signer sdk.AccAddress
		//check sign
		if etx, ok := tx.(*types2.MsgEthereumTx); ok {
			from, err := etx.VerifySig(big.NewInt(util.CHAINID))
			if err != nil {
				return newCtx, sdk.Result{}, sdkerrors.ErrorInvalidSigner, true
			}
			signer = sdk.HexToAddress(from.String())
		} else {
			eth := cryptosuit.NewETHSignIdentity()
			valid, err := eth.Verifier(tx.GetSignBytes(), tx.GetSignature(), nil, tx.GetFromAddress().Bytes())
			if !valid || err != nil {
				return newCtx, sdk.Result{}, sdkerrors.ErrorInvalidSigner, true
			}
			signer = tx.GetFromAddress()
		}

		acc := ak.GetAccount(ctx, signer)
		if acc == nil {
			newCtx := ctx.WithGasMeter(sdk.NewGasMeter(0))
			return newCtx, sdk.Result{}, sdkerrors.ErrAccountNotExist, true
		}
		accountSequence := acc.GetSequence()
		txNonce := tx.GetNonce()
		if txNonce != accountSequence {
			newCtx := ctx.WithGasMeter(sdk.NewGasMeter(0))
			return newCtx, sdk.Result{}, sdkerrors.ErrInvalidParam, true
		}

		params := authKeeper.GetParams(ctx)
		// Ensure that the provided fees meet a minimum threshold for the validator,
		// if this is a CheckTx. This is only for local mempool purposes, and thus
		// is only ran on check tx.
		if ctx.IsCheckTx() && !simulate {
			res := EnsureSufficientMempoolFees()
			if !res.IsOK() {
				return newCtx, res, sdkerrors.ErrInsufficientMempoolFess, true
			}
		}
		gas := tx.GetGas()//用户期望的gas值 g.limit
		//检查是否足够支付gas limit, 并预先扣除
		if acc.GetCoins().AmountOf(sdk.ChainCoinDenom).LT(sdk.NewIntFromBigInt(big.NewInt(int64(gas)))) {
			return newCtx, sdk.Result{}, sdkerrors.ErrInsufficientFunds,true
		}
		err = DeductFees(acc,sdk.NewUInt64Coin(sdk.ChainCoinDenom, gas),ak,ctx)
		if err != nil {
			return newCtx, sdk.Result{}, err, true
		}
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
		getFee := sdk.NewUInt64Coin(sdk.ChainCoinDenom, fee)

		//存储奖励金到feeCollector Module账户
		feeCollectorModuleAccount := sk.GetModuleAccount(ctx, auth.FeeCollectorName)
		newFee := feeCollectorModuleAccount.GetCoins().Add(sdk.NewCoins(getFee))
		err = feeCollectorModuleAccount.SetCoins(newFee)

		if err != nil {
			return newCtx, sdk.Result{}, sdkerrors.ErrModuleAccountSetCoinFailed, true
		}
		ak.SetAccount(ctx, feeCollectorModuleAccount)
		//fck.AddCollectedFees(newCtx, getFee)

		//account sequence + 1
		nowSequence := accountSequence + 1
		_ = acc.SetSequence(nowSequence)
		ak.SetAccount(ctx, acc)
		return newCtx, sdk.Result{GasWanted:gas,GasUsed:fee}, nil, false
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


func DeductFees(acc exported.Account, fee sdk.Coin, ak account.AccountKeeper, ctx sdk.Context) error{
	coin := acc.GetCoins()
	newCoins, hasNeg := coin.SafeSub(sdk.NewCoins(fee))
	if hasNeg {
		return sdkerrors.ErrInsufficientFunds
	}

	if err := acc.SetCoins(newCoins); err != nil {
		return sdkerrors.ErrAccountSetCoinFailed
	}
	ak.SetAccount(ctx, acc)

	return nil
}

func ReturnFees(acc exported.Account, restFee sdk.Coins, ak account.AccountKeeper, ctx sdk.Context) error {
	coin := acc.GetCoins()
	newCoins:= coin.Add(restFee)

	if err := acc.SetCoins(newCoins); err != nil {
		return sdkerrors.ErrAccountSetCoinFailed
	}
	ak.SetAccount(ctx, acc)

	return nil
}