package keeper

import (
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	"github.com/ci123chain/ci123chain/pkg/abci/store"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	capabilitykeeper "github.com/ci123chain/ci123chain/pkg/capability/keeper"
	"github.com/ci123chain/ci123chain/pkg/ibc/application/transfer/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/host"
	paramtypes "github.com/ci123chain/ci123chain/pkg/params/subspace"
	supplytypes "github.com/ci123chain/ci123chain/pkg/supply/exported"
	cmn "github.com/tendermint/tendermint/libs/common"

	"github.com/tendermint/tendermint/libs/log"
)

// Keeper defines the IBC fungible transfer keeper
type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        *codec.Codec
	paramSpace paramtypes.Subspace

	channelKeeper types.ChannelKeeper
	portKeeper    types.PortKeeper
	authKeeper    types.AccountKeeper
	bankKeeper    types.BankKeeper
	scopedKeeper  capabilitykeeper.ScopedKeeper
}


// NewKeeper creates a new IBC transfer Keeper instance
func NewKeeper(
	cdc *codec.Codec, key sdk.StoreKey, paramSpace paramtypes.Subspace,
	channelKeeper types.ChannelKeeper, portKeeper types.PortKeeper,
	authKeeper types.AccountKeeper, bankKeeper types.BankKeeper, scopedKeeper capabilitykeeper.ScopedKeeper,
) Keeper {

	// ensure ibc transfer module account is set
	if addr := authKeeper.GetModuleAddress(types.ModuleName); addr.Empty() {
		panic("the IBC transfer module account has not been set")
	}

	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:           cdc,
		storeKey:      key,
		paramSpace:    paramSpace,
		channelKeeper: channelKeeper,
		portKeeper:    portKeeper,
		authKeeper:    authKeeper,
		bankKeeper:    bankKeeper,
		scopedKeeper:  scopedKeeper,
	}
}


// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+host.ModuleName+"-"+types.ModuleName)
}

// GetTransferAccount returns the ICS20 - transfers ModuleAccount
func (k Keeper) GetTransferAccount(ctx sdk.Context) supplytypes.ModuleAccountI {
	return k.authKeeper.GetModuleAccount(ctx, types.ModuleName)
}

// GetDenomTrace retreives the full identifiers trace and base denomination from the store.
func (k Keeper) GetDenomTrace(ctx sdk.Context, denomTraceHash cmn.HexBytes) (types.DenomTrace, bool) {
	prefiSstore := store.NewPrefixStore(ctx.KVStore(k.storeKey), types.DenomTraceKey)
	bz := prefiSstore.Get(denomTraceHash)
	if bz == nil {
		return types.DenomTrace{}, false
	}

	denomTrace := k.MustUnmarshalDenomTrace(bz)
	return denomTrace, true
}


// HasDenomTrace checks if a the key with the given denomination trace hash exists on the store.
func (k Keeper) HasDenomTrace(ctx sdk.Context, denomTraceHash tmbytes.HexBytes) bool {
	store := store.NewPrefixStore(ctx.KVStore(k.storeKey), types.DenomTraceKey)
	return store.Has(denomTraceHash)
}

// SetDenomTrace sets a new {trace hash -> denom trace} pair to the store.
func (k Keeper) SetDenomTrace(ctx sdk.Context, denomTrace types.DenomTrace) {
	store := store.NewPrefixStore(ctx.KVStore(k.storeKey), types.DenomTraceKey)
	bz := k.MustMarshalDenomTrace(denomTrace)
	store.Set(denomTrace.Hash(), bz)
}