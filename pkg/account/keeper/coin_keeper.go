package keeper

import (
	"fmt"
	prefix "github.com/ci123chain/ci123chain/pkg/abci/store"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/account/types"
	"github.com/ci123chain/ci123chain/pkg/transaction"
)

func (ak AccountKeeper) GetBalance(ctx sdk.Context, addr sdk.AccAddress) sdk.Coin {
	coin := ak.getBalance(ctx, addr)
	if coin.IsZero() {
		return sdk.NewEmptyCoin()
	}
	return coin
}

func (ak AccountKeeper) AddBalance(ctx sdk.Context, addr sdk.AccAddress, amount sdk.Coin) (sdk.Coin, sdk.Error) {
	if !amount.IsValid() {
		return sdk.NewCoin(sdk.NewInt(0)), sdk.ErrInvalidCoins(amount.String())
	}
	oldCoin := ak.GetBalance(ctx, addr)
	newCoin := oldCoin.Add(amount)

	if newCoin.IsNegative() {
		return amount, sdk.ErrInsufficientCoins(
			fmt.Sprintf("insufficient account funds: %s < %s", oldCoin, amount),
		)
	}
	err := ak.SetCoin(ctx, addr, newCoin)
	return newCoin, err
}

func (ak AccountKeeper) SetCoin(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coin) sdk.Error {
	if !amt.IsValid() {
		return sdk.ErrInvalidCoins(amt.String())
	}

	acc := ak.GetAccount(ctx, addr)
	if acc == nil {
		acc = ak.NewAccountWithAddress(ctx, addr)
	}

	err := acc.SetCoin(amt)
	if err != nil {
		panic(err)
	}

	ak.SetAccount(ctx, acc)
	return nil
}

func (ak AccountKeeper) SubBalance(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coin) (sdk.Coin, sdk.Error) {
	if !amt.IsValid() {
		return sdk.NewCoin(sdk.NewInt(0)), sdk.ErrInvalidCoins(amt.String())
	}

	oldCoins, spendableCoins := sdk.NewCoin(sdk.NewInt(0)), sdk.NewCoin(sdk.NewInt(0))

	acc := ak.GetAccount(ctx, addr)
	if acc != nil {
		oldCoins = acc.GetCoin()
		spendableCoins = acc.SpendableCoins(ctx.BlockHeader().Time)
	} else {
		return sdk.NewCoin(sdk.NewInt(0)), transaction.ErrInvalidTx(types.DefaultCodespace, fmt.Sprintf("account not exist %s", addr.Hex()))
	}
	_, valid := spendableCoins.SafeSub(amt)
	if !valid {
		return amt, sdk.ErrInsufficientCoins(
			fmt.Sprintf("insufficient accounts funds; %s < %s", spendableCoins, amt),
		)
	}

	newCoin := oldCoins.Sub(amt)
	err := ak.SetCoin(ctx, addr, newCoin)
	return newCoin, err
}


func (ak AccountKeeper) Transfer(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, coin sdk.Coin) sdk.Error {
	_, err := ak.SubBalance(ctx, fromAddr, coin)
	if err != nil {
		return types.ErrSetAccount(types.DefaultCodespace, err)
	}

	_, err = ak.AddBalance(ctx, toAddr, coin)
	if err != nil {
		return types.ErrSetAccount(types.DefaultCodespace, err)
	}

	return nil
}

func (ak AccountKeeper) getStore(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(ak.key)
}


func (k AccountKeeper) getBalance(ctx sdk.Context, addr sdk.AccAddress) sdk.Coin {
	acc := k.GetAccount(ctx, addr)
	if acc == nil {
		return sdk.NewCoin(sdk.NewInt(0))
	}
	return acc.GetCoin()
}

func (k AccountKeeper) SetBalances(ctx sdk.Context, addr sdk.AccAddress, balances sdk.Coins) error {

	k.ClearBalances(ctx, addr)

	for _, balance := range balances {
		err := k.SetCoin(ctx, addr, balance)
		if err != nil {
			return err
		}
	}

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
		balances = balances.Add(balance)
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