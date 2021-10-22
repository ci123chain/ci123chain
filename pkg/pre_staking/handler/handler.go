package handler

import (
	"errors"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/pre_staking/keeper"
	"github.com/ci123chain/ci123chain/pkg/pre_staking/types"
	staking "github.com/ci123chain/ci123chain/pkg/staking/types"
	types2 "github.com/ci123chain/ci123chain/pkg/staking/types"
	"github.com/ci123chain/ci123chain/pkg/util"
	gogotypes "github.com/gogo/protobuf/types"
	"math/big"
	"time"
)

const (
	tokenManager = "0x5B1427075C0EF657a6F4f23A7EF7065E028cAd3b"
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
		case *types.MsgPrestakingCreateValidator:
			return CreateValidatorHandler(ctx, k, *msg)
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
	moduleAcc := k.SupplyKeeper.GetModuleAccount(ctx, types.ModuleName)
	err := k.AccountKeeper.Transfer(ctx, msg.FromAddress, moduleAcc.GetAddress(), sdk.NewCoins(msg.Amount))
	if err != nil {
		return nil, err
	}
	var z = msg.Amount.Amount.BigInt()
	util.AddDecimal(z, 18, 10)
	err = k.SupplyKeeper.Mint(ctx, sdk.HexToAddress(tokenManager), sdk.HexToAddress(msg.FromAddress.String()), types.ModuleName, z)
	if err != nil {
		return nil, err
	}
	//save to keeper.
	res := k.GetAccountPreStaking(ctx, msg.FromAddress)
	vat := types.NewVault(ctx.BlockTime(), ctx.BlockTime().Add(msg.DelegateTime), msg.DelegateTime, msg.Amount)
	res.AddVault(vat)
	k.SetAccountPreStaking(ctx, msg.FromAddress, res)

	//events.
	em := sdk.NewEventManager()
	em.EmitEvents(sdk.Events{
		sdk.NewEvent(types.EventsMsgPreStaking,
			sdk.NewAttribute([]byte(sdk.AttributeKeyMethod), []byte(types.EventsMsgPreStaking)),
			sdk.NewAttribute([]byte(sdk.AttributeKeyAmount), []byte(msg.Amount.Amount.String())),
			sdk.NewAttribute([]byte(sdk.AttributeKeyModule), []byte(types.ModuleName)),
			sdk.NewAttribute([]byte(sdk.AttributeKeySender), []byte(msg.FromAddress.String())),
			sdk.NewAttribute([]byte(types.AttributeKeyVaultID), []byte(res.LatestVaultID.String())),
		),
	})
	return &sdk.Result{Events:em.Events()}, nil
}

func StakingHandler(ctx sdk.Context, k keeper.PreStakingKeeper, msg types.MsgStaking) (*sdk.Result, error) {
	res := k.GetAccountPreStaking(ctx, msg.Delegator)
	if res.IsEmpty() {
		return nil, types.ErrAccountBalanceNotEnough
	}
	id, ok := new(big.Int).SetString(msg.VaultID, 10)
	if !ok {
		return nil, types.ErrInvalidVaultID
	}
	amount, endTime, err := res.PopVaultAmountAndEndTime(id)
	if err != nil {
		return nil, err
	}
	//update account prestaking.
	k.SetAccountPreStaking(ctx, msg.Delegator, res)

	Err := k.SetAccountStakingRecord(ctx, msg.Validator, msg.FromAddress, id, endTime, amount)
	if Err != nil {
		return nil, Err
	}

	validator, found := k.StakingKeeper.GetValidator(ctx, msg.Validator)
	if !found {
		return nil, types.ErrNoExpectedValidator
	}
	denom := k.StakingKeeper.BondDenom(ctx)
	if amount.Denom != denom {
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
	err = k.SupplyKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, sendName, amount)
	if err != nil {
		return nil, err
	}

	validator, newShares := k.StakingKeeper.AddValidatorTokensAndShares(ctx, validator, amount.Amount)
	//update delegation
	delegation.Shares = delegation.Shares.Add(newShares)
	k.StakingKeeper.SetDelegation(ctx, delegation)
	k.StakingKeeper.AfterDelegationModified(ctx, delegation.DelegatorAddress, delegation.ValidatorAddress)

	//events.
	em := sdk.NewEventManager()
	em.EmitEvents(sdk.Events{
		sdk.NewEvent(types.EventMsgStaking,
			sdk.NewAttribute([]byte(sdk.AttributeKeyMethod), []byte(types.EventMsgStaking)),
			sdk.NewAttribute([]byte(sdk.AttributeKeyAmount), []byte(amount.String())),
			sdk.NewAttribute([]byte(sdk.AttributeKeyModule), []byte(types.ModuleName)),
			sdk.NewAttribute([]byte(sdk.AttributeKeySender), []byte(msg.FromAddress.String())),
			sdk.NewAttribute([]byte(types.AttributeKeyVaultID), []byte(msg.VaultID)),
		),
	})
	return &sdk.Result{Events:em.Events()}, nil
}


func UndelegateHandler(ctx sdk.Context, k keeper.PreStakingKeeper, msg types.MsgUndelegate) (*sdk.Result, error) {

	res := k.GetAccountPreStaking(ctx, msg.FromAddress)
	if res.IsEmpty() {
		return nil, types.ErrNoBalanceLeft
	}
	id, ok := new(big.Int).SetString(msg.VaultID, 10)
	if !ok {
		return nil, types.ErrInvalidVaultID
	}
	if id.Uint64() < 1 {
		return nil, types.ErrInvalidVaultID
	}
	amount, et, err := res.PopVaultAmountAndEndTime(id)
	if err != nil {
		return nil, err
	}
	if !et.Before(ctx.BlockTime()) {
		return nil, errors.New("you can only undelegate the vault after endtime")
	}
	if !amount.IsPositive() {
		return nil, types.ErrNoEnoughBalanceLeft
	}

	k.SetAccountPreStaking(ctx, msg.FromAddress, res)

	moduleAcc := k.SupplyKeeper.GetModuleAccount(ctx, types.DefaultCodespace)
	err = k.AccountKeeper.Transfer(ctx, moduleAcc.GetAddress(), msg.FromAddress, sdk.NewCoins(amount))
	if err != nil {
		return nil, err
	}

	var z = amount.Amount.BigInt()
	util.AddDecimal(z, 18, 10)
	err = k.SupplyKeeper.BurnEVMCoin(ctx, types.ModuleName, sdk.HexToAddress(tokenManager), msg.FromAddress, z)
	if err != nil {
		return nil, err
	}

	//events.
	em := sdk.NewEventManager()
	em.EmitEvents(sdk.Events{
		sdk.NewEvent(types.EventsMsgPreStaking,
			sdk.NewAttribute([]byte(sdk.AttributeKeyMethod), []byte(types.EventUndelegate)),
			sdk.NewAttribute([]byte(sdk.AttributeKeyAmount), []byte(amount.Amount.String())),
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
			sdk.NewAttribute([]byte(types.AttributeKeyVaultID), []byte(msg.FromAddress.String())),
		),
	})

	return &sdk.Result{Data: completionTimeBz, Events: em.Events()}, nil
}


func CreateValidatorHandler(ctx sdk.Context, k keeper.PreStakingKeeper, msg types.MsgPrestakingCreateValidator) (*sdk.Result, error) {
	if _, found := k.StakingKeeper.GetValidator(ctx, msg.ValidatorAddress); found {
		return nil, types.ErrNoExpectedValidator
	}
	pk, err := util.ParsePubKey(msg.PublicKey)
	if err != nil {
		return nil, types.ErrInvalidPublicKey
	}

	if _, found := k.StakingKeeper.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(pk)); found {
		return nil, types.ErrPubkeyHasBonded
	}

	if _, err := msg.Description.EnsureLength(); err != nil {
		return nil, err
	}

	//get amount of preDelegate.
	res := k.GetAccountPreStaking(ctx, msg.FromAddress)
	if res.IsEmpty() {
		return nil, types.ErrAccountBalanceNotEnough
	}
	id, ok := new(big.Int).SetString(msg.VaultID, 10)
	if !ok {
		return nil, types.ErrInvalidVaultID
	}
	amount, endTime, err := res.PopVaultAmountAndEndTime(id)
	if err != nil {
		return nil, err
	}

	//update account prestaking.
	k.SetAccountPreStaking(ctx, msg.FromAddress, res)

	Err := k.SetAccountStakingRecord(ctx, msg.ValidatorAddress, msg.FromAddress, id, endTime, amount)
	if Err != nil {
		return nil, Err
	}

	validator, _ := staking.NewValidator(msg.ValidatorAddress, msg.PublicKey, msg.Description)

	commission := staking.NewCommissionWithTime(msg.Commission.Rate,
		msg.Commission.MaxRate, msg.Commission.MaxChangeRate, ctx.BlockHeader().Time)

	validator, err = validator.SetInitialCommission(commission)
	if err != nil {
		return nil, err
	}
	validator.MinSelfDelegation = msg.MinSelfDelegation

	err = k.StakingKeeper.SetValidator(ctx, validator)
	if err != nil {
		return nil, types.ErrSetValidatorFailed
	}
	k.StakingKeeper.SetValidatorByConsAddr(ctx, validator)
	k.StakingKeeper.SetNewValidatorByPowerIndex(ctx, validator)

	k.StakingKeeper.AfterValidatorCreated(ctx, validator.OperatorAddress)

	_, err = k.StakingKeeper.Delegate(ctx, msg.DelegatorAddress, amount.Amount, sdk.Unbonded, validator, true)
	if err != nil {
		return nil, err
	}

	em := sdk.NewEventManager()
	em.EmitEvents(sdk.Events{
		sdk.NewEvent(types.EventTypeCreateValidator,
			sdk.NewAttribute([]byte(sdk.AttributeKeyMethod), []byte(types.EventTypeCreateValidator)),
			sdk.NewAttribute([]byte(sdk.AttributeKeyAmount), []byte(amount.Amount.String())),
			sdk.NewAttribute([]byte(sdk.AttributeKeyModule), []byte(types.AttributeValueCategory)),
			sdk.NewAttribute([]byte(sdk.AttributeKeySender), []byte(msg.FromAddress.String())),
			sdk.NewAttribute([]byte(sdk.AttributeKeyVaultID), []byte(msg.VaultID)),
		),
	})

	return &sdk.Result{Events: em.Events()}, nil
}
