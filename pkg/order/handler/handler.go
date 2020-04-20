package handler

import (
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/order/keeper"
	order "github.com/ci123chain/ci123chain/pkg/order/types"
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

	orderbook, err := k.GetOrderBook(ctx)
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