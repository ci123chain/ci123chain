package mint

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/mint/keeper"
)

func InitGenesis(ctx sdk.Context, keeper keeper.MinterKeeper, data GenesisState) {
	keeper.SetMinter(ctx, data.Minter)
	keeper.SetParams(ctx, data.Params)
}


func ExportGenesis(ctx sdk.Context, keeper keeper.MinterKeeper) GenesisState {

	minter := keeper.GetMinter(ctx)
	params := keeper.GetParams(ctx)
	return GenesisState(NewGenesisState(minter, params))
}