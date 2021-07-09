package keeper

import (
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	"github.com/ci123chain/ci123chain/pkg/abci/store"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	capabilitykeeper "github.com/ci123chain/ci123chain/pkg/capability/keeper"
	capabilitytypes "github.com/ci123chain/ci123chain/pkg/capability/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/application/transfer/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/host"
	"github.com/ci123chain/ci123chain/pkg/params"
	supplytypes "github.com/ci123chain/ci123chain/pkg/supply/exported"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	"github.com/tendermint/tendermint/libs/log"
)

// Keeper defines the IBC fungible transfer keeper
type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        codec.BinaryMarshaler
	paramSpace params.Subspace

	channelKeeper types.ChannelKeeper
	portKeeper    types.PortKeeper
	supplyKeepr    types.SupplyKeeper
	scopedKeeper  capabilitykeeper.ScopedKeeper
}


// NewKeeper creates a new IBC transfer Keeper instance
func NewKeeper(
	cdc codec.BinaryMarshaler, key sdk.StoreKey, paramSpace params.Subspace,
	channelKeeper types.ChannelKeeper, portKeeper types.PortKeeper,
	supplyKeeper types.SupplyKeeper, scopedKeeper capabilitykeeper.ScopedKeeper,
) Keeper {

	// ensure ibc transfer module account is set
	if addr := supplyKeeper.GetModuleAddress(types.ModuleName); addr.Empty() {
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
		supplyKeepr:   supplyKeeper,
		scopedKeeper:  scopedKeeper,
	}
}


// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+host.ModuleName+"-"+types.ModuleName)
}

// GetTransferAccount returns the ICS20 - transfers ModuleAccount
func (k Keeper) GetTransferAccount(ctx sdk.Context) supplytypes.ModuleAccountI {
	return k.supplyKeepr.GetModuleAccount(ctx, types.ModuleName)
}

// GetDenomTrace retreives the full identifiers trace and base denomination from the store.
func (k Keeper) GetDenomTrace(ctx sdk.Context, denomTraceHash tmbytes.HexBytes) (types.DenomTrace, bool) {
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

// ClaimCapability allows the transfer module that can claim a capability that IBC module
// passes to it
func (k Keeper) ClaimCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) error {
	return k.scopedKeeper.ClaimCapability(ctx, cap, name)
}


// AuthenticateCapability wraps the scopedKeeper's AuthenticateCapability function
func (k Keeper) AuthenticateCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) bool {
	return k.scopedKeeper.AuthenticateCapability(ctx, cap, name)
}

// GetPort returns the portID for the transfer module. Used in ExportGenesis
func (k Keeper) GetPort(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)
	return string(store.Get(types.PortKey))
}

// SetPort sets the portID for the transfer module. Used in InitGenesis
func (k Keeper) SetPort(ctx sdk.Context, portID string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.PortKey, []byte(portID))
}

// IsBound checks if the transfer module is already bound to the desired port
func (k Keeper) IsBound(ctx sdk.Context, portID string) bool {
	_, ok := k.scopedKeeper.GetCapability(ctx, host.PortPath(portID))
	return ok
}

// BindPort defines a wrapper function for the ort Keeper's function in
// order to expose it to module's InitGenesis function
func (k Keeper) BindPort(ctx sdk.Context, portID string) error {
	cap := k.portKeeper.BindPort(ctx, portID)
	return k.ClaimCapability(ctx, cap, host.PortPath(portID))
}