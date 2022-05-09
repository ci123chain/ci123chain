package ante

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/account"
	types2 "github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/auth"
	"github.com/ci123chain/ci123chain/pkg/cryptosuite"
	"math/big"

	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/ci123chain/ci123chain/pkg/supply"
	"github.com/ci123chain/ci123chain/pkg/util"
)
const (
	GasSimulateCost sdk.Gas = 200
)
var simSecp256k1Sig [65]byte

//const unit = 1000
func NewAnteHandler( authKeeper auth.AuthKeeper, ak account.AccountKeeper, sk supply.Keeper) sdk.AnteHandler {
	return func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, res sdk.Result, err error, abort bool) {
		if simulate {
			params := authKeeper.GetParams(ctx)
			newCtx = SetGasMeter(simulate, ctx, 0)
			tx.SetSignature(simSecp256k1Sig[:])
			newCtx.GasMeter().ConsumeGas(params.TxSizeCostPerByte*sdk.Gas(len(tx.Bytes())), "txSize")
			newCtx.GasMeter().ConsumeGas(GasSimulateCost, "simulate cost")
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
			eth := cryptosuite.NewEccK1()
			valid, err := eth.Verify(tx.GetFromAddress().Bytes(), tx.GetSignBytes(), tx.GetSignature())
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
		newCtx = SetGasMeter(simulate, ctx, gas)//设置为GasMeter的gasLimit,成为用户可承受的gas上限.

		//检查是否足够支付gas limit, 并预先扣除
		gasTotal := sdk.CalculateGas(gas)

		if err := sk.SendCoinsFromAccountToModule(newCtx, signer, auth.FeeCollectorName, sdk.NewCoins(sdk.NewChainCoin(gasTotal))); err != nil {
			return newCtx, sdk.Result{}, err, true
		}

		//Calculate fee of tx size
		newCtx.GasMeter().ConsumeGas(params.TxSizeCostPerByte * sdk.Gas(len(newCtx.TxBytes())), "txSize")

		// account sequence{nonce} + 1
		if !ctx.IsCheckTx() {
			accountSequence := acc.GetSequence()
			txNonce := tx.GetNonce()
			if txNonce != accountSequence {
				newCtx := ctx.WithGasMeter(sdk.NewGasMeter(0))
				return newCtx, sdk.Result{}, sdkerrors.ErrInvalidParam.Wrap("nonce dismatch"), true
			}

			acc = ak.GetAccount(ctx, signer)
			nowSequence := accountSequence + 1
			_ = acc.SetSequence(nowSequence)
			ak.SetAccount(newCtx, acc)
		}

		return newCtx, sdk.Result{GasWanted:gas, GasUsed: newCtx.GasMeter().GasConsumed()}, nil, false
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

//
//func DeductFees(acc exported.Account, fee sdk.Coin, ak account.AccountKeeper, ctx sdk.Context) error{
//	coin := acc.GetCoins()
//	newCoins, hasNeg := coin.SafeSub(sdk.NewCoins(fee))
//	if hasNeg {
//		return sdkerrors.ErrInsufficientFunds
//	}
//
//	if err := acc.SetCoins(newCoins); err != nil {
//		return sdkerrors.ErrAccountSetCoinFailed
//	}
//	ak.SetAccount(ctx, acc)
//
//	return nil
//}
//
//func ReturnFees(acc exported.Account, restFee sdk.Coins, ak account.AccountKeeper, ctx sdk.Context) error {
//	coin := acc.GetCoins()
//	newCoins:= coin.Add(restFee)
//
//	if err := acc.SetCoins(newCoins); err != nil {
//		return sdkerrors.ErrAccountSetCoinFailed
//	}
//	ak.SetAccount(ctx, acc)
//
//	return nil
//}