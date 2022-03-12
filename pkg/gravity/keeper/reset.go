package keeper

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/gravity/types"
)

// reset gravity data
func (k Keeper) Reset(ctx sdk.Context)  {
	st := ctx.KVStore(k.storeKey)
	iterator := st.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {
		st.Delete(iterator.Key())
	}
	gs := types.DefaultGenesisState()
	InitGenesis(ctx, k, *gs)
}
