package pre_staking

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/pre_staking/keeper"
	"github.com/ci123chain/ci123chain/pkg/pre_staking/types"
)

const (
	moduleAcc = "0x3F43E75Aaba2c2fD6E227C10C6E7DC125A93DE3c"
)

func InitGenesis(ctx sdk.Context, k keeper.PreStakingKeeper, data GenesisState) {

	if data.DaoAddress != "" {
		k.SetWeeLinkDao(ctx, sdk.HexToAddress(data.DaoAddress))
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
	return types.NewGenesisState(records, k.GetWeeLinkDao(ctx))
}