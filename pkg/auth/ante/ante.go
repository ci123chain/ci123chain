package ante

import (
	"fmt"
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/account"
	"github.com/tanhuiya/ci123chain/pkg/account/exported"
	"github.com/tanhuiya/ci123chain/pkg/auth"
	fc "github.com/tanhuiya/ci123chain/pkg/fc"
	"github.com/tanhuiya/ci123chain/pkg/transaction"
)
const price uint64 = 1
func NewAnteHandler( authKeeper auth.AuthKeeper, ak account.AccountKeeper, fck fc.FcKeeper) types.AnteHandler {
	return func(ctx types.Context, tx types.Tx, simulate bool) (newCtx types.Context, res types.Result, abort bool) {

		stdTx, ok := tx.(transaction.Transaction)
		if !ok {
			// Set a gas meter with limit 0 as to prevent an infinite gas meter attack
			// during runTx.
			newCtx = SetGasMeter(simulate, ctx, 0)
			return newCtx, types.ErrInternal("tx must be StdTx").Result(), true
		}
		address := stdTx.GetFromAddress()
		acc := ak.GetAccount(ctx, address)
		if acc == nil {
			newCtx := ctx.WithGasMeter(types.NewGasMeter(0))
			return newCtx, types.ErrInternal("Invalid account").Result(), true
		}
		accountSequence := acc.GetSequence()
		txNonce := stdTx.GetNonce()
		if txNonce != accountSequence {
			newCtx := ctx.WithGasMeter(types.NewGasMeter(0))
			return newCtx, types.ErrInternal("Undefined transfer Type ").Result(), true
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
		gas := stdTx.GetGas()//用户期望的gas值 g.limit

		newCtx = SetGasMeter(simulate, ctx, gas)//设置为GasMeter的gasLimit,成为用户可承受的gas上限.
		//pms.TxSizeCostPerByte*types.Gas(len(newCtx.TxBytes()))
		var sg uint64 = 1
		newCtx.GasMeter().ConsumeGas(sg, "txSize") //计算最终的gas值. g.consumed

		//计算fee
		gasPrice := 2*price
		fee := newCtx.GasMeter().GasConsumed() * gasPrice
		getFee := types.Coin(fee)

		newCtx = SetGasMeter(simulate, ctx, gas)

		// AnteHandlers must have their own defer/recover in order for the BaseApp
		// to know how much gas was used! This is because the GasMeter is created in
		// the AnteHandler, but if it panics the context won't be set properly in
		// runTx's recover call.
		defer func() {
			if r := recover(); r != nil {
				switch rType := r.(type) {
				case types.ErrorOutOfGas:
					log := fmt.Sprintf(
						"out of gas in location: %v; gasWanted: %d, gasUsed: %d",
						rType.Descriptor, gas, newCtx.GasMeter().GasConsumed(),
					)
					res = types.ErrOutOfGas(log).Result()

					res.GasWanted = gas
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

		newCtx.GasMeter().ConsumeGas(params.TxSizeCostPerByte*types.Gas(len(newCtx.TxBytes())), "txSize")

		//if res := ValidateMemo(stdTx, params); !res.IsOK() {
		//	return newCtx, res, true
		//}

		// stdSigs contains the sequence number, account number, and signatures.
		// When simulating, this would just be a 0-length slice.
		//signerAddrs := stdTx.GetSigners()
		//signerAccs := make([]Account, len(signerAddrs))
		//isGenesis := ctx.BlockHeight() == 0

		// fetch first signer, who's going to pay the fees
		//signerAccs[0], res = GetSignerAcc(newCtx, ak, signerAddrs[0])
		//if !res.IsOK() {
		//	return newCtx, res, true
		//}
		res = DeductFees(acc, getFee)

		//signerAccs[0], res = DeductFees(ctx.BlockHeader().Time, signerAccs[0], stdTx.Fee)
		if !res.IsOK() {
			return newCtx, res, true
		}
		//存储奖励金

		fck.AddCollectedFees(newCtx, getFee)
		//g := fck.GetCollectedFees(newCtx)
		//fmt.Println(g)

		//给验证者的账户sequence+1
		//processSig(newCtx, acc, simulate, params)
		// stdSigs contains the sequence number, account number, and signatures.
		// When simulating, this would just be a 0-length slice.
		//stdSigs := stdTx.GetSignatures()

		//for i := 0; i < len(stdSigs); i++ {
		//	// skip the fee payer, account is cached and fees were deducted already
		//	if i != 0 {
		//		signerAccs[i], res = GetSignerAcc(newCtx, ak, signerAddrs[i])
		//		if !res.IsOK() {
		//			return newCtx, res, true
		//		}
		//	}
		//
		//	// check signature, return account with incremented nonce
		//	signBytes := GetSignBytes(newCtx.ChainID(), stdTx, signerAccs[i], isGenesis)
		//	signerAccs[i], res = processSig(newCtx, signerAccs[i], stdSigs[i], signBytes, simulate, params)
		//	if !res.IsOK() {
		//		return newCtx, res, true
		//	}
		//
		//	ak.SetAccount(newCtx, signerAccs[i])
		//}


		return newCtx, types.Result{GasWanted:gas,GasUsed:fee}, false
	}
}

func SetGasMeter(simulate bool, ctx types.Context, gaslimit uint64) types.Context {
	if simulate || ctx.BlockHeight() == 0 {
		return ctx.WithGasMeter(types.NewInfiniteGasMeter())
	}
	return ctx.WithGasMeter(types.NewGasMeter(gaslimit))
}

func EnsureSufficientMempoolFees() types.Result {
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

	return types.Result{}
}

func DeductFees(acc exported.Account, fee types.Coin) (types.Result) {
	coin := acc.GetCoin()
	newCoins, ok := coin.SafeSub(fee)
	if !ok {
		return types.ErrInsufficientFunds(
			fmt.Sprintf("insufficient funds to pay for fees; %s < %s", coin, fee),
		).Result()
	}

	// Validate the account has enough "spendable" coins as this will cover cases
	// such as vesting accounts.
	//spendableCoins := acc.SpendableCoins(blockTime)
	//if _, hasNeg := spendableCoins.SafeSub(feeAmount); hasNeg {
	//	return nil, sdk.ErrInsufficientFunds(
	//		fmt.Sprintf("insufficient funds to pay for fees; %s < %s", spendableCoins, feeAmount),
	//	).Result()
	//}

	if err := acc.SetCoin(newCoins); err != nil {
		return types.ErrInternal(err.Error()).Result()
	}

	return types.Result{}
}