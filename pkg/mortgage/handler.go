package mortgage

import (
	"encoding/hex"
	"fmt"
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/mortgage/types"
)

func NewHandler(k MortgageKeeper) sdk.Handler {
	return func(ctx sdk.Context, tx sdk.Tx) sdk.Result {
		switch tx := tx.(type) {
		case types.MsgMortgage:
			return handleMsgMortgage(ctx, k, tx)
		case types.MsgMortgageDone:
			return handleMsgMortgageSuccess(ctx, k, tx)
		case types.MsgMortgageCancel:
			return handleMsgMortgageCancel(ctx, k, tx)
		default:
			errMsg := fmt.Sprintf("unrecognized supply message type: %T", tx)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// 抵押消息
func handleMsgMortgage(ctx sdk.Context, k MortgageKeeper, tx types.MsgMortgage) sdk.Result {

	mort := getMortgage(ctx, k.StoreKey, tx.UniqueID)
	if mort != nil {
		return sdk.ErrInternal("uniqueID is exist").Result()
	}

	if err := k.SupplyKeeper.SendCoinsFromAccountToModule(ctx, tx.FromAddress, types.ModuleName, tx.Coin); err != nil {
		return err.Result()
	}

	setMortgage(ctx, k.StoreKey, types.Mortgage{
		MsgMortgage: tx,
		State:  types.StateMortgaged,
	})
	return sdk.Result{}
}

// 更新抵押取消交易
func handleMsgMortgageCancel (ctx sdk.Context, k MortgageKeeper, tx types.MsgMortgageCancel) sdk.Result {

	mort := getMortgage(ctx, k.StoreKey, tx.UniqueID)
	if mort == nil {
		return sdk.ErrInternal(fmt.Sprintf("mortgage record not exist :uniqueID = %s", hex.EncodeToString(tx.UniqueID))).Result()
	}

	if err := k.SupplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, mort.FromAddress, mort.Coin); err != nil {
		return err.Result()
	}
	if mort.State == types.StateMortgaged {
		mort.State = types.StateCancel
	}
	setMortgage(ctx, k.StoreKey, *mort)
	return sdk.Result{}
}

// 更新抵押状态为成功
func handleMsgMortgageSuccess (ctx sdk.Context, k MortgageKeeper, tx types.MsgMortgageDone) sdk.Result {

	mort := getMortgage(ctx, k.StoreKey, tx.UniqueID)
	if mort == nil {
		return sdk.ErrInternal(fmt.Sprintf("mortgage record not exist :uniqueID = %s", hex.EncodeToString(tx.UniqueID))).Result()
	}

	if err := k.SupplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, mort.ToAddress, mort.Coin); err != nil {
		return err.Result()
	}
	if mort.State == types.StateMortgaged {
		mort.State = types.StateSuccess
	}
	setMortgage(ctx, k.StoreKey, *mort)
	return sdk.Result{}
}

func getMortgage(ctx sdk.Context, key sdk.StoreKey, uniqueID []byte) (mort *types.Mortgage) {
	store := ctx.KVStore(key)
	mortbz := store.Get(uniqueID)
	if len(mortbz) < 1 {
		return nil
	}
	types.MortgageCdc.MustUnmarshalJSON(mortbz, mort)
	return
}

func setMortgage(ctx sdk.Context, key sdk.StoreKey, tx types.Mortgage)  {
	jsonbz := types.MortgageCdc.MustMarshalJSON(tx)
	store := ctx.KVStore(key)
	store.Set(tx.UniqueID, jsonbz)
}