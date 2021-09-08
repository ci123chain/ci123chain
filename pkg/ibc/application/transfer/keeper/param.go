package keeper

import (
	"bytes"
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

func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	ps := types.DefaultParams()
	for _, v := range ps.ParamSetPairs() {
		if bytes.Equal(v.Key, types.KeySendEnabled) {
			params.SendEnabled = v.Value.(bool)
		}else if bytes.Equal(v.Key, types.KeyReceiveEnabled) {
			params.ReceiveEnabled = v.Value.(bool)
		}
	}
	return
}
