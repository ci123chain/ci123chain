package auth

import "github.com/tanhuiya/ci123chain/pkg/abci/types"

func InitGenesis(ctx types.Context, ak AuthKeeper, data GenesisState) {
	ak.SetParams(ctx, data.Params)
}
