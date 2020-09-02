package keeper

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/distribution/types"
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
		case types.QueryValidatorOutstandingRewards:
			return queryValidatorOutstandingRewards(ctx, req, keeper)
		case types.QueryCommunityPool:
			return queryCommunityPool(ctx, req, keeper)
		case types.QueryWithdrawAddress:
			return queryDelegatorWithdrawAddress(ctx, req, keeper)
		case types.QueryValidatorCommission:
			return queryValidatorCommission(ctx, req, keeper)
		case types.QueryDelegationRewards:
			return queryDelegationRewards(ctx, req, keeper)

		default:
			return nil, sdk.ErrUnknownRequest("unknown nameservice query endpoint")
		}
	}
}

func queryRewards(ctx sdk.Context, path []string, _ abci.RequestQuery, keeper DistrKeeper) ([]byte, sdk.Error){

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

func queryValidatorOutstandingRewards(ctx sdk.Context, req abci.RequestQuery, k DistrKeeper) ([]byte, sdk.Error) {
	var params types.QueryValidatorOutstandingRewardsParams
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, types.ErrFailedMarshal(types.DefaultCodespace, err.Error())
	}

	rewards := k.GetValidatorOutstandingRewards(ctx, params.ValidatorAddress)

	res := types.DistributionCdc.MustMarshalJSON(rewards)

	return res, nil
}

func queryCommunityPool(ctx sdk.Context, _ abci.RequestQuery, k DistrKeeper) ([]byte, sdk.Error) {

	pool := k.GetFeePoolCommunity(ctx)

	res := types.DistributionCdc.MustMarshalJSON(pool)
	return res, nil
}

func queryDelegatorWithdrawAddress(ctx sdk.Context, req abci.RequestQuery, k DistrKeeper) ([]byte, sdk.Error) {

	var params types.QueryDelegatorWithdrawAddrParams
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, types.ErrFailedMarshal(types.DefaultCodespace, err.Error())
	}
	withdrawAddr := k.GetDelegatorWithdrawAddr(ctx, params.DelegatorAddress)

	res := types.DistributionCdc.MustMarshalJSON(withdrawAddr)
	return res, nil
}

func queryValidatorCommission(ctx sdk.Context, req abci.RequestQuery, k DistrKeeper) ([]byte, sdk.Error) {
	var params types.QueryValidatorCommissionParams
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, types.ErrFailedMarshal(types.DefaultCodespace, err.Error())
	}
	commission := k.GetValidatorAccumulatedCommission(ctx, params.ValidatorAddress)
	res := types.DistributionCdc.MustMarshalJSON(commission)
	return res, nil
}

func queryDelegationRewards(ctx sdk.Context, req abci.RequestQuery, k DistrKeeper) ([]byte, sdk.Error) {
	var params types.QueryDelegationRewardsParams
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, types.ErrFailedMarshal(types.DefaultCodespace, err.Error())
	}

	// cache-wrap context as to not persist state changes during querying
	ctx, _ = ctx.CacheContext()

	val := k.StakingKeeper.Validator(ctx, params.ValidatorAddress)
	if val == nil {
		return nil, types.ErrNoValidatorExist(types.DefaultCodespace, params.ValidatorAddress.String())
	}

	del := k.StakingKeeper.Delegation(ctx, params.DelegatorAddress, params.ValidatorAddress)
	if del == nil {
		return nil, types.ErrNoDelegationExist(types.DefaultCodespace, params.ValidatorAddress.String(), params.DelegatorAddress.String())
	}

	endingPeriod := k.incrementValidatorPeriod(ctx, val)
	fmt.Printf("endingPeriod: %d", endingPeriod)
	rewards := k.calculateDelegationRewards(ctx, val, del, endingPeriod)
	res := types.DistributionCdc.MustMarshalJSON(rewards)
	return res, nil
}