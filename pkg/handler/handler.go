package handler

import (
	"gitlab.oneitfarm.com/blockchain/ci123chain/pkg/abci/types"
	"gitlab.oneitfarm.com/blockchain/ci123chain/pkg/account"
	"gitlab.oneitfarm.com/blockchain/ci123chain/pkg/transaction"
	"gitlab.oneitfarm.com/blockchain/ci123chain/pkg/db"
	"reflect"
)

func NewHandler(
	txm transaction.TxIndexMapper,
	am account.AccountMapper,
	sm *db.StateManager) types.Handler {
	return func(ctx types.Context, tx types.Tx) types.Result{
		ctx = ctx.WithTxIndex(txm.Get(ctx))
		defer func() {
			txm.Incr(ctx)
		}()
		switch tx := tx.(type) {
		case *transaction.TransferTx:
			return handlerTransferTx(ctx, am, tx)
		// todo

		default:
			errMsg := "Unrecognized Tx type: " + reflect.TypeOf(tx).Name()
			return types.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handlerTransferTx(ctx types.Context, am account.AccountMapper, tx *transaction.TransferTx) types.Result {
	if err := am.Transfer(ctx, tx.Common.From, tx.Amount, tx.To); err != nil {
		return transaction.ErrFailTransfer(transaction.DefaultCodespace, err.Error()).Result()
	}
	return types.Result{}
}