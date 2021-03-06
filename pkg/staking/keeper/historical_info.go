package keeper

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/staking/types"
)

func (k StakingKeeper) GetHistoricalInfo(ctx sdk.Context, height int64) (types.HistoricalInfo, bool) {
	//
	var hi types.HistoricalInfo
	store := ctx.KVStore(k.storeKey)
	key := types.GetHistoricalInfoKey(height)

	value := store.Get(key)
	if value == nil {
		return types.HistoricalInfo{}, false
	}

	types.StakingCodec.MustUnmarshalBinaryLengthPrefixed(value, &hi)
	return hi, true
}

// SetHistoricalInfo sets the historical info at a given height
func (k StakingKeeper) SetHistoricalInfo(ctx sdk.Context, height int64, hi types.HistoricalInfo) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetHistoricalInfoKey(height)

	value := types.StakingCodec.MustMarshalBinaryLengthPrefixed(hi)
	store.Set(key, value)
}

// DeleteHistoricalInfo deletes the historical info at a given height
func (k StakingKeeper) DeleteHistoricalInfo(ctx sdk.Context, height int64) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetHistoricalInfoKey(height)

	store.Delete(key)
}

func (k StakingKeeper) TrackHistoricalInfo(ctx sdk.Context) {

	entryNum := k.HistoricalEntries(ctx)

	// Prune store to ensure we only have parameter-defined historical entries.
	// In most cases, this will involve removing a single historical entry.
	// In the rare scenario when the historical entries gets reduced to a lower value k'
	// from the original value k. k - k' entries must be deleted from the store.
	// Since the entries to be deleted are always in a continuous range, we can iterate
	// over the historical entries starting from the most recent version to be pruned
	// and then return at the first empty entry.
	for i := ctx.BlockHeight() - int64(entryNum); i >= 0; i-- {
		_, found := k.GetHistoricalInfo(ctx, i)
		if found {
			k.DeleteHistoricalInfo(ctx, i)
		} else {
			break
		}
	}

	// if there is no need to persist historicalInfo, return
	if entryNum == 0 {
		return
	}

	// Create HistoricalInfo struct
	lastVals := k.GetLastValidators(ctx)
	historicalEntry := types.NewHistoricalInfo(ctx.BlockHeader(), lastVals)

	// Set latest HistoricalInfo at current height
	k.SetHistoricalInfo(ctx, ctx.BlockHeight(), historicalEntry)
}