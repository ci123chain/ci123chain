package mortgage

import (
	"encoding/hex"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/mortgage/types"
	"github.com/pkg/errors"
)

func NewHandler(k MortgageKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case *types.MsgMortgage:
			return handleMsgMortgage(ctx, k, *msg)
		case *types.MsgMortgageDone:
			return handleMsgMortgageSuccess(ctx, k, *msg)
		case *types.MsgMortgageCancel:
			return handleMsgMortgageCancel(ctx, k, *msg)
		default:
			errMsg := fmt.Sprintf("unrecognized supply message type: %T", msg)
			return nil, errors.New(errMsg)
		}
	}
}

// 抵押消息
func handleMsgMortgage(ctx sdk.Context, k MortgageKeeper, msg types.MsgMortgage) (*sdk.Result, error) {

	mort := getMortgage(ctx, k.StoreKey, msg.UniqueID)
	if mort != nil {
		return nil, errors.New("uniqueID is exist")
	}

	if err := k.SupplyKeeper.SendCoinsFromAccountToModule(ctx, msg.FromAddress, types.ModuleName, msg.Coin); err != nil {
		return nil, err
	}
	setMortgage(ctx, k.StoreKey, types.Mortgage{
		MsgMortgage: msg,
		State:  types.StateMortgaged,
	})
	return &sdk.Result{}, nil
}

// 更新抵押取消交易
func handleMsgMortgageCancel (ctx sdk.Context, k MortgageKeeper, msg types.MsgMortgageCancel) (*sdk.Result, error) {

	mort := getMortgage(ctx, k.StoreKey, msg.UniqueID)
	if mort == nil {
		return nil, errors.New(fmt.Sprintf("mortgage record not exist :uniqueID = %s", hex.EncodeToString(msg.UniqueID)))
	}
	if !mort.FromAddress.Equal(msg.FromAddress) {
		return nil, errors.New(fmt.Sprintf("account address mismatch, expected %s, got %s", msg.FromAddress.String(), mort.FromAddress.String()))
	}

	if err := k.SupplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, mort.FromAddress, mort.Coin); err != nil {
		return nil, err
	}
	if mort.State == types.StateMortgaged {
		mort.State = types.StateCancel
	} else {
		return nil, errors.New("mortgage record state have done or canceled")
	}
	setMortgage(ctx, k.StoreKey, *mort)
	return &sdk.Result{}, nil
}

// 更新抵押状态为成功
func handleMsgMortgageSuccess (ctx sdk.Context, k MortgageKeeper, msg types.MsgMortgageDone) (*sdk.Result, error) {

	mort := getMortgage(ctx, k.StoreKey, msg.UniqueID)
	if mort == nil {
		return nil, errors.New(fmt.Sprintf("mortgage record not exist :uniqueID = %s", hex.EncodeToString(msg.UniqueID)))
	}

	if !mort.FromAddress.Equal(msg.FromAddress) {
		return nil, errors.New(fmt.Sprintf("account address mismatch, expected %s, got %s", msg.FromAddress.String(), mort.FromAddress.String()))
	}

	if err := k.SupplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, mort.ToAddress, mort.Coin); err != nil {
		return nil, err
	}
	if mort.State == types.StateMortgaged {
		mort.State = types.StateSuccess
	} else {
		return nil, errors.New("mortgage record state have done or canceled")
	}
	setMortgage(ctx, k.StoreKey, *mort)
	return &sdk.Result{}, nil
}

func getMortgage(ctx sdk.Context, key sdk.StoreKey, uniqueID []byte) (*types.Mortgage) {
	store := ctx.KVStore(key)
	mortbz := store.Get(uniqueID)
	if len(mortbz) < 1 {
		return nil
	}
	var mort types.Mortgage
	err := types.MortgageCdc.UnmarshalBinaryLengthPrefixed(mortbz, &mort)
	if err != nil {
		panic(err)
	}
	return &mort
}

func setMortgage(ctx sdk.Context, key sdk.StoreKey, tx types.Mortgage)  {
	jsonbz, err := types.MortgageCdc.MarshalBinaryLengthPrefixed(tx)
	store := ctx.KVStore(key)
	store.Set(tx.UniqueID, jsonbz)

	var mort types.Mortgage
	err = types.MortgageCdc.UnmarshalBinaryLengthPrefixed(jsonbz, &mort)
	if err != nil {
		panic(err)
	}
}
