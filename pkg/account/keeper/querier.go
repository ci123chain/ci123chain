package keeper

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/ci123chain/ci123chain/pkg/account/types"
	"github.com/ci123chain/ci123chain/pkg/util"
	abci "github.com/tendermint/tendermint/abci/types"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	"github.com/tendermint/tendermint/rpc/client/http"
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
	var height int64
	if accountParams.Height == -1 {
		height = ctx.BlockHeight() -1
	}else if accountParams.Height < 1 || accountParams.Height > ctx.BlockHeight() {
		return nil, sdkerrors.Wrap(sdkerrors.ErrParams, fmt.Sprintf("unexpected height: %v", accountParams.Height))
	}else {
		height = accountParams.Height
	}
	i := SearchHeight(ctx, k, acc, height)
	if i ==-2 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrParams, fmt.Sprintf("unexpected height: %v", accountParams.Height))
	}else if i ==-3 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInternal, "cdc unmarshal failed")
	}else {
		r := GetHistoryBalance(ctx, k, acc, i)
		key := types.AddressStoreKey(accountParams.AccountAddress)
		if r.Shard == ctx.ChainID() {
			queryable := ctx.MultiStore().(sdk.Queryable)
			req := abci.RequestQuery{
				Path:   fmt.Sprintf("/%s/key", types.StoreKey),
				Height: height,
				Data:   key,
				Prove:  req.Prove,
			}
			res := queryable.Query(req)
			ha := util.NewHistoryAccount("", res.Value, res.ProofOps)
			by := types.ModuleCdc.MustMarshalBinaryLengthPrefixed(ha)
			return by, nil
		}else {
			return sendRequest(ctx, types.StoreKey, "tcp://" + r.Shard + ":26657", key, height, req.Prove)
		}
	}
}

func sendRequest(ctx sdk.Context, storekey, host string, key []byte, height int64, prove bool) ([]byte, error) {
	rpc, err := http.New(host, "/websocket")
	if err != nil {
		return nil, err
	}
	req := abci.RequestQuery{
		Path:   fmt.Sprintf("store/%s/key", storekey),
		Height: height,
		Data:   key,
		Prove:  prove,
	}
	opts := rpcclient.ABCIQueryOptions{
		Height: height,
		Prove:  req.Prove,
	}
	res, err := rpc.ABCIQueryWithOptions(ctx, req.Path, req.Data, opts)
	if err != nil {
		return nil, err
	}
	ha := util.NewHistoryAccount("", res.Response.Value, res.Response.ProofOps)
	by := types.ModuleCdc.MustMarshalBinaryLengthPrefixed(ha)
	return by, nil
}