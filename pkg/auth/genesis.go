package auth

import "github.com/ci123chain/ci123chain/pkg/abci/types"

func InitGenesis(ctx types.Context, ak AuthKeeper, data GenesisState) {
	ak.SetParams(ctx, data.Params)
}
