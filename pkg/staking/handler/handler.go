package handler

import (
	"errors"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/staking/keeper"
	"github.com/ci123chain/ci123chain/pkg/staking/types"
	staking "github.com/ci123chain/ci123chain/pkg/staking/types"
	"github.com/ci123chain/ci123chain/pkg/util"
	gogotypes "github.com/gogo/protobuf/types"
	"time"
)

func NewHandler(k keeper.StakingKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *staking.MsgCreateValidator:
			return handleMsgCreateValidator(ctx, k, *msg)
		case *staking.MsgEditValidator:
			return handleMsgEditValidator(ctx, k, *msg)
		case *staking.MsgDelegate:
			return handleMsgDelegate(ctx, k, *msg)
		case *staking.MsgRedelegate:
			return handleMsgRedelegate(ctx, k, *msg)
		case *staking.MsgUndelegate:
			return handleMsgUndelegate(ctx, k, *msg)
		default:
			errMsg := fmt.Sprintf("unrecognized supply message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgCreateValidator(ctx sdk.Context, k keeper.StakingKeeper, msg staking.MsgCreateValidator) sdk.Result {
	if _, found := k.GetValidator(ctx, msg.ValidatorAddress); found {
		return types.ErrValidatorExisted(types.DefaultCodespace, errors.New(fmt.Sprintf("validator %s has existed", msg.ValidatorAddress.String()))).Result()
	}
	pk, err := util.ParsePubKey(msg.PublicKey)
	if err != nil {
		return types.ErrCheckParams(types.DefaultCodespace, "public_key").Result()
	}

	if _, found := k.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(pk)); found {
		return types.ErrValidatorExisted(types.DefaultCodespace, errors.New(fmt.Sprintf("the pubKey has been bonded"))).Result()
	}

	if _, err := msg.Description.EnsureLength(); err != nil {
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

	validator, err := staking.NewValidator(msg.ValidatorAddress, msg.PublicKey, msg.Description)
	if err != nil {
		return types.ErrSetValidatorFailed(types.DefaultCodespace, err).Result()
	}
	commission := staking.NewCommissionWithTime(msg.Commission.Rate,
		msg.Commission.MaxRate, msg.Commission.MaxChangeRate, ctx.BlockHeader().Time)

	validator, err = validator.SetInitialCommission(commission)
	if err != nil {
		return types.ErrSetCommissionFailed(types.DefaultCodespace, err).Result()
	}
	validator.MinSelfDelegation = msg.MinSelfDelegation

	err = k.SetValidator(ctx, validator)
	if err != nil {
		return types.ErrSetValidatorFailed(types.DefaultCodespace, err).Result()
	}
	k.SetValidatorByConsAddr(ctx, validator)
	k.SetNewValidatorByPowerIndex(ctx, validator)

	k.AfterValidatorCreated(ctx, validator.OperatorAddress)

	_, err = k.Delegate(ctx, msg.DelegatorAddress, msg.Value.Amount, sdk.Unbonded, validator, true)
	if err != nil {
		return types.ErrDelegateFailed(types.DefaultCodespace, err).Result()
	}

	em := sdk.NewEventManager()
	em.EmitEvents(sdk.Events{
		sdk.NewEvent(types.EventTypeCreateValidator,
			sdk.NewAttribute([]byte(sdk.AttributeKeyMethod), []byte(types.EventTypeCreateValidator)),
			sdk.NewAttribute([]byte(types.AttributeKeyValidator), []byte(msg.ValidatorAddress.String())),
			sdk.NewAttribute([]byte(sdk.AttributeKeyAmount), []byte(msg.Value.Amount.String())),
			sdk.NewAttribute([]byte(sdk.AttributeKeyModule), []byte(types.AttributeValueCategory)),
			sdk.NewAttribute([]byte(sdk.AttributeKeySender), []byte(msg.FromAddress.String())),
		),
		//sdk.NewEvent(
		//	sdk.EventTypeMessage,
		//	sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		//	sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAddress.String()),
		//	),
	})

	return sdk.Result{Events: em.Events()}
}

func handleMsgEditValidator(ctx sdk.Context, k keeper.StakingKeeper, msg staking.MsgEditValidator) sdk.Result {

	validator, found := k.GetValidator(ctx, msg.ValidatorAddress)
	if !found {
		return types.ErrNoExpectedValidator(types.DefaultCodespace, errors.New(fmt.Sprintf("%s not found", validator.OperatorAddress.String()))).Result()
	}
	description, err := validator.Description.UpdateDescription(msg.Description)
	if err != nil {
		return types.ErrCheckParams(types.DefaultCodespace, err.Error()).Result()
	}

	validator.Description = description

	//if commissionRate == nil, no change.
	if msg.CommissionRate != nil {
		commission, err := k.UpdateValidatorCommission(ctx, validator, *msg.CommissionRate)
		if err != nil {
			return types.ErrCheckParams(types.DefaultCodespace, err.Error()).Result()
		}

		k.BeforeValidatorModified(ctx, msg.ValidatorAddress)
		validator.Commission = commission
	}
	//if minSelfDelegation == nil, no change.
	if msg.MinSelfDelegation != nil {
		if !msg.MinSelfDelegation.GT(validator.MinSelfDelegation) {
			return types.ErrCheckParams(types.DefaultCodespace, fmt.Sprintf("new minSelfDelegation should be greater than before: %s", validator.MinSelfDelegation.String())).Result()
		}
		if msg.MinSelfDelegation.GT(validator.Tokens) {
			return types.ErrCheckParams(types.DefaultCodespace, fmt.Sprintf("new minSelfDelegation can not be greater than tokens that you hold on")).Result()
		}
		validator.MinSelfDelegation = *msg.MinSelfDelegation
	}
	err = k.SetValidator(ctx, validator)
	if err != nil {
		return types.ErrSetValidatorFailed(types.DefaultCodespace, err).Result()
	}
	em := sdk.NewEventManager()
	em.EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeEditValidator,
			sdk.NewAttribute([]byte(sdk.AttributeKeyMethod), []byte(types.EventTypeEditValidator)),
			sdk.NewAttribute([]byte(types.AttributeKeyCommissionRate), []byte(validator.Commission.String())),
			sdk.NewAttribute([]byte(types.AttributeKeyMinSelfDelegation), []byte(validator.MinSelfDelegation.String())),
			sdk.NewAttribute([]byte(sdk.AttributeKeyModule), []byte(types.AttributeValueCategory)),
			sdk.NewAttribute([]byte(sdk.AttributeKeySender), []byte(msg.FromAddress.String())),
		),
		//sdk.NewEvent(
		//	sdk.EventTypeMessage,
		//	sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		//	sdk.NewAttribute(sdk.AttributeKeySender, msg.ValidatorAddress.String()),
		//),
	})

	return sdk.Result{Events: em.Events()}
}

func handleMsgDelegate(ctx sdk.Context, k keeper.StakingKeeper, msg staking.MsgDelegate) sdk.Result {
	//
	validator, found := k.GetValidator(ctx, msg.ValidatorAddress)
	if !found {
		return types.ErrNoExpectedValidator(types.DefaultCodespace, nil).Result()
	}
	denom := k.BondDenom(ctx)
	if msg.Amount.Denom != denom {
		return types.ErrBondedDenomDiff(types.DefaultCodespace, nil).Result()
	}
	_, err := k.Delegate(ctx, msg.DelegatorAddress, msg.Amount.Amount, sdk.Unbonded, validator, true)
	if err != nil {
		return types.ErrDelegateFailed(types.ModuleName, err).Result()
	}
	em := sdk.NewEventManager()
	em.EmitEvents(sdk.Events{
		sdk.NewEvent(types.EventTypeDelegate,
			sdk.NewAttribute([]byte(sdk.AttributeKeyMethod), []byte(types.EventTypeDelegate)),
			sdk.NewAttribute([]byte(types.AttributeKeyValidator), []byte(msg.ValidatorAddress.String())),
			sdk.NewAttribute([]byte(sdk.AttributeKeyAmount), []byte(msg.Amount.Amount.String())),
			sdk.NewAttribute([]byte(sdk.AttributeKeyModule), []byte(types.AttributeValueCategory)),
			sdk.NewAttribute([]byte(sdk.AttributeKeySender), []byte(msg.FromAddress.String())),
			),
		//sdk.NewEvent(
		//	sdk.EventTypeMessage,
		//	sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		//	sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAddress.String()),
		//	),
	})

	return sdk.Result{Events: em.Events()}
}

func handleMsgRedelegate(ctx sdk.Context, k keeper.StakingKeeper, msg staking.MsgRedelegate) sdk.Result {
	shares, err := k.ValidateUnbondAmount(ctx, msg.DelegatorAddress, msg.ValidatorSrcAddress, msg.Amount.Amount)
	if err != nil {
		return types.ErrValidateUnBondAmountFailed(types.DefaultCodespace, err).Result()
	}
	if msg.Amount.Denom != k.BondDenom(ctx) {
		return types.ErrBondedDenomDiff(types.ModuleName, nil).Result()
	}

	completionTime, err := k.Redelegate(ctx, msg.DelegatorAddress, msg.ValidatorSrcAddress, msg.ValidatorDstAddress, shares)
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
			sdk.NewAttribute([]byte(sdk.AttributeKeyMethod), []byte(types.EventTypeRedelegate)),
			sdk.NewAttribute([]byte(types.AttributeKeySrcValidator), []byte(msg.ValidatorSrcAddress.String())),
			sdk.NewAttribute([]byte(types.AttributeKeyDstValidator), []byte(msg.ValidatorDstAddress.String())),
			sdk.NewAttribute([]byte(sdk.AttributeKeyAmount), []byte(msg.Amount.Amount.String())),
			sdk.NewAttribute([]byte(types.AttributeKeyCompletionTime), []byte(completionTime.Format(time.RFC3339))),
			sdk.NewAttribute([]byte(sdk.AttributeKeyModule), []byte(types.AttributeValueCategory)),
			sdk.NewAttribute([]byte(sdk.AttributeKeySender), []byte(msg.FromAddress.String())),
			),
		//sdk.NewEvent(
		//	sdk.EventTypeMessage,
		//	sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		//	sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAddress.String()),
		//	),
	})

	return sdk.Result{Data: completionTimeBz, Events: em.Events()}
}

func handleMsgUndelegate(ctx sdk.Context, k keeper.StakingKeeper, msg staking.MsgUndelegate) sdk.Result {
	//
	shares, err := k.ValidateUnbondAmount(ctx, msg.DelegatorAddress, msg.ValidatorAddress, msg.Amount.Amount)
	if err != nil {
		return types.ErrValidateUnBondAmountFailed(types.DefaultCodespace, err).Result()
	}
	if msg.Amount.Denom != k.BondDenom(ctx) {
		return types.ErrBondedDenomDiff(types.ModuleName, nil).Result()
	}
	completionTime, err := k.Undelegate(ctx, msg.DelegatorAddress, msg.ValidatorAddress, shares)
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
			sdk.NewEvent(types.EventTypeUndelegate,
				sdk.NewAttribute([]byte(sdk.AttributeKeyMethod), []byte(types.EventTypeUndelegate)),
				sdk.NewAttribute([]byte(types.AttributeKeyValidator), []byte(msg.ValidatorAddress.String())),
				sdk.NewAttribute([]byte(sdk.AttributeKeyAmount), []byte(msg.Amount.Amount.String())),
				sdk.NewAttribute([]byte(types.AttributeKeyCompletionTime), []byte(completionTime.Format(time.RFC3339))),
				sdk.NewAttribute([]byte(sdk.AttributeKeyModule), []byte(types.AttributeValueCategory)),
				sdk.NewAttribute([]byte(sdk.AttributeKeySender), []byte(msg.FromAddress.String())),
		),
		//sdk.NewEvent(
		//	sdk.EventTypeMessage,
		//	sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		//	sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAddress.String()),
		//	),
		})

	return sdk.Result{Data: completionTimeBz, Events: em.Events()}
}