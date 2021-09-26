package pre_staking

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/pre_staking/keeper"
)

func InitGenesis(ctx sdk.Context, k keeper.PreStakingKeeper, data GenesisState) {

	var contract sdk.AccAddress
	//deploy dao contract.
	if !data.DaoDeployed {
		var err error
		contract, err = k.SupplyKeeper.DeployDaoContract(ctx, ModuleName, nil)
		if err != nil {
			panic(err)
		}
	}

	//store account prestaking record.
	for _, v := range data.Records {
		var all = v.PrestakingAmount
		for _, val := range v.Delegations {
			err := k.SetAccountStakingRecord(ctx, val.Validator, v.Delegator, val.StorageTime, val.Amount)
			if err != nil {
				panic(err)
			}
			all = all.Sub(val.Amount.Amount.ToDec().RoundInt())
		}
		k.SetAccountPreStaking(ctx, v.Delegator, all)
		//mint token.
		if !data.DaoDeployed {
			err := k.SupplyKeeper.Mint(ctx, contract, v.Delegator, ModuleName, all.BigInt())
			if err != nil {
				panic(err)
			}
		}
	}
}