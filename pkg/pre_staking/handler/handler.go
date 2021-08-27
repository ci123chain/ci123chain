package handler

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/pre_staking/keeper"
	"github.com/ci123chain/ci123chain/pkg/pre_staking/types"
	types2 "github.com/ci123chain/ci123chain/pkg/staking/types"
)


func NewHandler(k keeper.PreStakingKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (result *sdk.Result, e error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgPreStaking:
			return PreStakingHandler(ctx, k, *msg)
		case *types.MsgStaking:
			return StakingHandler(ctx, k, *msg)
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
	k.SetAccountPreStaking(ctx, msg.Delegator, res.Sub(msg.Amount.Amount))

	validator, found := k.StakingKeeper.GetValidator(ctx, msg.Validator)
	if !found {
		//r := fmt.Sprintf("validator %s has existed", msg.ValidatorAddress.String())
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