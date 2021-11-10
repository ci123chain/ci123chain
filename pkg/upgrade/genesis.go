package upgrade

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/upgrade/keeper"
	"github.com/ci123chain/ci123chain/pkg/upgrade/types"
)

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) types.GenesisState {
	return types.NewGenesisState()
}