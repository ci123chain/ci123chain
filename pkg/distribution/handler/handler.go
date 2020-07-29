package handler

import (
	"errors"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/distribution/keeper"
	"github.com/ci123chain/ci123chain/pkg/distribution/types"
)


func NewHandler(k keeper.DistrKeeper) sdk.Handler {
	return func(ctx sdk.Context, tx sdk.Tx) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch tx := tx.(type) {
		case *types.SetWithdrawAddressTx:
			return handleMsgModifyWithdrawAddress(ctx, *tx, k)
		case *types.WithdrawDelegatorRewardTx:
			return handleMsgWithdrawDelegatorReward(ctx, *tx, k)
		case *types.WithdrawValidatorCommissionTx:
			return handleMsgWithdrawValidatorCommission(ctx, *tx, k)
		case *types.FundCommunityPoolTx:
			return handleMsgFundCommunityPool(ctx, *tx, k)
		default:
			errMsg := fmt.Sprintf("unrecognized supply message type: %T", tx)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}


func handleMsgModifyWithdrawAddress(ctx sdk.Context, msg types.SetWithdrawAddressTx, k keeper.DistrKeeper) sdk.Result {

	//verify identity
	if !msg.From.Equal(msg.DelegatorAddress) {
		return types.ErrWithdrawAddressInfoMismatch(types.DefaultCodespace, msg.From, msg.DelegatorAddress).Result()
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

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAddress.String()),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgWithdrawDelegatorReward(ctx sdk.Context, msg types.WithdrawDelegatorRewardTx, k keeper.DistrKeeper) sdk.Result {

	//verify identity
	if !msg.From.Equal(msg.DelegatorAddress) {
		return types.ErrWithdrawAddressInfoMismatch(types.DefaultCodespace, msg.From, msg.DelegatorAddress).Result()
	}

	_, err := k.WithdrawDelegationRewards(ctx, msg.DelegatorAddress, msg.ValidatorAddress)
	if err != nil {
		return types.ErrHandleTxFailed(types.DefaultCodespace, err).Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAddress.String()),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgWithdrawValidatorCommission(ctx sdk.Context, msg types.WithdrawValidatorCommissionTx, k keeper.DistrKeeper) sdk.Result {
	//verify identity
	if !msg.From.Equal(msg.ValidatorAddress) {
		return types.ErrWithdrawAddressInfoMismatch(types.DefaultCodespace, msg.From, msg.ValidatorAddress).Result()
	}

	_, ok := k.StakingKeeper.GetValidator(ctx, msg.ValidatorAddress)
	if !ok {
		return types.ErrNoValidatorExist(types.DefaultCodespace, msg.ValidatorAddress.String()).Result()
	}

	_, err := k.WithdrawValidatorCommission(ctx, msg.ValidatorAddress)
	if err != nil {
		return types.ErrHandleTxFailed(types.DefaultCodespace, err).Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.ValidatorAddress.String()),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

//send funds from personal account to communityPool.
func handleMsgFundCommunityPool(ctx sdk.Context, msg types.FundCommunityPoolTx, k keeper.DistrKeeper) sdk.Result {
	if err := k.FundCommunityPool(ctx, msg.Amount, msg.Depositor); err != nil {
		return types.ErrHandleTxFailed(types.DefaultCodespace, err).Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Depositor.String()),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}