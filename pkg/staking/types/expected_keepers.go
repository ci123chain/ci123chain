package types

import sdk "github.com/ci123chain/ci123chain/pkg/abci/types"

type StakingHooks interface {
	AferValidatorCreated(ctx sdk.Context, valAddr sdk.AccAddress)
	BeforeValidatorModified(ctx sdk.Context, valAddr sdk.AccAddress)
	AfterValidatorRemoved(ctx sdk.Context, consAddr sdk.AccAddress, valAddr sdk.AccAddress)

	AfterValidatorBonded(ctx sdk.Context, consAddr sdk.AccAddress, valAddr sdk.AccAddress)
	AfterValidatorBeginUnbonding(ctx sdk.Context, consAddr sdk.AccAddress, valAddr sdk.AccAddress) // Must be called when a validator begins unbonding

	BeforeDelegationCreated(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.AccAddress)        // Must be called when a delegation is created
	BeforeDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.AccAddress) // Must be called when a delegation's shares are modified
	BeforeDelegationRemoved(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.AccAddress)        // Must be called when a delegation is removed
	AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.AccAddress)
	BeforeValidatorSlashed(ctx sdk.Context, valAddr sdk.AccAddress, fraction sdk.Dec)
}
