package slashing

import (
	"time"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/ci123chain/ci123chain/pkg/telemetry"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/slashing/keeper"
	"github.com/ci123chain/ci123chain/pkg/slashing/types"
)

// BeginBlocker check for infraction evidence or downtime of validators
// on every begin block
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	// Iterate over all the validators which *should* have signed this block
	// store whether or not they have actually signed it and slash/unbond any
	// which have missed too many blocks in a row (downtime slashing)
	for _, voteInfo := range req.LastCommitInfo.GetVotes() {
		k.HandleValidatorSignature(ctx, sdk.ToAccAddress(voteInfo.Validator.Address), voteInfo.Validator.Power, voteInfo.SignedLastBlock)
	}
}
