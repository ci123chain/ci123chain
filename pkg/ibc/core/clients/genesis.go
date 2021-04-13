package clients

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/clients/keeper"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/clients/types"
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
		cs := client.ClientState

		if !gs.Params.IsAllowedClient(client.ClientState.ClientType()) {
			panic(fmt.Sprintf("client state type %s is not registered on the allowlist", cs.ClientType()))
		}

		k.SetClientState(ctx, client.ClientId, cs)
	}

	for _, cs := range gs.ClientsConsensus {
		for _, consState := range cs.ConsensusStates {
			consensusState := consState.ConsensusState
			k.SetClientConsensusState(ctx, cs.ClientId, consState.Height, consensusState)
		}
	}

	k.SetNextClientSequence(ctx, gs.NextClientSequence)

	// NOTE: localhost creation is specifically disallowed for the time being.
	// Issue: https://github.com/cosmos/cosmos-sdk/issues/7871
}
