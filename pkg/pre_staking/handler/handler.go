package handler

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/pre_staking/keeper"
	"github.com/ci123chain/ci123chain/pkg/pre_staking/types"
	staking "github.com/ci123chain/ci123chain/pkg/staking/types"
	"github.com/ci123chain/ci123chain/pkg/util"
	gogotypes "github.com/gogo/protobuf/types"
	"math/big"
	"time"
)

const (
	baseMonth = 720
)

func NewHandler(k keeper.PreStakingKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (result *sdk.Result, e error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgStakingDirect:
			return StakingDirectHandler(ctx, k, *msg)
		case *types.MsgRedelegate:
			return RedelegateHandler(ctx, k, *msg)
		case *types.MsgPrestakingCreateValidatorDirect:
			return CreateValidatorDirectHandler(ctx, k, *msg)
		case *types.MsgSetStakingToken:
			return SetStakingTokenHandler(ctx, k, *msg)
		default:
			return nil, nil
		}
	}
}

func SetStakingTokenHandler(ctx sdk.Context, k keeper.PreStakingKeeper, msg types.MsgSetStakingToken) (*sdk.Result, error) {
	if msg.FromAddress.String() == k.GetTokenManagerOwner(ctx) || k.GetTokenManager(ctx) == "" {
		k.SetTokenManager(ctx, msg.TokenAddress)
		k.SetTokenManagerOwner(ctx, msg.FromAddress)
		return &sdk.Result{}, nil
	}
	return nil, types.ErrSetStakingTokenFailed
}


func StakingDirectHandler(ctx sdk.Context, k keeper.PreStakingKeeper, msg types.MsgStakingDirect) (*sdk.Result, error) {
	//update account prestaking.
	denom := k.StakingKeeper.BondDenom(ctx)
	if msg.Amount.Denom != denom {
		return nil, types.ErrInvalidDenom
	}

	k.SetStakingVault(ctx, msg.Validator, msg.FromAddress, ctx.BlockTime().Add(msg.DelegateTime), msg.DelegateTime, msg.Amount,)

	validator, found := k.StakingKeeper.GetValidator(ctx, msg.Validator)
	if !found {
		return nil, types.ErrNoExpectedValidator
	}
	// sub delegator's amount directly
	k.StakingKeeper.Delegate(ctx, msg.Delegator, msg.Amount.Amount, sdk.Unbonded, validator, true)
	moduleAcc := k.SupplyKeeper.GetModuleAccount(ctx, types.ModuleName)
	err := k.AccountKeeper.Transfer(ctx, msg.FromAddress, moduleAcc.GetAddress(), sdk.NewCoins(msg.Amount))
	if err != nil {
		return nil, err
	}

	z := msg.Amount.Amount.BigInt()
	base := int64(msg.DelegateTime.Hours()) / baseMonth
	mintTokens := z.Mul(z, big.NewInt(base))
	if tokenmanager := k.GetTokenManager(ctx); len(tokenmanager) > 0 {
		err = k.SupplyKeeper.Mint(ctx, sdk.HexToAddress(tokenmanager), sdk.HexToAddress(msg.FromAddress.String()), types.ModuleName, mintTokens)
		if err != nil {
			return nil, err
		}
	} else {
		ctx.Logger().Warn("StakingToken Address not set")
	}
	em := sdk.NewEventManager()
	em.EmitEvents(sdk.Events{
		sdk.NewEvent(types.EventMsgStaking,
			sdk.NewAttribute([]byte(sdk.AttributeKeyMethod), []byte(types.EventMsgStaking)),
			sdk.NewAttribute([]byte(sdk.AttributeKeyAmount), []byte(msg.Amount.String())),
			sdk.NewAttribute([]byte(sdk.AttributeKeyModule), []byte(types.ModuleName)),
			sdk.NewAttribute([]byte(sdk.AttributeKeySender), []byte(msg.FromAddress.String())),
			sdk.NewAttribute([]byte(sdk.AttributeKeyAmount), []byte(msg.Amount.Amount.String())),
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

	delegation, found := k.StakingKeeper.GetDelegation(ctx, msg.FromAddress, srcValidator.OperatorAddress)
	if !found {
		return nil, types.ErrNoExpectedDelegation
	}


	if srcValidator.InvalidExRate() {
		return nil, types.ErrDelegatorShareExRateInvalid
	}

	completionTime, err := k.StakingKeeper.Redelegate(ctx, msg.FromAddress, msg.SrcValidator, msg.DstValidator, delegation.Shares)
	if err != nil {
		return nil, types.ErrRedelegateFailed
	}

	//Update staking record.
	if err := k.ChangeStakingRecordToNewValidator(ctx, msg.RecordID, msg.SrcValidator, msg.DstValidator); err != nil {
		return nil, err
	}

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
			sdk.NewAttribute([]byte(types.AttributeKeyCompletionTime), []byte(completionTime.Format(time.RFC3339))),
			sdk.NewAttribute([]byte(sdk.AttributeKeyModule), []byte(types.AttributeValueCategory)),
			sdk.NewAttribute([]byte(sdk.AttributeKeySender), []byte(msg.FromAddress.String())),
			sdk.NewAttribute([]byte(types.AttributeKeyVaultID), []byte(msg.FromAddress.String())),
		),
	})

	return &sdk.Result{Data: completionTimeBz, Events: em.Events()}, nil
}

func CreateValidatorDirectHandler(ctx sdk.Context, k keeper.PreStakingKeeper, msg types.MsgPrestakingCreateValidatorDirect) (*sdk.Result, error) {
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

	k.SetStakingVault(ctx, msg.ValidatorAddress, msg.FromAddress, ctx.BlockTime().Add(msg.DelegateTime), msg.DelegateTime, msg.Amount)

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

	_, err = k.StakingKeeper.Delegate(ctx, msg.DelegatorAddress, msg.Amount.Amount, sdk.Unbonded, validator, true)
	if err != nil {
		return nil, err
	}
	z := msg.Amount.Amount.BigInt()
	base := int64(msg.DelegateTime.Hours()) / baseMonth
	mintTokens := z.Mul(z, big.NewInt(base))
	if tokenmanager := k.GetTokenManager(ctx); len(tokenmanager) > 0 {
		err = k.SupplyKeeper.Mint(ctx, sdk.HexToAddress(tokenmanager), sdk.HexToAddress(msg.FromAddress.String()), types.ModuleName, mintTokens)
		if err != nil {
			return nil, err
		}
	} else {
		ctx.Logger().Warn("StakingToken Address not set")
	}
	em := sdk.NewEventManager()
	em.EmitEvents(sdk.Events{
		sdk.NewEvent(types.EventTypeCreateValidator,
			sdk.NewAttribute([]byte(sdk.AttributeKeyMethod), []byte(types.EventTypeCreateValidator)),
			sdk.NewAttribute([]byte(sdk.AttributeKeyAmount), []byte(msg.Amount.Amount.String())),
			sdk.NewAttribute([]byte(sdk.AttributeKeyModule), []byte(types.AttributeValueCategory)),
			sdk.NewAttribute([]byte(sdk.AttributeKeySender), []byte(msg.FromAddress.String())),
			sdk.NewAttribute([]byte(sdk.AttributeKeyAmount), []byte(msg.Amount.Amount.String())),
		),
	})

	return &sdk.Result{Events: em.Events()}, nil
}
