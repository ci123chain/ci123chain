package core

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/channel"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/clients"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/connection"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/keeper"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/types"
)

// InitGenesis initializes the ibc state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, createLocalhost bool, gs *types.GenesisState) {
	clients.InitGenesis(ctx, k.ClientKeeper, gs.ClientGenesis)
	connection.InitGenesis(ctx, k.ConnectionKeeper, gs.ConnectionGenesis)
	channel.InitGenesis(ctx, k.ChannelKeeper, gs.ChannelGenesis)
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) types.GenesisState {

	gs := types.GenesisState{
		ClientGenesis:     clients.ExportGenesis(ctx, k.ClientKeeper, ""),
		ConnectionGenesis: connection.ExportGenesis(ctx, k.ConnectionKeeper),
		ChannelGenesis:    channel.ExportGenesis(ctx, k.ChannelKeeper),
	}
	return gs
}