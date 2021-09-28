package pre_staking

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/pre_staking/keeper"
	"github.com/ci123chain/ci123chain/pkg/pre_staking/types"
	"github.com/umbracle/go-web3"
	"math/big"
)

func BeginBlock() {}

func EndBlock(ctx sdk.Context, k keeper.PreStakingKeeper) {

	var a = k.GetWeeLinkDao(ctx)
	if  a == "" {
		a, _ := new(big.Int).SetString("10000000000000000000", 10)

		v1, _ := new(big.Int).SetString("500000000000000000", 10)
		v2, _ := new(big.Int).SetString("150000000000000000", 10)
		v3, _ := new(big.Int).SetString("86400", 10)

		zero, _ := new(big.Int).SetString("0", 10)

		contractAddr, err := k.SupplyKeeper.DeployDaoContract(ctx, types.ModuleName, []interface{}{[]web3.Address{web3.HexToAddress(moduleAcc)}, []*big.Int{a}, [3]*big.Int{v1, v2, v3}, zero})

		if err != nil {
			panic(err)
		}
		k.SetWeeLinkDao(ctx, contractAddr)
		ctx.Logger().Info("deployed weeLink dao contract", "contract address", contractAddr)
	}
	k.UpdateDeadlineRecord(ctx)
}
