package keeper

import (
	"encoding/json"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	types2 "github.com/ci123chain/ci123chain/pkg/pre_staking/types"
	"github.com/tendermint/tendermint/abci/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(k PreStakingKeeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req types.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case types2.StakingRecordQuery:
			return StakingRecord(ctx, req, k)
		case types2.PreStakingTokenQuery:
			return PreStakingToken(ctx, req, k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown request endpoint")
		}
	}
}

func PreStakingToken(ctx sdk.Context, req abci.RequestQuery, k PreStakingKeeper) ([]byte, error) {
	res := k.GetTokenManager(ctx)
	by, err := json.Marshal(res)
	return by, err
}

func StakingRecord(ctx sdk.Context, req abci.RequestQuery, k PreStakingKeeper) ([]byte, error) {
	var params types2.QueryStakingRecord
	err := k.Cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, err
	}
	records := k.GetAllStakingVault(ctx)
	var results []types2.StakingVault
	for _, r := range records{
		if r.Delegator == params.DelegatorAddr {
			results = append(results, r)
		}
	}
	return k.Cdc.MarshalBinaryBare(results)
}