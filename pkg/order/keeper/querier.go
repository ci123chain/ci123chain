package keeper

import (
	"encoding/json"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/order/types"
	abci "github.com/tendermint/tendermint/abci/types"
)


func NewQuerier(orderKeeper *OrderKeeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryState:
			return queryState(ctx, orderKeeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown nameservice query endpoint")
		}
	}
}


func queryState(ctx sdk.Context, k *OrderKeeper) ([]byte, sdk.Error) {
	order, err := k.GetOrderBook(ctx)
	if err != nil {
		panic(err)
	}
	obytes, err := json.Marshal(order)
	return obytes, nil
}