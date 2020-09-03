package keeper

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/staking/types"
)

// Implements StakingHooks interface
var _ types.StakingHooks = StakingKeeper{}

func (k StakingKeeper) AfterValidatorCreated(ctx sdk.Context, valAddr sdk.AccAddress) {
	if k.hooks != nil {
		k.hooks.AfterValidatorCreated(ctx, valAddr)

	}
}

// AfterValidatorRemoved - call hook if registered
func (k StakingKeeper) AfterValidatorRemoved(ctx sdk.Context, consAddr sdk.AccAddress, valAddr sdk.AccAddress) {
	if k.hooks != nil {
		k.hooks.AfterValidatorRemoved(ctx, consAddr, valAddr)
	}
}

// BeforeValidatorModified - call hook if registered
func (k StakingKeeper) BeforeValidatorModified(ctx sdk.Context, valAddr sdk.AccAddress) {
	if k.hooks != nil {
		k.hooks.BeforeValidatorModified(ctx, valAddr)
	}
}

// BeforeDelegationCreated - call hook if registered
func (k StakingKeeper) BeforeDelegationCreated(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.AccAddress) {
	if k.hooks != nil {
		k.hooks.BeforeDelegationCreated(ctx, delAddr, valAddr)
	}
}


func (k StakingKeeper) BeforeDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.AccAddress) {
	if k.hooks != nil {
		k.hooks.BeforeDelegationSharesModified(ctx, delAddr, valAddr)
	}
}

// BeforeDelegationRemoved - call hook if registered
func (k StakingKeeper) BeforeDelegationRemoved(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.AccAddress) {
	if k.hooks != nil {
		k.hooks.BeforeDelegationRemoved(ctx, delAddr, valAddr)
	}
}

// AfterDelegationModified - call hook if registered
func (k StakingKeeper) AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.AccAddress) {
	if k.hooks != nil {
		k.hooks.AfterDelegationModified(ctx, delAddr, valAddr)
	}
}

// BeforeValidatorSlashed - call hook if registered
func (k StakingKeeper) BeforeValidatorSlashed(ctx sdk.Context, valAddr sdk.AccAddress, fraction sdk.Dec) {
	if k.hooks != nil {
		k.hooks.BeforeValidatorSlashed(ctx, valAddr, fraction)
	}
}

// AfterValidatorBonded - call hook if registered
func (k StakingKeeper) AfterValidatorBonded(ctx sdk.Context, consAddr sdk.AccAddress, valAddr sdk.AccAddress) {
	if k.hooks != nil {
		k.hooks.AfterValidatorBonded(ctx, consAddr, valAddr)
	}
}

// AfterValidatorBeginUnbonding - call hook if registered
func (k StakingKeeper) AfterValidatorBeginUnbonding(ctx sdk.Context, consAddr sdk.AccAddress, valAddr sdk.AccAddress) {
	if k.hooks != nil {
		k.hooks.AfterValidatorBeginUnbonding(ctx, consAddr, valAddr)
	}
}