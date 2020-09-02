package keeper

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/distribution/types"
)

func (k DistrKeeper) GetParams(ctx sdk.Context) (params types.Params){
	k.ParamSpace.GetParamSet(ctx, &params)
	return params
}

func (k DistrKeeper) SetParams(ctx sdk.Context, params types.Params) {
	k.ParamSpace.SetParamSet(ctx, &params)
}

// GetBaseProposerReward returns the current distribution base proposer rate.
func (k DistrKeeper) GetBaseProposerReward(ctx sdk.Context) (percent sdk.Dec) {
	k.ParamSpace.Get(ctx, types.ParamStoreKeyBaseProposerReward, &percent)
	return percent
}


// GetBonusProposerReward returns the current distribution bonus proposer reward
// rate.
func (k DistrKeeper) GetBonusProposerReward(ctx sdk.Context) (percent sdk.Dec) {
	k.ParamSpace.Get(ctx, types.ParamStoreKeyBonusProposerReward, &percent)
	return percent
}

func (k DistrKeeper) GetCommunityTax(ctx sdk.Context) (percent sdk.Dec) {
	k.ParamSpace.Get(ctx, types.ParamStoreKeyCommunityTax, &percent)
	return percent
}

// GetWithdrawAddrEnabled returns the current distribution withdraw address
// enabled parameter.
func (k DistrKeeper) GetWithdrawAddrEnabled(ctx sdk.Context) (enabled bool) {
	k.ParamSpace.Get(ctx, types.ParamStoreKeyWithdrawAddrEnabled, &enabled)
	return enabled
}