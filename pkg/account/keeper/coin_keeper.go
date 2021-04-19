package keeper

import (
	"fmt"
	prefix "github.com/ci123chain/ci123chain/pkg/abci/store"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/ci123chain/ci123chain/pkg/account/types"
)

func (ak AccountKeeper) GetBalance(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	coin := ak.getBalance(ctx, addr)
	if coin.IsZero() {
		return sdk.NewCoins()
	}
	return coin
}

func (ak AccountKeeper) AddBalance(ctx sdk.Context, addr sdk.AccAddress, amount sdk.Coins) (sdk.Coins, error) {
	if !amount.IsValid() {
		return sdk.NewCoins(), sdkerrors.ErrInvalidCoins
	}
	oldCoins := ak.GetBalance(ctx, addr)
	//newCoin := oldCoin.Add(amount)
	//
	//if newCoin.IsNegative() {
	//	return amount, sdk.ErrInsufficientCoins(
	//		fmt.Sprintf("insufficient account funds: %s < %s", oldCoin, amount),
	//	)
	//}
	newCoins := oldCoins.Add(amount)
	err := ak.SetCoin(ctx, addr, newCoins)
	return newCoins, err
}

func (ak AccountKeeper) SetCoin(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) error {
	if !amt.IsValid() {
		return sdkerrors.ErrInvalidCoins
	}

	acc := ak.GetAccount(ctx, addr)
	if acc == nil {
		acc = ak.NewAccountWithAddress(ctx, addr)
	}

	err := acc.SetCoins(amt)
	if err != nil {
		panic(err)
	}

	ak.SetAccount(ctx, acc)
	return nil
}

func (ak AccountKeeper) SubBalance(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, error) {
	if !amt.IsValid() {
		return sdk.NewCoins(), sdkerrors.ErrInvalidCoins
	}

	oldCoins, spendableCoins := sdk.NewCoins(), sdk.NewCoins()//sdk.NewChainCoin(sdk.NewInt(0)), sdk.NewChainCoin(sdk.NewInt(0))

	acc := ak.GetAccount(ctx, addr)
	if acc != nil {
		oldCoins = acc.GetCoins()
		spendableCoins = acc.SpendableCoins(ctx.BlockHeader().Time)
	} else {
		return sdk.NewCoins(), sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("account not exist %s", addr.Hex()))
	}
	_, hasNeg := spendableCoins.SafeSub(amt)
	if hasNeg {
		return amt, sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, fmt.Sprintf("insufficient accounts funds; %s < %s", spendableCoins, amt))
	}

	newCoin := oldCoins.Sub(amt)
	err := ak.SetCoin(ctx, addr, newCoin)
	return newCoin, err
}


func (ak AccountKeeper) Transfer(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, coin sdk.Coins) error {
	if ctx.EventManager() != nil {
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeTransfer,
				sdk.NewAttributeString(types.AttributeKeyRecipient, toAddr.String()),
				sdk.NewAttributeString(types.AttributeKeySender, fromAddr.String()),
				sdk.NewAttributeString(sdk.AttributeKeyAmount, coin.String()),
			),
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttributeString(types.AttributeKeySender, fromAddr.String()),
			),
		})
	}

	_, err := ak.SubBalance(ctx, fromAddr, coin)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInternal, err.Error())
	}

	_, err = ak.AddBalance(ctx, toAddr, coin)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInternal, err.Error())
	}

	return nil
}

func (ak AccountKeeper) getStore(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(ak.key)
}


func (k AccountKeeper) getBalance(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	acc := k.GetAccount(ctx, addr)
	if acc == nil {
		return sdk.NewCoins(sdk.NewChainCoin(sdk.NewInt(0)))
	}
	return acc.GetCoins()
}

func (k AccountKeeper) SetBalances(ctx sdk.Context, addr sdk.AccAddress, balances sdk.Coins) error {

	k.ClearBalances(ctx, addr)
	err := k.SetCoin(ctx, addr, balances)
	if err != nil {
		return err
	}

	//for _, balance := range balances {
	//	err := k.SetCoin(ctx, addr, balance)
	//	if err != nil {
	//		return err
	//	}
	//}

	return nil
}

func (k AccountKeeper) ClearBalances(ctx sdk.Context, addr sdk.AccAddress) {

	keys := [][]byte{}
	k.IterateAccountBalances(ctx, addr, func(balance sdk.Coin) bool {
		keys = append(keys, []byte(balance.Denom))
		return false
	})

	store := ctx.KVStore(k.key)
	balancesStore := prefix.NewPrefixStore(store, types.BalancesPrefix)
	accountStore := prefix.NewPrefixStore(balancesStore, addr.Bytes())

	for _, key := range keys {
		accountStore.Delete(key)
	}

}

// IterateAccountBalances iterates over the balances of a single account and
// provides the token balance to a callback. If true is returned from the
// callback, iteration is halted.
func (k AccountKeeper) IterateAccountBalances(ctx sdk.Context, addr sdk.AccAddress, cb func(sdk.Coin) bool) {

	store := ctx.KVStore(k.key)
	balancesStore := prefix.NewPrefixStore(store, types.BalancesPrefix)
	accountStore := prefix.NewPrefixStore(balancesStore, addr.Bytes())

	iterator := accountStore.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var balance sdk.Coin
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &balance)

		if cb(balance) {
			break
		}
	}

}

func (am AccountKeeper) GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {

	balances := sdk.NewCoins()
	am.IterateAccountBalances(ctx, addr, func(balance sdk.Coin) bool {
		balances = balances.Add(sdk.NewCoins(balance))
		return false
	})

	return balances.Sort()
}

//func (ak *AccountKeeper) SetSequence(ctx sdk.Context, addr sdk.AccAddress, nonce uint64) sdk.Error {
//	//err := ak.SetSequence(ctx, addr, nonce)
//	//if err != nil {
//	//	return err
//	//}
//	return nil
//}