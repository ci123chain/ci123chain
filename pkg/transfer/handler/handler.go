package handler

import (
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/account/keeper"
	"github.com/tanhuiya/ci123chain/pkg/db"
	n "github.com/tanhuiya/ci123chain/pkg/nonce"
	"github.com/tanhuiya/ci123chain/pkg/transaction"
	"github.com/tanhuiya/ci123chain/pkg/transfer"
	"reflect"
)

func NewHandler(
	txm transaction.TxIndexMapper,
	am keeper.AccountKeeper,
	sm *db.StateManager) types.Handler {
	return func(ctx types.Context, tx types.Tx) types.Result{
		ctx = ctx.WithTxIndex(txm.Get(ctx))
		defer func() {
			txm.Incr(ctx)
		}()
		switch tx := tx.(type) {
		case *transfer.TransferTx:
			return handlerTransferTx(ctx, am, tx)
		// todo

		default:
			errMsg := "Unrecognized Tx type: " + reflect.TypeOf(tx).Name()
			return types.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handlerTransferTx(ctx types.Context, am keeper.AccountKeeper, tx *transfer.TransferTx) types.Result {
	//----------
	//对transaction.CommonTx的处理
	savedSequence := am.GetAccount(ctx, tx.Common.From).GetSequence()
	checkResult := n.CheckTransferNonce(ctx, savedSequence, tx.Common.Nonce)
	if checkResult != true {
		return types.ErrInvalidSequence("Unexpected nonce").Result()
	}
	//----------
	if err := am.Transfer(ctx, tx.Common.From, tx.To, tx.Amount); err != nil {
		return err.Result()
	}

	//交易成功，nonce+1
	saveErr := am.GetAccount(ctx, tx.Common.From).SetSequence(tx.Common.Nonce + 1)
	if saveErr != nil {
		return types.ErrInvalidSequence("Unexpected nonce of transaction").Result()
	}
	//
	return types.Result{}
}