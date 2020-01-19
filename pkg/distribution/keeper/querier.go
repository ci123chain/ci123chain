package keeper

import (
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/distribution/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"strconv"
)

const (
	QueryRewards = "rewards"
)

func NewQuerier(keeper DistrKeeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryRewards:
			return queryRewards(ctx, path[1:], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown nameservice query endpoint")
		}
	}
}

func queryRewards(ctx sdk.Context, path []string, req abci.RequestQuery, keeper DistrKeeper) ([]byte, sdk.Error){

	accountAddress := path[0]
	height := path[1]
	if height == "now" {
		h := ctx.BlockHeight()
		height = strconv.FormatInt(h, 10)
	}else {
		_, Err := strconv.ParseInt(height, 10, 64)
		if Err != nil {
			return nil, types.ErrBadHeight(types.DefaultCodespace, Err)
		}
	}

	key := accountAddress + height
	address := []byte(key)
	addr := sdk.AccAddr(address)
	rewards, err := keeper.GetValCurrentRewards(ctx, addr)
	if err != nil {
		return nil, types.ErrBadHeight(types.DefaultCodespace, err)
	}

	amount := uint64(rewards.Amount.Int64())
	retbz, err := types.DistributionCdc.MarshalBinaryLengthPrefixed(amount)
	if err != nil {
		return nil, types.ErrFailedMarshal(types.DefaultCodespace, err.Error())
	}
	return retbz, nil
}
