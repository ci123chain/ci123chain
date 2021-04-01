package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
)

type ClientKeeper interface {
	GetClientState(ctx sdk.Context, clientID string) (exported.ClientState, bool)
	GetClientConsensusState(ctx sdk.Context, clientID string, height exported.Height) (exported.ConsensusState, bool)
	GetSelfConsensusState(ctx sdk.Context, height exported.Height) (exported.ConsensusState, bool)
	ValidateSelfClient(ctx sdk.Context, clientState exported.ClientState) error
	IterateClients(ctx sdk.Context, cb func(string, exported.ClientState) bool)
	ClientStore(ctx sdk.Context, clientID string) sdk.KVStore
}
