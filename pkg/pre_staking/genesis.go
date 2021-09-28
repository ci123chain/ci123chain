package pre_staking

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/pre_staking/keeper"
	"github.com/ci123chain/ci123chain/pkg/pre_staking/types"
	"github.com/umbracle/go-web3"
	"math/big"
)

const (
	moduleAcc = "0x3F43E75Aaba2c2fD6E227C10C6E7DC125A93DE3c"
)

func InitGenesis(ctx sdk.Context, k keeper.PreStakingKeeper, data GenesisState) {

	if !data.DaoDeployed {
		a, _ := new(big.Int).SetString("10000000000000000000", 10)

		v1, _ := new(big.Int).SetString("500000000000000000", 10)
		v2, _ := new(big.Int).SetString("150000000000000000", 10)
		v3, _ := new(big.Int).SetString("86400", 10)

		zero, _ := new(big.Int).SetString("0", 10)

		contractAddr, err := k.SupplyKeeper.DeployDaoContract(ctx, types.ModuleName, []interface{}{[]web3.Address{web3.HexToAddress(moduleAcc)}, []*big.Int{a}, [3]*big.Int{v1, v2, v3}, zero})
		if err != nil {
			panic(err)
		}
		ctx.Logger().Info("deployed weeLink dao contract", "contract address", contractAddr)
	}

	for _, v := range data.Records.PrestakingRecord {
		k.SetAccountPreStaking(ctx, v.Delegator, v.Staking)
	}

	for _, v := range data.Records.DelStakingRecords {
		k.SetAccountStakingRecords(ctx, v.Delegator, v.Validator, v.Records)
	}
}


func ExportGenesis(ctx sdk.Context, k keeper.PreStakingKeeper) types.GenesisState {
	var records types.DelegationRecord
	pr := k.GetAllAccountPreStaking(ctx)
	sr := k.GetAllStakingRecords(ctx)
	records.PrestakingRecord = pr
	records.DelStakingRecords = sr
	return types.NewGenesisState(records, true)
}