package params

import (
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/params/subspace"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper of the global paramstore
type Keeper struct {
	cdc       *codec.Codec
	key       types.StoreKey
	tkey      types.StoreKey
	codespace types.CodespaceType
	spaces    map[string]*Subspace

}

// NewKeeper constructs a params keeper
func NewKeeper(cdc *codec.Codec, key *types.KVStoreKey, tkey *types.TransientStoreKey, codespace types.CodespaceType) (k Keeper) {
	k = Keeper{
		cdc:       cdc,
		key:       key,
		tkey:      tkey,
		codespace: codespace,
		spaces:    make(map[string]*Subspace),
	}
	return k
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx types.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", ModuleName))
}

// Allocate subspace used for keepers
func (k Keeper) Subspace(s string) Subspace {
	_, ok := k.spaces[s]
	if ok {
		panic("subspace already occupied")
	}

	if s == "" {
		panic("cannot use empty string for subspace")
	}

	space := subspace.NewSubspace(k.cdc, k.key, k.tkey, s)
	k.spaces[s] = &space

	return space
}

// Get existing substore from keeper
func (k Keeper) GetSubspace(s string) (Subspace, bool) {
	space, ok := k.spaces[s]
	if !ok {
		return Subspace{}, false
	}
	return *space, ok
}
