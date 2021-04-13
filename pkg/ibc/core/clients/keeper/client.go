package keeper

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/clients/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/host"
)

func (k Keeper)CreateClient(ctx sdk.Context, clientState exported.ClientState,
	consensusState exported.ConsensusState ) (string, error) {
	params := k.GetParams(ctx)
	if !params.IsAllowedClient(clientState.ClientType()) {
		return "", sdkerrors.Wrapf(
			types.ErrInvalidClientType,
			"client state type %s is not registered in the allowlist", clientState.ClientType(),
		)	}

	clientID := k.GenerateClientIdentifier(ctx, clientState.ClientType())
	k.SetClientState(ctx, clientID, clientState)
	k.Logger(ctx).Info("client created at height", "client-id", clientID, "height", clientState.GetLatestHeight().String())

	if err := clientState.Initialize(ctx, k.ClientStore(ctx, clientID), consensusState); err != nil {
		return "", err
	}

	if consensusState != nil {
		k.SetClientConsensusState(ctx, clientID, clientState.GetLatestHeight(), consensusState)
	}

	k.Logger(ctx).Info("client created at height", "client-id", clientID, "height", clientState.GetLatestHeight().String())

	// todo  for telemetry
	return clientID, nil
}


// UpdateClient updates the consensus state and the state root from a provided header.
func (k Keeper) UpdateClient(ctx sdk.Context, clientID string, header exported.Header) error {
	clientState, found := k.GetClientState(ctx, clientID)
	if !found {
		return sdkerrors.Wrapf(types.ErrClientNotFound, "cannot update client with ID %s", clientID)
	}

	// prevent update if the client is frozen before or at header height
	if clientState.IsFrozen() && clientState.GetFrozenHeight().LTE(header.GetHeight()) {
		return sdkerrors.Wrapf(types.ErrClientFrozen, "cannot update client with ID %s", clientID)
	}

	clientState, consensusState, err := clientState.CheckHeaderAndUpdateState(ctx, k.cdc, k.ClientStore(ctx, clientID), header)
	if err != nil {
		return sdkerrors.Wrapf(err, "cannot update client with ID %s", clientID)
	}

	k.SetClientState(ctx, clientID, clientState)

	var consensusHeight exported.Height

	// we don't set consensus state for localhost client
	if header != nil && clientID != exported.Localhost {
		k.SetClientConsensusState(ctx, clientID, header.GetHeight(), consensusState)
		consensusHeight = header.GetHeight()
	} else {
		consensusHeight = types.GetSelfHeight(ctx)
	}

	k.Logger(ctx).Info("client state updated", "client-id", clientID, "height", consensusHeight.String())

	//defer func() {
	//	telemetry.IncrCounterWithLabels(
	//		[]string{"ibc", "client", "update"},
	//		1,
	//		[]metrics.Label{
	//			telemetry.NewLabel("client-type", clientState.ClientType()),
	//			telemetry.NewLabel("client-id", clientID),
	//			telemetry.NewLabel("update-type", "msg"),
	//		},
	//	)
	//}()

	// emitting events in the keeper emits for both begin block and handler client updates
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUpdateClient,
			sdk.NewAttributeString(types.AttributeKeyClientID, clientID),
			sdk.NewAttributeString(types.AttributeKeyClientType, clientState.ClientType()),
			sdk.NewAttributeString(types.AttributeKeyConsensusHeight, consensusHeight.String()),
		),
	)

	return nil
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
