package keeper

import (
	"encoding/json"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/order/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(orderKeeper *OrderKeeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case types.QueryState:
			return queryState(ctx, orderKeeper)
		default:
			return nil, types.ErrInvalidEndPoint
		}
	}
}

func queryState(ctx sdk.Context, k *OrderKeeper) ([]byte, error) {
	order, err := k.GetOrderBook(ctx)
	if err != nil {
		return nil, types.ErrGetOrderBookFailed
	}
	obytes, err := json.Marshal(order)
	if err != nil {
		return nil, types.ErrFailedMarshal
	}
	return obytes, nil
}