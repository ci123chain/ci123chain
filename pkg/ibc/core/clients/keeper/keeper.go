package keeper

import (
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	prefix "github.com/ci123chain/ci123chain/pkg/abci/store"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/clients/types"
	commitmenttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/commitment/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/host"
	ibcclienttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/light-clients/07-tendermint/types"
	paramtypes "github.com/ci123chain/ci123chain/pkg/params/subspace"
	"github.com/tendermint/tendermint/libs/log"
	"reflect"
	"strings"
)

// Keeper represents a types that grants read and write permissions to any client
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


func (k Keeper) GetClientConsensusState(ctx sdk.Context, clientID string, height exported.Height) (exported.ConsensusState, bool) {
	store := k.ClientStore(ctx, clientID)
	bz := store.Get(host.ConsensusStateKey(height))
	if bz == nil {
		return nil, false
	}

	consensusState := k.MustUnmarshalConsensusState(bz)
	return consensusState, true
}


func (k Keeper) GetSelfConsensusState(ctx sdk.Context, height exported.Height) (exported.ConsensusState, bool) {
	selfHeight, ok := height.(types.Height)
	if !ok {
		return nil, false
	}
	// check that height revision matches chainID revision
	revision := types.ParseChainID(ctx.ChainID())
	if revision != height.GetRevisionNumber() {
		return nil, false
	}
	histInfo, found := k.stakingKeeper.GetHistoricalInfo(ctx, int64(selfHeight.RevisionHeight))
	if !found {
		return nil, false
	}

	consensusState := &ibcclienttypes.ConsensusState{
		Timestamp:          histInfo.Header.Time,
		Root:               commitmenttypes.NewMerkleRoot(histInfo.Header.GetAppHash()),
		NextValidatorsHash: histInfo.Header.NextValidatorsHash,
	}
	return consensusState, true
}

func (k Keeper) ValidateSelfClient(ctx sdk.Context, clientState exported.ClientState) error {
	tmClient, ok := clientState.(*ibcclienttypes.ClientState)
	if !ok {
		return sdkerrors.Wrapf(types.ErrInvalidClient, "client must be a Tendermint client, expected: %T, got: %T",
			&ibcclienttypes.ClientState{}, tmClient)
	}

	if clientState.IsFrozen() {
		return types.ErrClientFrozen
	}

	if ctx.ChainID() != tmClient.ChainId {
		return sdkerrors.Wrapf(types.ErrInvalidClient, "invalid chain-id. expected: %s, got: %s",
			ctx.ChainID(), tmClient.ChainId)
	}

	revision := types.ParseChainID(ctx.ChainID())

	// client must be in the same revision as executing chain
	if tmClient.LatestHeight.RevisionNumber != revision {
		return sdkerrors.Wrapf(types.ErrInvalidClient, "client is not in the same revision as the chain. expected revision: %d, got: %d",
			tmClient.LatestHeight.RevisionNumber, revision)
	}

	selfHeight := types.NewHeight(revision, uint64(ctx.BlockHeight()))
	if tmClient.LatestHeight.GTE(selfHeight) {
		return sdkerrors.Wrapf(types.ErrInvalidClient, "client has LatestHeight %d greater than or equal to chain height %d",
			tmClient.LatestHeight, selfHeight)
	}

	expectedProofSpecs := commitmenttypes.GetSDKSpecs()
	if !reflect.DeepEqual(expectedProofSpecs, tmClient.ProofSpecs) {
		return sdkerrors.Wrapf(types.ErrInvalidClient, "client has invalid proof specs. expected: %v got: %v",
			expectedProofSpecs, tmClient.ProofSpecs)
	}

	//if err := light.ValidateTrustLevel(tmClient.TrustLevel.ToTendermint()); err != nil {
	//	return sdkerrors.Wrapf(types.ErrInvalidClient, "trust-level invalid: %v", err)
	//}

	expectedUbdPeriod := k.stakingKeeper.UnbondingTime(ctx)
	if expectedUbdPeriod != tmClient.UnbondingPeriod {
		return sdkerrors.Wrapf(types.ErrInvalidClient, "invalid unbonding period. expected: %s, got: %s",
			expectedUbdPeriod, tmClient.UnbondingPeriod)
	}

	if tmClient.UnbondingPeriod < tmClient.TrustingPeriod {
		return sdkerrors.Wrapf(types.ErrInvalidClient, "unbonding period must be greater than trusting period. unbonding period (%d) < trusting period (%d)",
			tmClient.UnbondingPeriod, tmClient.TrustingPeriod)
	}

	//if len(tmClient.UpgradePath) != 0 {
	//	// For now, SDK IBC implementation assumes that upgrade path (if defined) is defined by SDK upgrade module
	//	expectedUpgradePath := []string{upgradetypes.StoreKey, upgradetypes.KeyUpgradedIBCState}
	//	if !reflect.DeepEqual(expectedUpgradePath, tmClient.UpgradePath) {
	//		return sdkerrors.Wrapf(types.ErrInvalidClient, "upgrade path must be the upgrade path defined by upgrade module. expected %v, got %v",
	//			expectedUpgradePath, tmClient.UpgradePath)
	//	}
	//}
	return nil
}

func (k Keeper) IterateClients(ctx sdk.Context, cb func(string, exported.ClientState) bool) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, host.KeyClientStorePrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		keySplit := strings.Split(string(iterator.Key()), "/")
		if keySplit[len(keySplit)-1] != host.KeyClientState {
			continue
		}
		clientState := k.MustUnmarshalClientState(iterator.Value())

		// key is ibc/{clientid}/clientState
		// Thus, keySplit[1] is clientID
		if cb(keySplit[1], clientState) {
			break
		}
	}
}

func (k Keeper) SetClientConsensusState(ctx sdk.Context, clientID string, height exported.Height,
	consensusState exported.ConsensusState) {
	store := k.ClientStore(ctx, clientID)
	store.Set(host.ConsensusStateKey(height), k.MustMarshalConsensusState(consensusState))
}

// GetLatestClientConsensusState gets the latest ConsensusState stored for a given client
func (k Keeper) GetLatestClientConsensusState(ctx sdk.Context, clientID string) (exported.ConsensusState, bool) {
	clientState, ok := k.GetClientState(ctx, clientID)
	if !ok {
		return nil, false
	}
	return k.GetClientConsensusState(ctx, clientID, clientState.GetLatestHeight())
}