package keeper

import (
	"errors"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
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
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error)  {
		switch path[0] {
		case types.QueryAccount:
			return queryAccount(ctx, req, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown nameservice query endpoint")
		}
	}
}

func queryAccount(ctx sdk.Context, req abci.RequestQuery, k AccountKeeper) ([]byte, sdk.Error) {
	var accountParams QueryAccountParams
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &accountParams)
	if err != nil {
		return nil, types.ErrGetAccount(types.DefaultCodespace, errors.New("unmarshal json failed"))
	}
	acc := k.GetAccount(ctx, accountParams.AccountAddress)
	if acc == nil {
		return nil, types.ErrGetAccount(types.DefaultCodespace, errors.New("account not found"))
	}
	by := types.ModuleCdc.MustMarshalBinaryLengthPrefixed(acc)
	return by, nil
}