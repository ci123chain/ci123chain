package account

import (
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/account/exported"
	"github.com/ci123chain/ci123chain/pkg/account/keeper"
)

func InitGenesis(ctx types.Context, _ *codec.Codec, accountKeeper keeper.AccountKeeper, genesisState GenesisState) {
	for _, gacc := range genesisState {
		//acc := gacc.ToAccount()
		//acc = accountKeeper.NewAccount(ctx, acc)
		//if gacc.IsModule {
		//	accountKeeper.
		//}
		//if gacc.Name != "" {
		//	//acc := gacc.ToModuleAccount()
		//	acc := types3.NewModuleAccountFromBaseAccount(gacc.BaseAccount, gacc.Name, gacc.Permissions...)
		//	accountKeeper.SetAccount(ctx, acc)
		//}else {
		//	acc := gacc.ToAccount()
		//	accountKeeper.SetAccount(ctx, acc)
		//}
		accountKeeper.SetAccount(ctx, gacc)
	}
}

func ExportGenesis(ctx types.Context, ak keeper.AccountKeeper) GenesisState {
	var genAccounts GenesisAccounts
	ak.IterateAccounts(ctx, func(account exported.Account) bool {
		//var genAccount GenesisAccount
		//macc, ok := account.(exported2.ModuleAccountI)
		//if ok {
		//	genAccount = NewGenesisAccountRaw(types2.NewBaseAccountFromExportAccount(account), macc.GetName(), macc.GetPermissions()...)
		//}else {
		//	genAccount = NewGenesisAccountRaw(types2.NewBaseAccountFromExportAccount(account), "", "")
		//}
		//genAccount := NewGenesisAccountRaw(account)
		//genAccount.Sequence = account.GetSequence()
		//genAccount.AccountNumber = account.GetAccountNumber()
		genAccounts = append(genAccounts, account)
		return false
	})
	return NewGensisState(genAccounts)
}