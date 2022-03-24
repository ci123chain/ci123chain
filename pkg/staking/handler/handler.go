package handler

import (
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
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		//case *staking.MsgCreateValidator:
		//	return handleMsgCreateValidator(ctx, k, *msg)
		case *staking.MsgEditValidator:
			return handleMsgEditValidator(ctx, k, *msg)
		//case *staking.MsgDelegate:
		//	return handleMsgDelegate(ctx, k, *msg)
		//case *staking.MsgRedelegate:
		//	return handleMsgRedelegate(ctx, k, *msg)
		//case *staking.MsgUndelegate:
		//	return handleMsgUndelegate(ctx, k, *msg)
		//default:
		//	return nil, types.ErrInvalidTxType
		}
		return nil, nil
	}
}

func handleMsgCreateValidator(ctx sdk.Context, k keeper.StakingKeeper, msg staking.MsgCreateValidator) (*sdk.Result, error) {
	if _, found := k.GetValidator(ctx, msg.ValidatorAddress); found {
		return nil, types.ErrNoExpectedValidator
	}
	pk, err := util.ParsePubKey(msg.PublicKey)
	if err != nil {
		return nil, types.ErrInvalidPublicKey
	}

	if _, found := k.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(pk)); found {
		return nil, types.ErrPubkeyHasBonded
	}

	if _, err := msg.Description.EnsureLength(); err != nil {
		return nil, err
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

	validator, _ := staking.NewValidator(msg.ValidatorAddress, msg.PublicKey, msg.Description)

	commission := staking.NewCommissionWithTime(msg.Commission.Rate,
		msg.Commission.MaxRate, msg.Commission.MaxChangeRate, ctx.BlockHeader().Time)

	validator, err = validator.SetInitialCommission(commission)
	if err != nil {
		return nil, err
	}
	validator.MinSelfDelegation = msg.MinSelfDelegation

	err = k.SetValidator(ctx, validator)
	if err != nil {
		return nil, types.ErrSetValidatorFailed
	}
	k.SetValidatorByConsAddr(ctx, validator)
	k.SetNewValidatorByPowerIndex(ctx, validator)

	k.AfterValidatorCreated(ctx, validator.OperatorAddress)

	_, err = k.Delegate(ctx, msg.DelegatorAddress, msg.Value.Amount, sdk.Unbonded, validator, true)
	if err != nil {
		return nil, err
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

	return &sdk.Result{Events: em.Events()}, nil
}

func handleMsgEditValidator(ctx sdk.Context, k keeper.StakingKeeper, msg staking.MsgEditValidator) (*sdk.Result, error) {

	validator, found := k.GetValidator(ctx, msg.ValidatorAddress)
	if !found {
		_ = fmt.Sprintf("validator %s has existed", msg.ValidatorAddress.String())
		return nil, types.ErrNoExpectedValidator
	}
	description, err := validator.Description.UpdateDescription(msg.Description)
	if err != nil {
		return nil, err
	}

	validator.Description = description

	//if commissionRate == nil, no change.
	if msg.CommissionRate != nil {
		commission, err := k.UpdateValidatorCommission(ctx, validator, *msg.CommissionRate)
		if err != nil {
			return nil, err
		}

		k.BeforeValidatorModified(ctx, msg.ValidatorAddress)
		validator.Commission = commission
	}
	//if minSelfDelegation == nil, no change.
	if msg.MinSelfDelegation != nil {
		if !msg.MinSelfDelegation.GT(validator.MinSelfDelegation) {
			return nil, types.ErrInvalidParam
		}
		if msg.MinSelfDelegation.GT(validator.Tokens) {
			return nil, types.ErrInvalidParam
		}
		validator.MinSelfDelegation = *msg.MinSelfDelegation
	}
	err = k.SetValidator(ctx, validator)
	if err != nil {
		return nil, types.ErrSetValidatorFailed
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

	return &sdk.Result{Events: em.Events()}, nil
}

func handleMsgDelegate(ctx sdk.Context, k keeper.StakingKeeper, msg staking.MsgDelegate) (*sdk.Result, error) {
	//
	validator, found := k.GetValidator(ctx, msg.ValidatorAddress)
	if !found {
		//r := fmt.Sprintf("validator %s has existed", msg.ValidatorAddress.String())
		return nil, types.ErrNoExpectedValidator
	}
	denom := k.BondDenom(ctx)
	if msg.Amount.Denom != denom {
		return nil, types.ErrInvalidParam
	}
	_, err := k.Delegate(ctx, msg.DelegatorAddress, msg.Amount.Amount, sdk.Unbonded, validator, true)
	if err != nil {
		return nil, err
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

	return &sdk.Result{Events: em.Events()}, nil
}

func handleMsgRedelegate(ctx sdk.Context, k keeper.StakingKeeper, msg staking.MsgRedelegate) (*sdk.Result, error) {
	shares, err := k.ValidateUnbondAmount(ctx, msg.DelegatorAddress, msg.ValidatorSrcAddress, msg.Amount.Amount)
	if err != nil {
		return nil, err
	}
	if msg.Amount.Denom != k.BondDenom(ctx) {
		return nil, types.ErrInvalidParam
	}

	completionTime, err := k.Redelegate(ctx, msg.DelegatorAddress, msg.ValidatorSrcAddress, msg.ValidatorDstAddress, shares)
	if err != nil {
		return nil, err
	}

	ts, err := gogotypes.TimestampProto(completionTime)
	if err != nil {
		return nil, types.ErrInvalidParam
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

	return &sdk.Result{Data: completionTimeBz, Events: em.Events()}, nil
}

func handleMsgUndelegate(ctx sdk.Context, k keeper.StakingKeeper, msg staking.MsgUndelegate) (*sdk.Result, error) {
	//
	shares, err := k.ValidateUnbondAmount(ctx, msg.DelegatorAddress, msg.ValidatorAddress, msg.Amount.Amount)
	if err != nil {
		return nil, err
	}
	if msg.Amount.Denom != k.BondDenom(ctx) {
		return nil, types.ErrInvalidParam
	}
	completionTime, err := k.Undelegate(ctx, msg.DelegatorAddress, msg.ValidatorAddress, shares)
	if err != nil {
		return nil, err
	}
	ts, err := gogotypes.TimestampProto(completionTime)
	if err != nil {
		return nil, types.ErrInvalidParam
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

	return &sdk.Result{Data: completionTimeBz, Events: em.Events()}, nil
}