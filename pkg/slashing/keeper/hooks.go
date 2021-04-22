package keeper

import (
	"time"

	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/slashing/types"
)

func (k Keeper) AfterValidatorBonded(ctx sdk.Context, address sdk.AccAddress, _ sdk.AccAddress) {
	// Update the signing info start height or create a new signing info
	_, found := k.GetValidatorSigningInfo(ctx, address)
	if !found {
		signingInfo := types.NewValidatorSigningInfo(
			address,
			ctx.BlockHeight(),
			0,
			time.Unix(0, 0),
			false,
			0,
		)
		k.SetValidatorSigningInfo(ctx, address, signingInfo)
	}
}

// AfterValidatorCreated adds the address-pubkey relation when a validator is created.
func (k Keeper) AfterValidatorCreated(ctx sdk.Context, valAddr sdk.AccAddress) error {
	validator := k.sk.Validator(ctx, valAddr)
	consPk := validator.GetConsPubKey()
	k.AddPubkey(ctx, consPk)

	return nil
}

// AfterValidatorRemoved deletes the address-pubkey relation when a validator is removed,
func (k Keeper) AfterValidatorRemoved(ctx sdk.Context, address sdk.AccAddress) {
	k.deleteAddrPubkeyRelation(ctx, address)
}

//_________________________________________________________________________________________

// Hooks wrapper struct for slashing keeper
type Hooks struct {
	k Keeper
}

var _ types.StakingHooks = Hooks{}

// Return the wrapper struct
func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

// Implements sdk.ValidatorHooks
func (h Hooks) AfterValidatorBonded(ctx sdk.Context, consAddr sdk.AccAddress, valAddr sdk.AccAddress) {
	h.k.AfterValidatorBonded(ctx, consAddr, valAddr)
}

// Implements sdk.ValidatorHooks
func (h Hooks) AfterValidatorRemoved(ctx sdk.Context, consAddr sdk.AccAddress, _ sdk.AccAddress) {
	h.k.AfterValidatorRemoved(ctx, consAddr)
}

// Implements sdk.ValidatorHooks
func (h Hooks) AfterValidatorCreated(ctx sdk.Context, valAddr sdk.AccAddress) {
	h.k.AfterValidatorCreated(ctx, valAddr)
}

func (h Hooks) AfterValidatorBeginUnbonding(_ sdk.Context, _ sdk.AccAddress, _ sdk.AccAddress)  {}
func (h Hooks) BeforeValidatorModified(_ sdk.Context, _ sdk.AccAddress)                          {}
func (h Hooks) BeforeDelegationCreated(_ sdk.Context, _ sdk.AccAddress, _ sdk.AccAddress)        {}
func (h Hooks) BeforeDelegationSharesModified(_ sdk.Context, _ sdk.AccAddress, _ sdk.AccAddress) {}
func (h Hooks) BeforeDelegationRemoved(_ sdk.Context, _ sdk.AccAddress, _ sdk.AccAddress)        {}
func (h Hooks) AfterDelegationModified(_ sdk.Context, _ sdk.AccAddress, _ sdk.AccAddress)        {}
func (h Hooks) BeforeValidatorSlashed(_ sdk.Context, _ sdk.AccAddress, _ sdk.Dec)                {}
