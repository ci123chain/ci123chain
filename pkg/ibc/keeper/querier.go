package keeper

import (
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/ibc/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the nameservice Querier
const (
	QueryState = "state"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper IBCKeeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryState:
			return queryResolve(ctx, path[1:], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown nameservice query endpoint")
		}
	}
}

// nolint: unparam
func queryResolve(ctx sdk.Context, path []string, req abci.RequestQuery, keeper IBCKeeper) ([]byte, sdk.Error) {
	if path[0] != StateReady {
		return nil, sdk.ErrUnknownRequest("Parameter State Error")
	}

	value := *keeper.GetFirstReadyIBCMsg(ctx)

	retbz, err := types.IbcCdc.MarshalBinaryLengthPrefixed(value)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(err.Error())
	}
	return retbz, nil
}