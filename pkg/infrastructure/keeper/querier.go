package keeper

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/infrastructure/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(k InfrastructureKeeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryContent:
			return QueryContent(ctx, req,  k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown nameservice query endpoint")
		}
	}
}


func QueryContent(ctx sdk.Context, req abci.RequestQuery, k InfrastructureKeeper) ([]byte, sdk.Error) {
	var params types.QueryContentParams
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, types.ErrCdcUnmarshalFailed(types.DefaultCodespace, err)
	}
	res, err := k.GetContent(ctx, params.Key)
	if err != nil {
		return nil, types.ErrGetInvalidResponse(types.DefaultCodespace, err)
	}
	return res, nil
}