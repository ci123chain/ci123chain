package keeper

import (
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	"github.com/ci123chain/ci123chain/pkg/abci/store"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	commitmenttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/commitment/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/connection/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/host"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/log"
)

type Keeper struct {
	storeKey	sdk.StoreKey
	cdc 		*codec.Codec
	clientKeeper types.ClientKeeper
}


func NewKeeper(cdc *codec.Codec, key store.StoreKey, ck types.ClientKeeper) Keeper {
	return Keeper{
		storeKey: key,
		cdc: cdc,
		clientKeeper: ck,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", host.ModuleName + types.SubModuleName)
}

func (k Keeper) GetCommitmentPrefix() exported.Prefix {
	return commitmenttypes.NewMerklePrefix([]byte(k.storeKey.Name()))
}

func (k Keeper) GenerateConnectionIdentifier(ctx sdk.Context) string {
	nextConnSeq := k.GetNextConnectionSequence(ctx)
	connectionID := types.FormatConnectionIdentifier(nextConnSeq)
	nextConnSeq++
	k.SetNextConnectionSequence(ctx, nextConnSeq)
	return connectionID
}

func (k Keeper) SetConnection(ctx sdk.Context, connectionID string, connection types.ConnectionEnd)  {
	connectionStore := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(&connection)
	connectionStore.Set(host.ConnectionKey(connectionID), bz)
}

func (k Keeper) GetConnection(ctx sdk.Context, connectionID string) (types.ConnectionEnd, bool) {
	connectionStore := ctx.KVStore(k.storeKey)
	bz := connectionStore.Get(host.ConnectionKey(connectionID))
	if bz == nil {
		return types.ConnectionEnd{}, false
	}

	var connection types.ConnectionEnd
	k.cdc.MustUnmarshalBinaryBare(bz, &connection)
	return connection, true
}

func (k Keeper) GetTimestampAtHeight(ctx sdk.Context, connection types.ConnectionEnd, height exported.Height) (uint64, error) {
	consensusState, found := k.clientKeeper.GetClientConsensusState(
		ctx, connection.GetClientID(), height,
	)

	if !found {
		return 0, errors.Errorf("clientID (%s), height (%s)", connection.GetClientID(), height)
	}

	return consensusState.GetTimestamp(), nil
}


// GetNextConnectionSequence gets the next connection sequence from the store.
func (k Keeper) GetNextConnectionSequence(ctx sdk.Context) uint64 {
	connectionStore := ctx.KVStore(k.storeKey)
	bz := connectionStore.Get([]byte(types.KeyNextConnectionSequence))
	if bz == nil {
		panic("next connection sequence is nil")
	}

	return sdk.BigEndianToUint64(bz)
}
// SetNextConnectionSequence sets the next connection sequence to the store.
func (k Keeper) SetNextConnectionSequence(ctx sdk.Context, sequence uint64) {
	connectionStore := ctx.KVStore(k.storeKey)
	bz := sdk.Uint64ToBigEndian(sequence)
	connectionStore.Set([]byte(types.KeyNextConnectionSequence), bz)
}

// GetClientConnectionPaths returns all the connection paths stored under a
// particular client
func (k Keeper) GetClientConnectionPaths(ctx sdk.Context, clientID string) ([]string, bool) {
	connectionStore := ctx.KVStore(k.storeKey)
	bz := connectionStore.Get(host.ClientConnectionsKey(clientID))
	if bz == nil {
		return nil, false
	}

	var clientPaths types.ClientPaths
	k.cdc.MustUnmarshalBinaryBare(bz, &clientPaths)
	return clientPaths.Paths, true
}

// SetClientConnectionPaths sets the connections paths for client
func (k Keeper) SetClientConnectionPaths(ctx sdk.Context, clientID string, paths []string) {
	connectionStore := ctx.KVStore(k.storeKey)
	clientPaths := types.ClientPaths{Paths: paths}
	bz := k.cdc.MustMarshalBinaryBare(&clientPaths)
	connectionStore.Set(host.ClientConnectionsKey(clientID), bz)
}

// addConnectionToClient is used to add a connection identifier to the set of
// connections associated with a client.
func (k Keeper) addConnectionToClient(ctx sdk.Context, clientID, connectionID string) error {
	_, found := k.clientKeeper.GetClientState(ctx, clientID)
	if !found {
		return errors.Errorf("client not found for client_id: %s", clientID)
	}

	conns, found := k.GetClientConnectionPaths(ctx, clientID)
	if !found {
		conns = []string{}
	}

	conns = append(conns, connectionID)
	k.SetClientConnectionPaths(ctx, clientID, conns)
	return nil
}
