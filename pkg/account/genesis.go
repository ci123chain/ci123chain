package account

import (
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/account/exported"
	"github.com/ci123chain/ci123chain/pkg/account/keeper"
)

func InitGenesis(ctx types.Context, _ *codec.Codec, accountKeeper keeper.AccountKeeper, genesisState GenesisState) {
	for _, gacc := range genesisState {
		acc := gacc.ToAccount()
		acc = accountKeeper.NewAccount(ctx, acc)
		accountKeeper.SetAccount(ctx, acc)
	}
}

func ExportGenesis(ctx types.Context, ak keeper.AccountKeeper) GenesisState {
	var genAccounts GenesisAccounts
	ak.IterateAccounts(ctx, func(account exported.Account) bool {
		genAccount := NewGenesisAccountRaw(account.GetAddress(), account.GetCoins())
		genAccount.Sequence = account.GetSequence()
		genAccount.AccountNumber = account.GetAccountNumber()
		genAccounts = append(genAccounts, genAccount)
		return false
	})
	return NewGensisState(genAccounts)
}