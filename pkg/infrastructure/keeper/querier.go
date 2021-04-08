package keeper

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/ci123chain/ci123chain/pkg/infrastructure/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(k InfrastructureKeeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case types.QueryContent:
			return QueryContent(ctx, req,  k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown query endpoint")
		}
	}
}


func QueryContent(ctx sdk.Context, req abci.RequestQuery, k InfrastructureKeeper) ([]byte, error) {
	var params types.QueryContentParams
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInternal, "cdc unmarshal failed")
	}
	res, err := k.GetContent(ctx, params.Key)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInternal, err.Error())
	}
	return res, nil
}