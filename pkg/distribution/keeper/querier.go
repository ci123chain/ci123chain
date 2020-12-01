package keeper

import (
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
		case types.QueryAccountInfo:
			return queryDelegatorAccountInfo(ctx, req, keeper)

		default:
			return nil, sdk.ErrUnknownRequest("unknown query endpoint")
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
	rewards := k.calculateDelegationRewards(ctx, val, del, endingPeriod)
	res := types.DistributionCdc.MustMarshalJSON(rewards)
	return res, nil
}

func queryDelegatorAccountInfo(ctx sdk.Context, req abci.RequestQuery, k DistrKeeper) ([]byte, sdk.Error) {
	var params types.QueryDelegatorBalanceParams

	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal("unmarshal failed")
	}

	ctx, _ = ctx.CacheContext()

	balance := k.AccountKeeper.GetBalance(ctx, params.AccountAddress)

	validator, found := k.StakingKeeper.GetValidator(ctx, params.AccountAddress)
	if !found {
		zero := sdk.NewEmptyCoin()
		result := types.NewDelegatorAccountInfo(balance, zero, zero, zero, zero)
		res := types.DistributionCdc.MustMarshalJSON(result)
		return res, nil
	}

	cs := k.GetValidatorAccumulatedCommission(ctx, params.AccountAddress).Commission
	commission := sdk.NewCoin(cs.Amount.RoundInt())

	del := k.StakingKeeper.Delegation(ctx, params.AccountAddress, params.AccountAddress)
	if del == nil {
		return nil, types.ErrNoDelegationExist(types.DefaultCodespace, params.AccountAddress.String(), params.AccountAddress.String())
	}

	endingPeriod := k.incrementValidatorPeriod(ctx, validator)
	rw := k.calculateDelegationRewards(ctx, validator, del, endingPeriod)
	rewards := sdk.NewCoin(rw.Amount.RoundInt())

	ctxTime := ctx.BlockHeader().Time
	matureUnbonds := k.StakingKeeper.DequeueAllMatureUBDQueue(ctx, ctxTime)
	unbondings := sdk.NewEmptyCoin()//all unboding balance
	delegated := sdk.NewEmptyCoin()//all delegated
	for _, dvPair := range matureUnbonds {
		if dvPair.DelegatorAddress == params.AccountAddress {
			ubd, found := k.StakingKeeper.GetUnbondingDelegation(ctx, dvPair.DelegatorAddress, dvPair.ValidatorAddress)
			if !found {
				continue
			}
			for i := 0; i < len(ubd.Entries); i++ {
				entry := ubd.Entries[i]
				if entry.IsMature(ctxTime) {
					if !entry.Balance.IsZero() {
						amt := sdk.NewCoin(entry.Balance)
						unbondings = unbondings.Add(amt)
					}
				}
				i--
			}
		}
	}

	delegations := k.StakingKeeper.GetAllDelegatorDelegations(ctx, params.AccountAddress)
	for _, v := range delegations {
		amt := sdk.NewCoin(v.GetShares().RoundInt())
		delegated = delegated.Add(amt)
	}

	result := types.NewDelegatorAccountInfo(balance, delegated, unbondings, rewards, commission)
	res := types.DistributionCdc.MustMarshalJSON(result)
	return res, nil
}