package upgrade

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/upgrade/types"
)

// NewSoftwareUpgradeProposalHandler creates a governance handler to manage new proposal types.
// It enables SoftwareUpgradeProposal to propose an Upgrade, and CancelSoftwareUpgradeProposal
// to abort a previously voted upgrade.
func NewSoftwareUpgradeProposalHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch c := msg.(type) {
		case *types.SoftwareUpgradeProposal:
			return handleSoftwareUpgradeProposal(ctx, k, *c)

		case *types.CancelSoftwareUpgradeProposal:
			return handleCancelSoftwareUpgradeProposal(ctx, k, *c)

		default:
			return nil, nil
		}
	}
}

func handleSoftwareUpgradeProposal(ctx sdk.Context, k Keeper, p types.SoftwareUpgradeProposal) (*sdk.Result, error) {
	if p.Proposer.String() != InnerAccount {
		return nil, types.ErrPermission
	}
	err := k.ScheduleUpgrade(ctx, p.Plan)
	return &sdk.Result{}, err
}

func handleCancelSoftwareUpgradeProposal(ctx sdk.Context, k Keeper, p types.CancelSoftwareUpgradeProposal) (*sdk.Result, error) {
	if p.Proposer.String() != InnerAccount {
		return nil, types.ErrPermission
	}

	k.ClearUpgradePlan(ctx)
	return &sdk.Result{}, nil
}

