package keeper

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	types2 "github.com/ci123chain/ci123chain/pkg/pre_staking/types"
	sktypes "github.com/ci123chain/ci123chain/pkg/staking/types"
	"github.com/tendermint/tendermint/abci/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(k PreStakingKeeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req types.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case types2.PreStakingRecordQuery:
			return PreStakingRecord(ctx, req, k)
		case types2.StakingRecordQuery:
			return StakingRecord(ctx, req, k)
		default:
			return nil, nil
		}
	}
}


func PreStakingRecord(ctx sdk.Context, req abci.RequestQuery, k PreStakingKeeper) ([]byte, error) {
	var params types2.QueryPreStakingRecord
	err := types2.PreStakingCodec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, err
	}
	res := k.GetAccountPreStaking(ctx, params.Delegator)
	return types2.PreStakingCodec.MarshalJSON(types2.QueryPreStakingResult{Amount:res.String()})
}

func StakingRecord(ctx sdk.Context, req abci.RequestQuery, k PreStakingKeeper) ([]byte, error) {
	var params types2.QueryStakingRecord
	err := types2.PreStakingCodec.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, err
	}
	delegation, found := k.StakingKeeper.GetDelegation(ctx, params.DelegatorAddr, params.ValidatorAddr)
	if !found {
		return nil, nil
	}
	val, found := k.StakingKeeper.GetValidator(ctx, delegation.ValidatorAddress)
	if !found {
		return nil, nil
	}
	res := sktypes.NewDelegationResp(delegation.DelegatorAddress,
		delegation.ValidatorAddress,
		delegation.GetShares(),
		sdk.NewChainCoin(val.TokensFromShares(delegation.Shares).TruncateInt()),)
	return types2.PreStakingCodec.MarshalJSON(res)
}