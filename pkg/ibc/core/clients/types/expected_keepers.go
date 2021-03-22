package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	stakingtypes "github.com/ci123chain/ci123chain/pkg/staking/types"
	"time"
)

type StakingKeeper interface {
	GetHistoricalInfo(ctx sdk.Context, height int64) (stakingtypes.HistoricalInfo, bool)
	UnbondingTime(ctx sdk.Context) time.Duration
}
