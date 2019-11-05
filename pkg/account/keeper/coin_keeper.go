package keeper

import (
	"fmt"
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
)

func (ak AccountKeeper) GetBalance(ctx types.Context, addr types.AccAddress) types.Coin {
	coin := ak.getBalance(ctx, addr)
	return coin
}

func (ak AccountKeeper) AddBalance(ctx types.Context, addr types.AccAddress, amount types.Coin) (types.Coin, types.Error) {
	if !amount.IsValid() {
		return types.NewCoin(), types.ErrInvalidCoins(amount.String())
	}
	oldCoin := ak.GetBalance(ctx, addr)
	newCoin := oldCoin.Add(amount)

	if newCoin.IsAnyNegative() {
		return amount, types.ErrInsufficientCoins(
			fmt.Sprintf("insufficient account funds: %s < %s", oldCoin, amount),
		)
	}
	err := ak.SetCoin(ctx, addr, newCoin)
	return newCoin, err
}

func (ak AccountKeeper) SetCoin(ctx types.Context, addr types.AccAddress, amt types.Coin) types.Error {
	if !amt.IsValid() {
		return types.ErrInvalidCoins(amt.String())
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

func (ak AccountKeeper) SubBalance(ctx types.Context, addr types.AccAddress, amt types.Coin) (types.Coin, types.Error) {
	if !amt.IsValid() {
		return types.NewCoin(), types.ErrInvalidCoins(amt.String())
	}

	oldCoins, spendableCoins := types.NewCoin(), types.NewCoin()

	acc := ak.GetAccount(ctx, addr)
	if acc != nil {
		oldCoins = acc.GetCoin()
		spendableCoins = acc.SpendableCoins(ctx.BlockHeader().Time)
	}
	_, valid := spendableCoins.SafeSub(amt)
	if !valid {
		return amt, types.ErrInsufficientCoins(
			fmt.Sprintf("insufficient accounts funds; %s < %s", spendableCoins, amt),
		)
	}

	newCoin := oldCoins.Sub(amt)
	err := ak.SetCoin(ctx, addr, newCoin)
	return newCoin, err
}


func (ak AccountKeeper) Transfer(ctx types.Context, fromAddr types.AccAddress, toAddr types.AccAddress, coin types.Coin) types.Error {
	_, err := ak.SubBalance(ctx, fromAddr, coin)
	if err != nil {
		return err
	}

	_, err = ak.AddBalance(ctx, toAddr, coin)
	if err != nil {
		return err
	}

	return nil
}

func (ak *AccountKeeper) getStore(ctx types.Context) types.KVStore {
	return ctx.KVStore(ak.key)
}


func (am *AccountKeeper) getBalance(ctx types.Context, addr types.AccAddress) types.Coin {
	acc := am.GetAccount(ctx, addr)
	if acc == nil {
		return types.NewCoin()
	}
	return acc.GetCoin()

}