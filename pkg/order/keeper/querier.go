package keeper

import (
	"encoding/json"
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	//"github.com/tanhuiya/ci123chain/pkg/order/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	QueryState = "shardState"
)

func NewQuerier(orderKeeper OrderKeeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryState:
			return queryState(ctx,orderKeeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown nameservice query endpoint")
		}
	}
}


func queryState(ctx sdk.Context, k OrderKeeper) ([]byte, sdk.Error) {

	_, orderBook := k.GetOrderBook()
	obByte, err := json.Marshal(orderBook)
	if err != nil {
		//return nil, types.ErrFailedMarshal(types.DefaultCodespace, err.Error())
	}
	return obByte , nil
}