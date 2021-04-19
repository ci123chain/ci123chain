package keeper

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/distribution/types"
	stakingtypes "github.com/ci123chain/ci123chain/pkg/staking/types"
)


type Hooks struct {
	k   DistrKeeper
}

var _ stakingtypes.StakingHooks = Hooks{}

// Create new distribution hooks
func (k DistrKeeper) Hooks() Hooks { return Hooks{k} }



// initialize validator distribution record
func (h Hooks) AfterValidatorCreated(ctx sdk.Context, valAddr sdk.AccAddress) {
	val := h.k.StakingKeeper.Validator(ctx, valAddr)
	h.k.initializeValidator(ctx, val)
}

// cleanup for after validator is removed
func (h Hooks) AfterValidatorRemoved(ctx sdk.Context, _ sdk.AccAddress, valAddr sdk.AccAddress) {
	// fetch outstanding
	outstanding := h.k.GetValidatorOutstandingRewardsCoins(ctx, valAddr)

	// force-withdraw commission
	commission := h.k.GetValidatorAccumulatedCommission(ctx, valAddr).Commission
	if !commission.IsZero() {
		// subtract from outstanding
		outstanding = outstanding.Sub(commission)

		// split into integral & remainder
		coins, remainder := commission.TruncateDecimal()

		// remainder to community pool
		feePool := h.k.GetFeePool(ctx)
		feePool.CommunityPool = feePool.CommunityPool.Add(remainder)
		h.k.SetFeePool(ctx, feePool)

		// add to validator account
		if !coins.IsZero() {

			accAddr := valAddr
			withdrawAddr := h.k.GetDelegatorWithdrawAddr(ctx, accAddr)
			err := h.k.SupplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, withdrawAddr, sdk.NewCoins(coins))
			if err != nil {
				panic(err)
			}
		}
	}

	// add outstanding to community pool
	feePool := h.k.GetFeePool(ctx)
	feePool.CommunityPool = feePool.CommunityPool.Add(outstanding)
	h.k.SetFeePool(ctx, feePool)

	// delete outstanding
	h.k.DeleteValidatorOutstandingRewards(ctx, valAddr)

	// remove commission record
	h.k.DeleteValidatorAccumulatedCommission(ctx, valAddr)

	// clear slashes
	h.k.DeleteValidatorSlashEvents(ctx, valAddr)

	// clear historical rewards
	h.k.DeleteValidatorHistoricalRewards(ctx, valAddr)

	// clear current rewards
	h.k.DeleteValidatorCurrentRewards(ctx, valAddr)
}

// increment period
func (h Hooks) BeforeDelegationCreated(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.AccAddress) {
	val := h.k.StakingKeeper.Validator(ctx, valAddr)
	h.k.incrementValidatorPeriod(ctx, val)
}

// withdraw delegation rewards (which also increments period)
func (h Hooks) BeforeDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.AccAddress) {
	val := h.k.StakingKeeper.Validator(ctx, valAddr)
	del := h.k.StakingKeeper.Delegation(ctx, delAddr, valAddr)
	if _, err := h.k.withdrawDelegationRewards(ctx, val, del); err != nil {
		panic(err)
	}
}

// create new delegation period record
func (h Hooks) AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.AccAddress) {
	h.k.initializeDelegation(ctx, valAddr, delAddr)
}

// record the slash event
func (h Hooks) BeforeValidatorSlashed(ctx sdk.Context, valAddr sdk.AccAddress, fraction sdk.Dec) {
	//h.k.updateValidatorSlashFraction(ctx, valAddr, fraction)
}

// nolint - unused hooks
func (h Hooks) BeforeValidatorModified(_ sdk.Context, _ sdk.AccAddress)                         {}
func (h Hooks) AfterValidatorBonded(_ sdk.Context, _ sdk.AccAddress, _ sdk.AccAddress)         {}
func (h Hooks) AfterValidatorBeginUnbonding(_ sdk.Context, _ sdk.AccAddress, _ sdk.AccAddress) {}
func (h Hooks) BeforeDelegationRemoved(_ sdk.Context, _ sdk.AccAddress, _ sdk.AccAddress)       {}