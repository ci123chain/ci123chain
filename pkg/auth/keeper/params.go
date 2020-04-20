package keeper

import (
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	auth_types "github.com/ci123chain/ci123chain/pkg/auth/types"
)

func (ak AuthKeeper) SetParams(ctx types.Context, params auth_types.Params) {
	ak.paramSubspace.SetParamSet(ctx, &params)
}

// GetParams gets the auth module's parameters.
func (ak AuthKeeper) GetParams(ctx types.Context) (params auth_types.Params) {
	ak.paramSubspace.GetParamSet(ctx, &params)
	return
}
