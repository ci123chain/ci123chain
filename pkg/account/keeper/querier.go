package keeper

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
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
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown query endpoint")
		}
	}
}

func queryAccountNonce(ctx sdk.Context, req abci.RequestQuery, k AccountKeeper) ([]byte, error) {
	var accountParams QueryAccountParams
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &accountParams)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("cdc marshal failed: %v", err.Error()))
	}
	acc := k.GetAccount(ctx, accountParams.AccountAddress)
	if acc == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("get account faield"))
	}
	by := types.ModuleCdc.MustMarshalBinaryLengthPrefixed(acc.GetSequence())
	return by, nil
}

func queryAccount(ctx sdk.Context, req abci.RequestQuery, k AccountKeeper) ([]byte, error) {

	var accountParams QueryAccountParams
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &accountParams)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("cdc marshal failed: %v", err.Error()))
	}
	acc := k.GetAccount(ctx, accountParams.AccountAddress)
	if acc == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("the account not exist"))
	}
	by := types.ModuleCdc.MustMarshalBinaryLengthPrefixed(acc)
	return by, nil
}