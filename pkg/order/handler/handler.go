package handler

import (
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
	order "github.com/tanhuiya/ci123chain/pkg/order"
	"github.com/tanhuiya/ci123chain/pkg/order/keeper"
	"reflect"
)

func NewHandler(keeper keeper.OrderKeeper) types.Handler {
	return func(ctx types.Context, tx types.Tx) types.Result {
		switch tx := tx.(type) {
		case *order.UpgradeTx:
			return handlerUpgradeTx(ctx,keeper, tx)
		default:
			errMsg := "Unrecognized Tx type: " + reflect.TypeOf(tx).Name()
			return types.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handlerUpgradeTx(ctx types.Context,k keeper.OrderKeeper, tx *order.UpgradeTx) types.Result {
	///扩展容量交易的处理
	_, orderBook := k.GetOrderBook()
	//现在是新添加一个分片
	var action keeper.Actions
	action.Name = tx.Name
	action.Height = tx.Height
	action.Type = tx.Type
	orderBook.Actions = append(orderBook.Actions, action)

	k.SetEventBook(orderBook)

	return types.Result{}
}