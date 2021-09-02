package connection

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/connection/keeper"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/connection/types"
)

// InitGenesis initializes the ibc connection submodule's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, gs types.GenesisState) {
	for _, connection := range gs.Connections {
		conn := types.NewConnectionEnd(connection.State, connection.ClientId, connection.Counterparty, connection.Versions, connection.DelayPeriod)
		k.SetConnection(ctx, connection.Id, conn)
	}
	for _, connPaths := range gs.ClientConnectionPaths {
		k.SetClientConnectionPaths(ctx, connPaths.ClientId, connPaths.Paths)
	}
	k.SetNextConnectionSequence(ctx, gs.NextConnectionSequence)
}


func ExportGenesis(ctx sdk.Context, k keeper.Keeper) types.GenesisState {

	css := make([]types.ConnectionEnd, 0)
	k.GetConnects(ctx, func(end types.ConnectionEnd) (stop bool) {
		css = append(css, end)
		return false
	})
	ics := make([]types.IdentifiedConnection, 0)
	cps := make([]types.ConnectionPaths, 0)
	for _, v := range css {
		ic := types.NewIdentifiedConnection(v.ClientId, v)
		ics = append(ics, ic)
		cp, _ := k.GetClientConnectionPaths(ctx, v.ClientId)
		cps = append(cps, types.ConnectionPaths{
			ClientId: v.ClientId,
			Paths:    cp,
		})
	}
	gs := types.GenesisState{
		Connections:            ics,
		ClientConnectionPaths:  cps,
		NextConnectionSequence: k.GetNextConnectionSequence(ctx),
	}
	return gs
}