package pre_staking

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/pre_staking/keeper"
	"github.com/ci123chain/ci123chain/pkg/pre_staking/types"
)

func InitGenesis(ctx sdk.Context, k keeper.PreStakingKeeper, data GenesisState) {

	if data.StakingToken != "" {
		k.SetTokenManager(ctx, sdk.HexToAddress(data.StakingToken))
	}

	if data.Owner != "" {
		k.SetTokenManagerOwner(ctx, sdk.HexToAddress(data.Owner))
	}

	for _, v := range data.Records.StakingRecord {
		k.SetStakingVault(ctx, v.Validator, v.Delegator, v.EndTime, v.StorageTime, v.Amount)
	}
}


func ExportGenesis(ctx sdk.Context, k keeper.PreStakingKeeper) types.GenesisState {
	var records types.DelegationRecord
	sr := k.GetAllStakingVault(ctx)
	records.StakingRecord = sr
	return types.NewGenesisState(records, k.GetTokenManager(ctx), k.GetTokenManagerOwner(ctx))
}