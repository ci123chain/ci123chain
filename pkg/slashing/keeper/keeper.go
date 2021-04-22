package keeper

import (
	"fmt"
	staking "github.com/ci123chain/ci123chain/pkg/staking/keeper"
	"github.com/tendermint/tendermint/crypto"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/ci123chain/ci123chain/pkg/abci/codec"

	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/slashing/types"
)

// Keeper of the slashing store
type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        *codec.Codec
	sk         staking.StakingKeeper
	paramspace types.ParamSubspace
}

// NewKeeper creates a slashing keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, sk staking.StakingKeeper, paramspace types.ParamSubspace) Keeper {
	// set KeyTable if it has not already been set
	if !paramspace.HasKeyTable() {
		paramspace = paramspace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		storeKey:   key,
		cdc:        cdc,
		sk:         sk,
		paramspace: paramspace,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// AddPubkey sets a address-pubkey relation
func (k Keeper) AddPubkey(ctx sdk.Context, pubkey crypto.PubKey) error {
	bz, err := k.cdc.MarshalBinaryLengthPrefixed(pubkey)
	if err != nil {
		return err
	}
	store := ctx.KVStore(k.storeKey)
	key := types.AddrPubkeyRelationKey(pubkey.Address())
	store.Set(key, bz)
	return nil
}

// GetPubkey returns the pubkey from the adddress-pubkey relation
func (k Keeper) GetPubkey(ctx sdk.Context, a sdk.AccAddress) (crypto.PubKey, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.AddrPubkeyRelationKey(a.Bytes()))
	if bz == nil {
		return nil, fmt.Errorf("address %s not found", sdk.AccAddress(a))
	}
	var pk crypto.PubKey
	return pk, k.cdc.UnmarshalBinaryLengthPrefixed(bz, &pk)
}

// Slash attempts to slash a validator. The slash is delegated to the staking
// module to make the necessary validator changes.
func (k Keeper) Slash(ctx sdk.Context, consAddr sdk.AccAddress, fraction sdk.Dec, power, distributionHeight int64) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSlash,
			sdk.NewAttribute([]byte(types.AttributeKeyAddress), consAddr.Bytes()),
			sdk.NewAttribute([]byte(types.AttributeKeyPower), []byte(fmt.Sprintf("%d", power))),
			sdk.NewAttribute([]byte(types.AttributeKeyReason), []byte(types.AttributeValueDoubleSign)),
		),
	)

	k.sk.Slash(ctx, consAddr, distributionHeight, power, fraction)
}

// Jail attempts to jail a validator. The slash is delegated to the staking module
// to make the necessary validator changes.
func (k Keeper) Jail(ctx sdk.Context, consAddr sdk.AccAddress) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSlash,
			sdk.NewAttribute([]byte(types.AttributeKeyJailed), consAddr.Bytes()),
		),
	)

	k.sk.Jail(ctx, consAddr)
}

func (k Keeper) deleteAddrPubkeyRelation(ctx sdk.Context, addr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.AddrPubkeyRelationKey(addr.Bytes()))
}
