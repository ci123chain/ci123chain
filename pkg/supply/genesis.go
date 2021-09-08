package supply

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/supply/types"
)

func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {

	k.SetSupply(ctx, types.NewSupply(data.Supply))
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	return types.NewGenesisState(k.GetSupply(ctx).GetTotal())
}