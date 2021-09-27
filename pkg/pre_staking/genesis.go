package pre_staking

import (
	//"errors"
	//"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/pre_staking/keeper"
	//"github.com/ci123chain/ci123chain/pkg/pre_staking/types"
	//"math/big"
)

func InitGenesis(ctx sdk.Context, k keeper.PreStakingKeeper, data GenesisState) {

	//var contract sdk.AccAddress
	//deploy dao contract.
	//if !data.DaoDeployed {
	//	var err error
	//	contract, err = k.SupplyKeeper.DeployDaoContract(ctx, ModuleName, nil)
	//	if err != nil {
	//		panic(err)
	//	}
	//}

	//store account prestaking record.
	for _, _ = range data.Records {
		//var all = make([]types.Vault, 0)
		//var latestId = new(big.Int).SetUint64(0)
		//for key, val := range v.Delegations {
		//	id, ok := new(big.Int).SetString(key, 64)
		//	if !ok {
		//		panic(errors.New(fmt.Sprintf("invalid vault_id: %v", key)))
		//	}
		//	err := k.SetAccountStakingRecord(ctx, val.Validator, v.Delegator, id, val.EndTime, val.Amount)
		//	if err != nil {
		//		panic(err)
		//	}
		//	sr := types.NewStakingRecord()
		//	v := types.NewVault(val.StartTime, val.EndTime, val.StorageTime, val.Amount)
		//	if latestId.Cmp(id) == -1 {
		//		latestId = id
		//	}
		//}
		//k.SetAccountPreStaking(ctx, v.Delegator, all)
		//mint token.
		//if !data.DaoDeployed {
		//	err := k.SupplyKeeper.Mint(ctx, contract, v.Delegator, ModuleName, all.BigInt())
		//	if err != nil {
		//		panic(err)
		//	}
		//}
	}
}