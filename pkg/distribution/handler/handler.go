package handler

import (
	"errors"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/distribution/keeper"
	"github.com/ci123chain/ci123chain/pkg/distribution/types"
)


func NewHandler(k keeper.DistrKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgSetWithdrawAddress:
			return handleMsgModifyWithdrawAddress(ctx, *msg, k)
		case *types.MsgWithdrawDelegatorReward:
			return handleMsgWithdrawDelegatorReward(ctx, *msg, k)
		case *types.MsgWithdrawValidatorCommission:
			return handleMsgWithdrawValidatorCommission(ctx, *msg, k)
		case *types.MsgFundCommunityPool:
			return handleMsgFundCommunityPool(ctx, *msg, k)
		default:
			errMsg := fmt.Sprintf("unrecognized supply message types: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}


func handleMsgModifyWithdrawAddress(ctx sdk.Context, msg types.MsgSetWithdrawAddress, k keeper.DistrKeeper) sdk.Result {

	//verify identity
	if !msg.FromAddress.Equal(msg.DelegatorAddress) {
		return types.ErrWithdrawAddressInfoMismatch(types.DefaultCodespace, msg.FromAddress, msg.DelegatorAddress).Result()
	}
	//check validator that is bonded to delegator account.
	validators, Err := k.StakingKeeper.GetDelegatorValidators(ctx, msg.DelegatorAddress, 3)
	if Err != nil || validators == nil {
		return types.ErrBadAddress(types.DefaultCodespace, errors.New(fmt.Sprintf("got no validator that is bonded to %s", msg.DelegatorAddress.String()))).Result()
	}

	err := k.SetWithdrawAddr(ctx, msg.DelegatorAddress, msg.WithdrawAddress)
	if err != nil {
		return types.ErrHandleTxFailed(types.DefaultCodespace, err).Result()
	}

	//ctx.EventManager().EmitEvent(
	//	sdk.NewEvent(
	//		types.EventTypeModifyWithdrawAddress,
	//		sdk.NewAttribute([]byte(sdk.AttributeKeyMethod), []byte(types.EventTypeModifyWithdrawAddress)),
	//		sdk.NewAttribute([]byte(sdk.AttributeKeyModule), []byte(types.AttributeValueCategory)),
	//		sdk.NewAttribute([]byte(sdk.AttributeKeySender), []byte(msg.FromAddress.String())),
	//	),
	//)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgWithdrawDelegatorReward(ctx sdk.Context, msg types.MsgWithdrawDelegatorReward, k keeper.DistrKeeper) sdk.Result {

	//verify identity
	if !msg.FromAddress.Equal(msg.DelegatorAddress) {
		return types.ErrWithdrawAddressInfoMismatch(types.DefaultCodespace, msg.FromAddress, msg.DelegatorAddress).Result()
	}

	_, err := k.WithdrawDelegationRewards(ctx, msg.DelegatorAddress, msg.ValidatorAddress)
	if err != nil {
		return types.ErrHandleTxFailed(types.DefaultCodespace, err).Result()
	}

	//ctx.EventManager().EmitEvent(
	//	sdk.NewEvent(
	//		types.EventTypeWithdrawRewards,
	//		sdk.NewAttribute([]byte(sdk.AttributeKeyMethod), []byte(types.EventTypeWithdrawRewards)),
	//		sdk.NewAttribute([]byte(sdk.AttributeKeyModule), []byte(types.AttributeValueCategory)),
	//		sdk.NewAttribute([]byte(types.AttributeKeyWithdrawAddress), []byte(msg.FromAddress.String())),
	//		sdk.NewAttribute([]byte(sdk.AttributeKeySender), []byte(msg.FromAddress.String())),
	//	),
	//)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgWithdrawValidatorCommission(ctx sdk.Context, msg types.MsgWithdrawValidatorCommission, k keeper.DistrKeeper) sdk.Result {
	//verify identity
	if !msg.FromAddress.Equal(msg.ValidatorAddress) {
		return types.ErrWithdrawAddressInfoMismatch(types.DefaultCodespace, msg.FromAddress, msg.ValidatorAddress).Result()
	}

	_, ok := k.StakingKeeper.GetValidator(ctx, msg.ValidatorAddress)
	if !ok {
		return types.ErrNoValidatorExist(types.DefaultCodespace, msg.ValidatorAddress.String()).Result()
	}

	_, err := k.WithdrawValidatorCommission(ctx, msg.ValidatorAddress)
	if err != nil {
		return types.ErrHandleTxFailed(types.DefaultCodespace, err).Result()
	}

	//ctx.EventManager().EmitEvent(
	//	sdk.NewEvent(
	//		types.EventTypeWithdrawCommission,
	//		sdk.NewAttribute([]byte(sdk.AttributeKeyMethod), []byte(types.EventTypeWithdrawCommission)),
	//		sdk.NewAttribute([]byte(sdk.AttributeKeyModule), []byte(types.AttributeValueCategory)),
	//		sdk.NewAttribute([]byte(types.AttributeKeyWithdrawAddress), []byte(msg.FromAddress.String())),
	//		sdk.NewAttribute([]byte(sdk.AttributeKeySender), []byte(msg.FromAddress.String())),
	//	),
	//)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

//send funds from personal account to communityPool.
func handleMsgFundCommunityPool(ctx sdk.Context, msg types.MsgFundCommunityPool, k keeper.DistrKeeper) sdk.Result {
	if err := k.FundCommunityPool(ctx, msg.Amount, msg.Depositor); err != nil {
		return types.ErrHandleTxFailed(types.DefaultCodespace, err).Result()
	}

	//ctx.EventManager().EmitEvent(
	//	sdk.NewEvent(
	//		types.EventTypeFundCommunityPool,
	//		sdk.NewAttribute([]byte(sdk.AttributeKeyMethod), []byte(types.EventTypeFundCommunityPool)),
	//		sdk.NewAttribute([]byte(sdk.AttributeKeyModule), []byte(types.AttributeValueCategory)),
	//		sdk.NewAttribute([]byte(sdk.AttributeKeySender), []byte(msg.FromAddress.String())),
	//	),
	//)

	return sdk.Result{Events: ctx.EventManager().Events()}
}