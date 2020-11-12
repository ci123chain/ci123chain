package keeper

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/mint/types"
	abci "github.com/tendermint/tendermint/abci/types"
)


func NewQuerier(k MinterKeeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryAnnualProvisions:
			return queryAnnualProvisions(ctx, k)
		case types.QueryInflation:
			return QueryInflation(ctx, k)
		case types.QueryParameters:
			return QueryParams(ctx, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown query endpoint")
		}
	}
}


func QueryParams(ctx sdk.Context,k MinterKeeper) ([]byte, sdk.Error) {

	params := k.GetParams(ctx)

	res := types.MintCdc.MustMarshalJSON(params)
	return res, nil
}


func QueryInflation(ctx sdk.Context,k MinterKeeper) ([]byte, sdk.Error) {

	minter := k.GetMinter(ctx)
	res := types.MintCdc.MustMarshalJSON(minter.Inflation)

	return res, nil
}


func queryAnnualProvisions(ctx sdk.Context,k MinterKeeper) ([]byte, sdk.Error) {
	minter := k.GetMinter(ctx)
	res := types.MintCdc.MustMarshalJSON(minter.AnnualProvisions)

	return res, nil
}