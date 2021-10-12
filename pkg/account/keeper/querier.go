package keeper

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/account/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type QueryAccountParams struct {
	AccountAddress   sdk.AccAddress `json:"account_address"`
	Height           int64          `json:"height"`
}

func NewQueryAccountParams(accountAddress sdk.AccAddress, h int64) QueryAccountParams {
	params := QueryAccountParams{AccountAddress: accountAddress, Height: h}
	return params
}

func NewQuerier(k AccountKeeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error)  {
		switch path[0] {
		case types.QueryAccount:
			return queryAccount(ctx, req, k)
		case types.QueryAccountNonce:
			return queryAccountNonce(ctx, req, k)
		default:
			return nil, types.ErrInvalidEndPoint
		}
	}
}

func queryAccount(ctx sdk.Context, req abci.RequestQuery, k AccountKeeper) ([]byte, error) {
	var accountParams QueryAccountParams
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &accountParams)
	if err != nil {
		return nil, types.ErrCdcUnmarshalFailed
	}
	acc := k.GetAccount(ctx, accountParams.AccountAddress)
	if acc == nil {
		return nil, types.ErrAccountNotExisted
	}
	if acc.GetIsModule() {
		acc = types.NewBaseAccountFromExportAccount(acc)
	}
	by := types.ModuleCdc.MustMarshalBinaryLengthPrefixed(acc)
	return by, nil
}

func queryAccountNonce(ctx sdk.Context, req abci.RequestQuery, k AccountKeeper) ([]byte, error) {
	var accountParams QueryAccountParams
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &accountParams)
	if err != nil {
		return nil, types.ErrCdcUnmarshalFailed
	}
	acc := k.GetAccount(ctx, accountParams.AccountAddress)
	if acc == nil {
		return nil, types.ErrAccountNotExisted
	}
	if acc.GetIsModule() {
		acc = types.NewBaseAccountFromExportAccount(acc)
	}
	by := types.ModuleCdc.MustMarshalBinaryLengthPrefixed(acc.GetSequence())
	return by, nil
}