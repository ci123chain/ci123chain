package keeper

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/application/transfer/types"
)

// GetSendEnabled retrieves the send enabled boolean from the paramstore
func (k Keeper) GetSendEnabled(ctx sdk.Context) bool {
	var res bool
	k.paramSpace.Get(ctx, types.KeySendEnabled, &res)
	return res
}

// GetReceiveEnabled retrieves the receive enabled boolean from the paramstore
func (k Keeper) GetReceiveEnabled(ctx sdk.Context) bool {
	var res bool
	k.paramSpace.Get(ctx, types.KeyReceiveEnabled, &res)
	return res
}

// SetParams sets the total set of ibc-transfer parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}
