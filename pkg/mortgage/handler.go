package mortgage

import (
	"encoding/hex"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/mortgage/types"
	"github.com/ci123chain/ci123chain/pkg/transaction"
)

func NewHandler(k MortgageKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case *types.MsgMortgage:
			return handleMsgMortgage(ctx, k, *msg)
		case *types.MsgMortgageDone:
			return handleMsgMortgageSuccess(ctx, k, *msg)
		case *types.MsgMortgageCancel:
			return handleMsgMortgageCancel(ctx, k, *msg)
		default:
			errMsg := fmt.Sprintf("unrecognized supply message types: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// 抵押消息
func handleMsgMortgage(ctx sdk.Context, k MortgageKeeper, msg types.MsgMortgage) sdk.Result {

	mort := getMortgage(ctx, k.StoreKey, msg.UniqueID)
	if mort != nil {
		return transaction.ErrInvalidTx(types.DefaultCodespace,"uniqueID is exist").Result()
	}

	if err := k.SupplyKeeper.SendCoinsFromAccountToModule(ctx, msg.FromAddress, types.ModuleName, msg.Coin); err != nil {
		return err.Result()
	}
	setMortgage(ctx, k.StoreKey, types.Mortgage{
		MsgMortgage: msg,
		State:  types.StateMortgaged,
	})
	return sdk.Result{}
}

// 更新抵押取消交易
func handleMsgMortgageCancel (ctx sdk.Context, k MortgageKeeper, msg types.MsgMortgageCancel) sdk.Result {

	mort := getMortgage(ctx, k.StoreKey, msg.UniqueID)
	if mort == nil {
		return transaction.ErrInvalidTx(types.DefaultCodespace, fmt.Sprintf("mortgage record not exist :uniqueID = %s", hex.EncodeToString(msg.UniqueID))).Result()
	}
	if !mort.FromAddress.Equal(msg.FromAddress) {
		return transaction.ErrInvalidTx(types.DefaultCodespace, "from address error").Result()
	}

	if err := k.SupplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, mort.FromAddress, mort.Coin); err != nil {
		return err.Result()
	}
	if mort.State == types.StateMortgaged {
		mort.State = types.StateCancel
	} else {
		return transaction.ErrInvalidTx(types.DefaultCodespace, "mortgage record state have done or canceled").Result()
	}
	setMortgage(ctx, k.StoreKey, *mort)
	return sdk.Result{}
}

// 更新抵押状态为成功
func handleMsgMortgageSuccess (ctx sdk.Context, k MortgageKeeper, msg types.MsgMortgageDone) sdk.Result {

	mort := getMortgage(ctx, k.StoreKey, msg.UniqueID)
	if mort == nil {
		return transaction.ErrInvalidTx(types.DefaultCodespace, fmt.Sprintf("mortgage record not exist :uniqueID = %s", hex.EncodeToString(msg.UniqueID))).Result()
	}

	if !mort.FromAddress.Equal(msg.FromAddress) {
		return transaction.ErrInvalidTx(types.DefaultCodespace, "from address error").Result()
	}

	if err := k.SupplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, mort.ToAddress, mort.Coin); err != nil {
		return transaction.ErrSendCoin(types.DefaultCodespace, err).Result()
	}
	if mort.State == types.StateMortgaged {
		mort.State = types.StateSuccess
	} else {
		return transaction.ErrInvalidTx(types.DefaultCodespace, "mortgage record state have done or canceled").Result()
	}
	setMortgage(ctx, k.StoreKey, *mort)
	return sdk.Result{}
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
