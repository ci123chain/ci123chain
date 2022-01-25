package pre_staking

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/pre_staking/keeper"
)

func BeginBlock() {}

func EndBlock(ctx sdk.Context, k keeper.PreStakingKeeper) {
	k.UpdateDeadlineRecord(ctx)
}
