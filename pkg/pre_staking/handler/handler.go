package handler

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/pre_staking/keeper"
	"github.com/ci123chain/ci123chain/pkg/pre_staking/types"
	types2 "github.com/ci123chain/ci123chain/pkg/staking/types"
	gogotypes "github.com/gogo/protobuf/types"
	"time"
)


func NewHandler(k keeper.PreStakingKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (result *sdk.Result, e error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgPreStaking:
			return PreStakingHandler(ctx, k, *msg)
		case *types.MsgStaking:
			return StakingHandler(ctx, k, *msg)
		case *types.MsgUndelegate:
			return UndelegateHandler(ctx, k, *msg)
		case *types.MsgRedelegate:
			return RedelegateHandler(ctx, k, *msg)
		default:
			return nil, nil
		}
	}
}


func PreStakingHandler(ctx sdk.Context, k keeper.PreStakingKeeper, msg types.MsgPreStaking) (*sdk.Result, error) {

	balance := k.AccountKeeper.GetBalance(ctx, msg.FromAddress).AmountOf(sdk.ChainCoinDenom)
	if balance.LT(msg.Amount.Amount) {
		return nil, types.ErrAccountBalanceNotEnough
	}
	//pay to module account.
	moduleAcc := k.SupplyKeeper.GetModuleAccount(ctx, types.DefaultCodespace)
	err := k.AccountKeeper.Transfer(ctx, msg.FromAddress, moduleAcc.GetAddress(), sdk.NewCoins(msg.Amount))
	if err != nil {
		return nil, err
	}

	//call contract.
	//TODO

	//save to keeper.
	res := k.GetAccountPreStaking(ctx, msg.FromAddress)
	k.SetAccountPreStaking(ctx, msg.FromAddress, msg.Amount.Amount.Add(res))

	//events.
	em := sdk.NewEventManager()
	em.EmitEvents(sdk.Events{
		sdk.NewEvent(types.EventsMsgPreStaking,
			sdk.NewAttribute([]byte(sdk.AttributeKeyMethod), []byte(types.EventsMsgPreStaking)),
			sdk.NewAttribute([]byte(sdk.AttributeKeyAmount), []byte(msg.Amount.Amount.String())),
			sdk.NewAttribute([]byte(sdk.AttributeKeyModule), []byte(types.ModuleName)),
			sdk.NewAttribute([]byte(sdk.AttributeKeySender), []byte(msg.FromAddress.String())),
		),
	})
	return &sdk.Result{Events:em.Events()}, nil
}

func StakingHandler(ctx sdk.Context, k keeper.PreStakingKeeper, msg types.MsgStaking) (*sdk.Result, error) {
	//check
	res := k.GetAccountPreStaking(ctx, msg.Delegator)
	if res.IsZero() || res.LT(msg.Amount.Amount) {
		return nil, types.ErrAccountBalanceNotEnough
	}
	//update account prestaking.
	k.SetAccountPreStaking(ctx, msg.Delegator, res.Sub(msg.Amount.Amount))

	//var update = ctx.BlockTime()
	//t, err := time.ParseDuration(msg.StorageTime.String())
	//var end = update.Add(t)
	//var record = types.NewStakingRecord(msg.StorageTime, update, end, msg.Amount)
	//var key = types.GetStakingRecordKey(msg.FromAddress, msg.Validator)
	Err := k.SetAccountStakingRecord(ctx, msg.Validator, msg.FromAddress, msg.StorageTime, msg.Amount)
	if Err != nil {
		return nil, Err
	}

	validator, found := k.StakingKeeper.GetValidator(ctx, msg.Validator)
	if !found {
		return nil, types.ErrNoExpectedValidator
	}
	denom := k.StakingKeeper.BondDenom(ctx)
	if msg.Amount.Denom != denom {
		return nil, types.ErrInvalidDenom
	}

	if validator.InvalidExRate() {
		return nil, types.ErrDelegatorShareExRateInvalid
	}

	delegation, found := k.StakingKeeper.GetDelegation(ctx, msg.Delegator, validator.OperatorAddress)
	if !found {
		delegation = types2.NewDelegation(msg.Delegator, validator.OperatorAddress, sdk.ZeroDec())
	}

	if found {
		k.StakingKeeper.BeforeDelegationSharesModified(ctx, msg.Delegator, validator.OperatorAddress)
	}else {
		k.StakingKeeper.BeforeDelegationCreated(ctx, msg.Delegator, validator.OperatorAddress)
	}

	var sendName string
	switch {
	case validator.IsBonded():
		sendName = types2.BondedPoolName
	case validator.IsUnbonding(), validator.IsUnbonded():
		sendName = types2.NotBondedPoolName
	default:
		panic("invalid validator status")
	}

	//pay to validator account.
	err := k.SupplyKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, sendName, msg.Amount)
	if err != nil {
		return nil, err
	}

	validator, newShares := k.StakingKeeper.AddValidatorTokensAndShares(ctx, validator, msg.Amount.Amount)
	//update delegation
	delegation.Shares = delegation.Shares.Add(newShares)
	k.StakingKeeper.SetDelegation(ctx, delegation)
	k.StakingKeeper.AfterDelegationModified(ctx, delegation.DelegatorAddress, delegation.ValidatorAddress)

	//events.
	em := sdk.NewEventManager()
	em.EmitEvents(sdk.Events{
		sdk.NewEvent(types.EventMsgStaking,
			sdk.NewAttribute([]byte(sdk.AttributeKeyMethod), []byte(types.EventMsgStaking)),
			sdk.NewAttribute([]byte(sdk.AttributeKeyAmount), []byte(msg.Amount.Amount.String())),
			sdk.NewAttribute([]byte(sdk.AttributeKeyModule), []byte(types.ModuleName)),
			sdk.NewAttribute([]byte(sdk.AttributeKeySender), []byte(msg.FromAddress.String())),
		),
	})
	return &sdk.Result{Events:em.Events()}, nil
}


func UndelegateHandler(ctx sdk.Context, k keeper.PreStakingKeeper, msg types.MsgUndelegate) (*sdk.Result, error) {

	res := k.GetAccountPreStaking(ctx, msg.FromAddress)
	if res.IsZero() {
		return nil, types.ErrNoBalanceLeft
	}
	if res.LT(msg.Amount.Amount) {
		return nil, types.ErrNoEnoughBalanceLeft
	}

	k.SetAccountPreStaking(ctx, msg.FromAddress, res.Sub(msg.Amount.Amount))

	moduleAcc := k.SupplyKeeper.GetModuleAccount(ctx, types.DefaultCodespace)
	err := k.AccountKeeper.Transfer(ctx, moduleAcc.GetAddress(), msg.FromAddress, sdk.NewCoins(msg.Amount))
	if err != nil {
		return nil, err
	}
	//TODO
	//call contract.

	//events.
	em := sdk.NewEventManager()
	em.EmitEvents(sdk.Events{
		sdk.NewEvent(types.EventsMsgPreStaking,
			sdk.NewAttribute([]byte(sdk.AttributeKeyMethod), []byte(types.EventUndelegate)),
			sdk.NewAttribute([]byte(sdk.AttributeKeyAmount), []byte(msg.Amount.Amount.String())),
			sdk.NewAttribute([]byte(sdk.AttributeKeyModule), []byte(types.ModuleName)),
			sdk.NewAttribute([]byte(sdk.AttributeKeySender), []byte(msg.FromAddress.String())),
		),
	})
	return &sdk.Result{Events:em.Events()}, nil
}


func RedelegateHandler(ctx sdk.Context, k keeper.PreStakingKeeper, msg types.MsgRedelegate) (*sdk.Result, error) {

	srcValidator, found := k.StakingKeeper.GetValidator(ctx, msg.SrcValidator)
	if !found {
		return nil, types.ErrNoExpectedValidator
	}
	_, found = k.StakingKeeper.GetValidator(ctx, msg.DstValidator)
	if !found {
		return nil, types.ErrNoExpectedValidator
	}
	//denom := k.StakingKeeper.BondDenom(ctx)
	//if msg.Amount.Denom != denom {
	//	return nil, types.ErrInvalidDenom
	//}

	delegation, found := k.StakingKeeper.GetDelegation(ctx, msg.FromAddress, srcValidator.OperatorAddress)
	if !found {
		return nil, types.ErrNoExpectedDelegation
	}
	//if delegation.Shares.LT(msg.Amount.Amount.ToDec()) {
	//	return nil, types.ErrNoEnoughSharesToRedelegate
	//}

	if srcValidator.InvalidExRate() {
		return nil, types.ErrDelegatorShareExRateInvalid
	}
	//shares, err := k.StakingKeeper.ValidateUnbondAmount(ctx, msg.FromAddress, msg.SrcValidator, msg.Amount.Amount)
	//if err != nil {
	//	return nil, err
	//}

	completionTime, err := k.StakingKeeper.Redelegate(ctx, msg.FromAddress, msg.SrcValidator, msg.DstValidator, delegation.Shares)
	if err != nil {
		return nil, types.ErrRedelegateFailed
	}

	//update staking record.
	srcRecord := k.GetAccountStakingRecord(ctx, msg.SrcValidator, msg.FromAddress)
	k.ClearStakingRecord(ctx, msg.SrcValidator, msg.DstValidator)
	//dstRecord := k.GetAccountStakingRecord(ctx, msg.DstValidator, msg.FromAddress)
	//k.SetAccountStakingRecord(ctx, msg.DstValidator, msg.FromAddress, srcRecord)
	k.UpdateStakingRecord(ctx, msg.SrcValidator, msg.FromAddress, srcRecord)

	ts, err := gogotypes.TimestampProto(completionTime)
	if err != nil {
		return nil, types.ErrTimestampProto
	}

	completionTimeBz := types.PreStakingCodec.MustMarshalBinaryLengthPrefixed(ts)
	em := sdk.NewEventManager()
	em.EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRedelegate,
			sdk.NewAttribute([]byte(sdk.AttributeKeyMethod), []byte(types.EventTypeRedelegate)),
			sdk.NewAttribute([]byte(types.AttributeKeySrcValidator), []byte(msg.SrcValidator.String())),
			sdk.NewAttribute([]byte(types.AttributeKeyDstValidator), []byte(msg.DstValidator.String())),
			//sdk.NewAttribute([]byte(sdk.AttributeKeyAmount), []byte(msg.Amount.Amount.String())),
			sdk.NewAttribute([]byte(types.AttributeKeyCompletionTime), []byte(completionTime.Format(time.RFC3339))),
			sdk.NewAttribute([]byte(sdk.AttributeKeyModule), []byte(types.AttributeValueCategory)),
			sdk.NewAttribute([]byte(sdk.AttributeKeySender), []byte(msg.FromAddress.String())),
		),
	})

	return &sdk.Result{Data: completionTimeBz, Events: em.Events()}, nil
}