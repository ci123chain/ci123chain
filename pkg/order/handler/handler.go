package handler

import (
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/order/keeper"
	order "github.com/ci123chain/ci123chain/pkg/order/types"
)

func NewHandler(keeper *keeper.OrderKeeper) types.Handler {
	return func(ctx types.Context, msg types.Msg) (*types.Result, error) {
		ctx = ctx.WithEventManager(types.NewEventManager())
		switch msg := msg.(type) {
		case *order.MsgUpgrade:
			return handlerMsgUpgrade(ctx, keeper, msg)
		default:
			//errMsg := "Unrecognized msg type: " + reflect.TypeOf(msg).Name()
			return nil, order.ErrInvalidEndPoint
		}
	}
}

func handlerMsgUpgrade(ctx types.Context,k *keeper.OrderKeeper, msg *order.MsgUpgrade) (*types.Result, error) {
	///扩展容量交易的处理

	orderbook, err := k.GetOrderBook(ctx)
	if err != nil {
		panic(err)
	}

	//现在是新添加一个分片
	var action order.Actions
	action.Name = msg.Name
	action.Height = msg.Height
	action.Type = msg.Type

	k.UpdateOrderBook(ctx, orderbook, &action)

	em := ctx.EventManager()
	em.EmitEvents(types.Events{
		types.NewEvent(order.EventType,
			types.NewAttribute([]byte(types.AttributeKeyMethod), []byte(order.AttributeValueAddShard)),
			types.NewAttribute([]byte(types.AttributeKeyModule), []byte(order.AttributeValueCategory)),
			types.NewAttribute([]byte(types.AttributeKeySender), []byte(msg.FromAddress.String())),
		),
	})
	return &types.Result{Events: em.Events()}, nil
}