package keeper

import (
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	prefix "github.com/ci123chain/ci123chain/pkg/abci/store"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/clients/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/host"
	paramtypes "github.com/ci123chain/ci123chain/pkg/params/subspace"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper represents a type that grants read and write permissions to any client
// state information
type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           *codec.Codec
	paramSpace    paramtypes.Subspace
	stakingKeeper types.StakingKeeper
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramSpace paramtypes.Subspace, sk types.StakingKeeper) Keeper {
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}
	return Keeper{
		storeKey:      key,
		cdc:           cdc,
		paramSpace:    paramSpace,
		stakingKeeper: sk,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+host.ModuleName+"/"+types.SubModuleName)
}

// GetClientState gets a particular client from the store
func (k Keeper) GetClientState(ctx sdk.Context, clientID string) (exported.ClientState, bool) {
	store := k.ClientStore(ctx, clientID)
	bz := store.Get(host.ClientStateKey())
	if bz == nil {
		return nil, false
	}

	clientState := k.MustUnmarshalClientState(bz)
	return clientState, true
}

func (k Keeper) ClientStore(ctx sdk.Context, clientID string) sdk.KVStore {
	clientPrefix := []byte(fmt.Sprintf("%s/%s/", host.KeyClientStorePrefix, clientID))
	return prefix.NewPrefixStore(ctx.KVStore(k.storeKey), clientPrefix)
}



func (k Keeper) SetClientConsensusState(ctx sdk.Context, clientID string, height exported.Height,
	consensusState exported.ConsensusState) {
	store := k.ClientStore(ctx, clientID)
	store.Set(host.ConsensusStateKey(height), k.MustMarshalConsensusState(consensusState))
}