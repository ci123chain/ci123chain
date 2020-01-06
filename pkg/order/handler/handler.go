package handler

import (
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
	order "github.com/tanhuiya/ci123chain/pkg/order"
	"github.com/tanhuiya/ci123chain/pkg/order/keeper"
	"reflect"
)

func NewHandler(keeper *keeper.OrderKeeper) types.Handler {
	return func(ctx types.Context, tx types.Tx) types.Result {
		switch tx := tx.(type) {
		case *order.UpgradeTx:
			return handlerUpgradeTx(ctx, keeper, tx)
		default:
			errMsg := "Unrecognized Tx type: " + reflect.TypeOf(tx).Name()
			return types.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handlerUpgradeTx(ctx types.Context,k *keeper.OrderKeeper, tx *order.UpgradeTx) types.Result {
	///扩展容量交易的处理
	var orderbook keeper.OrderBook
	store := ctx.KVStore(k.StoreKey)
	bz := store.Get([]byte(keeper.OrderBookKey))
	err := keeper.ModuleCdc.UnmarshalBinaryLengthPrefixed(bz, &orderbook)
	if err != nil {
		panic(err)
	}
	//现在是新添加一个分片
	var action keeper.Actions
	action.Name = tx.Name
	action.Height = tx.Height
	action.Type = tx.Type

	k.UpdateOrderBook(ctx, orderbook, &action)

	return types.Result{}
}