package account

import (
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/account/keeper"
)

func InitGenesis(ctx types.Context, _ *codec.Codec, accountKeeper keeper.AccountKeeper, genesisState GenesisState) {
	for _, gacc := range genesisState {
		acc := gacc.ToAccount()
		acc = accountKeeper.NewAccount(ctx, acc)
		accountKeeper.SetAccount(ctx, acc)
	}
}