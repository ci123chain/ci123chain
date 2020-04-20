package keeper

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/tendermint/tendermint/libs/log")

func Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "ibc")
}
