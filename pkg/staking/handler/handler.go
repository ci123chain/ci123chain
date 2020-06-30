package handler

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/staking/keeper"
	"github.com/ci123chain/ci123chain/pkg/staking/types"
	staking "github.com/ci123chain/ci123chain/pkg/staking/types"
	gogotypes "github.com/gogo/protobuf/types"
	"time"
)

func NewHandler(k keeper.StakingKeeper) sdk.Handler {
	return func(ctx sdk.Context, tx sdk.Tx) sdk.Result {
		//ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch tx := tx.(type) {
		case *staking.CreateValidatorTx:
			return handleCreateValidatorTx(ctx, k, *tx)
		case *staking.DelegateTx:
			return handleDelegateTx(ctx, k, *tx)
		case *staking.RedelegateTx:
			return handleRedelegateTx(ctx, k, *tx)
		case *staking.UndelegateTx:
			return handleUndelegateTx(ctx, k, *tx)
		default:
			errMsg := fmt.Sprintf("unrecognized supply message type: %T", tx)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleCreateValidatorTx(ctx sdk.Context, k keeper.StakingKeeper, tx staking.CreateValidatorTx) sdk.Result {
	//
	if _, found := k.GetValidator(ctx, tx.ValidatorAddress); found {
		return types.ErrValidatorExisted(types.DefaultCodespace, nil).Result()
	}
	pk := tx.PublicKey

	if _, found := k.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(pk)); found {
		return types.ErrValidatorExisted(types.DefaultCodespace, nil).Result()
	}

	if _, err := tx.Description.EnsureLength(); err != nil {
		return types.ErrDescriptionOutOfLength(types.DefaultCodespace, err).Result()
	}

	//--------------
	//if ctx.ConsensusParams() != nil {
		//
		//tmPubKey := tmtypes.TM2PB.PubKey(pk)
		//fmt.Println(tmPubKey)
		/*if !tmstrings.StringInSlice(tmPubKey.Type, ctx.ConsensusParams().Validator.PubKeyTypes) {
			//
		}*/
	//}

	validator := staking.NewValidator(tx.ValidatorAddress, tx.PublicKey, tx.Description)
	commission := staking.NewCommissionWithTime(tx.Commission.Rate,
		tx.Commission.MaxRate, tx.Commission.MaxChangeRate, ctx.BlockHeader().Time)

	validator, err := validator.SetInitialCommission(commission)
	if err != nil {
		return types.ErrSetCommissionFailed(types.DefaultCodespace, err).Result()
	}
	validator.MinSelfDelegation = tx.MinSelfDelegation

	err = k.SetValidator(ctx, validator)
	if err != nil {
		return types.ErrSetValidatorFailed(types.DefaultCodespace, err).Result()
	}
	k.SetValidatorByConsAddr(ctx, validator)
	k.SetNewValidatorByPowerIndex(ctx, validator)

	k.AfterValidatorCreated(ctx, validator.OperatorAddress)

	_, err = k.Delegate(ctx, tx.DelegatorAddress, tx.Value.Amount, sdk.Unbonded, validator, true)
	if err != nil {
		return types.ErrDelegateFailed(types.DefaultCodespace, err).Result()
	}

	em := sdk.NewEventManager()
	em.EmitEvents(sdk.Events{
		sdk.NewEvent(types.EventTypeCreateValidator,
			sdk.NewAttribute(types.AttributeKeyValidator, tx.ValidatorAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, tx.Value.Amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, tx.DelegatorAddress.String()),
			),
	})

	return sdk.Result{Events: em.Events()}
}

func handleDelegateTx(ctx sdk.Context, k keeper.StakingKeeper, tx staking.DelegateTx) sdk.Result {
	//
	validator, found := k.GetValidator(ctx, tx.ValidatorAddress)
	if !found {
		return types.ErrNoExpectedValidator(types.DefaultCodespace, nil).Result()
	}
	denom := k.BondDenom(ctx)
	if tx.Amount.Denom != denom {
		return types.ErrBondedDenomDiff(types.DefaultCodespace, nil).Result()
	}
	_, err := k.Delegate(ctx, tx.DelegatorAddress, tx.Amount.Amount, sdk.Unbonded, validator, true)
	if err != nil {
		return types.ErrDelegateFailed(types.ModuleName, err).Result()
	}
	em := sdk.NewEventManager()
	em.EmitEvents(sdk.Events{
		sdk.NewEvent(types.EventTypeDelegate,
			sdk.NewAttribute(types.AttributeKeyValidator, tx.ValidatorAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, tx.Amount.Amount.String()),
			),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, tx.DelegatorAddress.String()),
			),
	})

	return sdk.Result{Events: em.Events()}
}

func handleRedelegateTx(ctx sdk.Context, k keeper.StakingKeeper, tx staking.RedelegateTx) sdk.Result {
	shares, err := k.ValidateUnbondAmount(ctx, tx.DelegatorAddress, tx.ValidatorSrcAddress, tx.Amount.Amount)
	if err != nil {
		return types.ErrValidateUnBondAmountFailed(types.DefaultCodespace, err).Result()
	}
	if tx.Amount.Denom != k.BondDenom(ctx) {
		return types.ErrBondedDenomDiff(types.ModuleName, nil).Result()
	}

	completionTime, err := k.Redelegate(ctx, tx.DelegatorAddress, tx.ValidatorSrcAddress, tx.ValidatorDstAddress, shares)
	if err != nil {
		return types.ErrRedelegationFailed(types.ModuleName, err).Result()
	}

	ts, err := gogotypes.TimestampProto(completionTime)
	if err != nil {
		return types.ErrGotTimeFailed(types.ModuleName, err).Result()
	}

	completionTimeBz := types.StakingCodec.MustMarshalBinaryLengthPrefixed(ts)

	em := sdk.NewEventManager()
	em.EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRedelegate,
			sdk.NewAttribute(types.AttributeKeySrcValidator, tx.ValidatorSrcAddress.String()),
			sdk.NewAttribute(types.AttributeKeyDstValidator, tx.ValidatorDstAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, tx.Amount.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyCompletionTime, completionTime.Format(time.RFC3339)),
			),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, tx.DelegatorAddress.String()),
			),
	})

	return sdk.Result{Data: completionTimeBz, Events: em.Events()}
}

func handleUndelegateTx(ctx sdk.Context, k keeper.StakingKeeper, tx staking.UndelegateTx) sdk.Result {
	//
	shares, err := k.ValidateUnbondAmount(ctx, tx.DelegatorAddress, tx.ValidatorAddress, tx.Amount.Amount)
	if err != nil {
		return types.ErrValidateUnBondAmountFailed(types.DefaultCodespace, err).Result()
	}
	if tx.Amount.Denom != k.BondDenom(ctx) {
		return types.ErrBondedDenomDiff(types.ModuleName, nil).Result()
	}
	completionTime, err := k.Undelegate(ctx, tx.DelegatorAddress, tx.ValidatorAddress, shares)
	if err != nil {
		//
		return types.ErrUndelegateFailed(types.DefaultCodespace, err).Result()
	}
	ts, err := gogotypes.TimestampProto(completionTime)
	if err != nil {
		return types.ErrGotTimeFailed(types.DefaultCodespace, err).Result()
	}
	completionTimeBz := types.StakingCodec.MustMarshalBinaryLengthPrefixed(ts)
	em := sdk.NewEventManager()
	em.EmitEvents(
		sdk.Events{
			sdk.NewEvent(types.EventTypeUnbond,
				sdk.NewAttribute(types.AttributeKeyValidator, tx.ValidatorAddress.String()),
				sdk.NewAttribute(sdk.AttributeKeyAmount, tx.Amount.Amount.String()),
				sdk.NewAttribute(types.AttributeKeyCompletionTime, completionTime.Format(time.RFC3339)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, tx.DelegatorAddress.String()),
			),
		})

	return sdk.Result{Data: completionTimeBz, Events: em.Events()}
}