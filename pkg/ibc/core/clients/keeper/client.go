package keeper

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/clients/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/host"
)

func (k Keeper)CreateClient(ctx sdk.Context, clientState exported.ClientState,
	consensusState exported.ConsensusState ) (string, sdk.Error) {
	params := k.GetParams(ctx)
	if params.IsAllowedClient(clientState.ClientType()) {
		return "", types.ErrInvalidClientType(types.DefaultCodespace, fmt.Errorf("client state type %s is not registered in the allowlist", clientState.ClientType()))
	}

	clientID := k.GenerateClientIdentifier(ctx, clientState.ClientType())
	k.SetClientState(ctx, clientID, clientState)
	k.Logger(ctx).Info("client created at height", "client-id", clientID, "height", clientState.GetLastHeight().String())

	if err := clientState.Initialize(ctx, k.ClientStore(ctx, clientID), consensusState); err != nil {
		return "", types.ErrInitClientState(types.DefaultCodespace, err)
	}

	if consensusState != nil {
		k.SetClientConsensusState(ctx, clientID, clientState.GetLatestHeight(), consensusState)
	}

	k.Logger(ctx).Info("client created at height", "client-id", clientID, "height", clientState.GetLastHeight().String())

	// todo  for telemetry
	return clientID, nil
}

// GenerateClientIdentifier returns the next client identifier.
func (k Keeper) GenerateClientIdentifier(ctx sdk.Context, clientType string) string {
	nextClientSeq := k.GetNextClientSequence(ctx)
	clientID := types.FormatClientIdentifier(clientType, nextClientSeq)

	nextClientSeq++
	k.SetNextClientSequence(ctx, nextClientSeq)
	return clientID
}

// SetClientState sets a particular Client to the store
func (k Keeper) SetClientState(ctx sdk.Context, clientID string, clientState exported.ClientState) {
	store := k.ClientStore(ctx, clientID)
	store.Set(host.ClientStateKey(), k.MustMarshalClientState(clientState))
}

// GetNextClientSequence gets the next client sequence from the store.
func (k Keeper) GetNextClientSequence(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte(types.KeyNextClientSequence))
	if bz == nil {
		panic("next client sequence is nil")
	}

	return sdk.BigEndianToUint64(bz)
}

// SetNextClientSequence sets the next client sequence to the store.
func (k Keeper) SetNextClientSequence(ctx sdk.Context, sequence uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := sdk.Uint64ToBigEndian(sequence)
	store.Set([]byte(types.KeyNextClientSequence), bz)
}
