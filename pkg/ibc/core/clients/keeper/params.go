package keeper

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/clients/types"
)

func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(k.getAllowedClients(ctx)...)
}

func (k Keeper) getAllowedClients(ctx sdk.Context) []string {
	var res []string
	k.paramSpace.Get(ctx, types.KeyAllowedClients, &res)
	return res
}