package clients

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/clients/keeper"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/clients/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
)

// InitGenesis initializes the ibc client submodule's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, gs types.GenesisState) {
	k.SetParams(ctx, gs.Params)

	// Set all client metadata first. This will allow client keeper to overwrite client and consensus state keys
	// if clients accidentally write to ClientKeeper reserved keys.
	//if len(gs.ClientsMetadata) != 0 {
	//	k.SetAllClientMetadata(ctx, gs.ClientsMetadata)
	//}

	for _, client := range gs.Clients {
		cs, ok := client.ClientState.GetCachedValue().(exported.ClientState)
		if !ok {
			panic("invalid client state")
		}

		if !gs.Params.IsAllowedClient(cs.ClientType()) {
			panic(fmt.Sprintf("client state type %s is not registered on the allowlist", cs.ClientType()))
		}

		k.SetClientState(ctx, client.ClientId, cs)
	}

	for _, cs := range gs.ClientsConsensus {
		for _, consState := range cs.ConsensusStates {
			consensusState, ok := consState.ConsensusState.GetCachedValue().(exported.ConsensusState)
			if !ok {
				panic(fmt.Sprintf("invalid consensus state with client ID %s at height %s", cs.ClientId, consState.Height))
			}
			k.SetClientConsensusState(ctx, cs.ClientId, consState.Height, consensusState)
		}
	}

	k.SetNextClientSequence(ctx, gs.NextClientSequence)

	// NOTE: localhost creation is specifically disallowed for the time being.
	// Issue: https://github.com/cosmos/cosmos-sdk/issues/7871
}


func ExportGenesis(ctx sdk.Context, k keeper.Keeper, clientID string) types.GenesisState {
	ps := k.GetParams(ctx)

	gs := types.GenesisState{
		Clients:            types.IdentifiedClientStates{},
		ClientsConsensus:   types.ClientsConsensusStates{},
		ClientsMetadata:    []types.IdentifiedGenesisMetadata{},
		Params:             ps,
		NextClientSequence: k.GetNextClientSequence(ctx),
	}
	return gs
}