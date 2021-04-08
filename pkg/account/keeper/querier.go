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
}

func NewQueryAccountParams(accountAddress sdk.AccAddress) QueryAccountParams {
	params := QueryAccountParams{AccountAddress: accountAddress}
	return params
}

func NewQuerier(k AccountKeeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error)  {
		switch path[0] {
		case types.QueryAccount:
			return queryAccount(ctx, req, k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown query endpoint")
		}
	}
}

func queryAccount(ctx sdk.Context, req abci.RequestQuery, k AccountKeeper) ([]byte, error) {
	var accountParams QueryAccountParams
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &accountParams)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("cdc marshal failed: %v", err.Error()))
	}
	acc := k.GetAccount(ctx, accountParams.AccountAddress)
	if acc == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("get account faield"))
	}
	by := types.ModuleCdc.MustMarshalBinaryLengthPrefixed(acc)
	return by, nil
}