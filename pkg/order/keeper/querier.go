package keeper

import (
	"encoding/json"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/ci123chain/ci123chain/pkg/order/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(orderKeeper *OrderKeeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case types.QueryState:
			return queryState(ctx, orderKeeper)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown query endpoint")
		}
	}
}

func queryState(ctx sdk.Context, k *OrderKeeper) ([]byte, error) {
	order, err := k.GetOrderBook(ctx)
	if err != nil {
		panic(err)
	}
	obytes, err := json.Marshal(order)
	return obytes, nil
}