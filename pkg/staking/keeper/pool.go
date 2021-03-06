package keeper

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/staking/types"
	"github.com/ci123chain/ci123chain/pkg/supply/exported"
)

func (k StakingKeeper) bondedTokensToNotBonded(ctx sdk.Context, tokens sdk.Int) error {

	//coins := sdk.NewCoins(sdk.NewCoin(tokens))
	coin := sdk.NewCoin(tokens)
	err := k.SupplyKeeper.SendCoinsFromModuleToModule(ctx, types.BondedPoolName, types.NotBondedPoolName, coin)
	if err != nil {
		return err
	}
	return nil
}

func (k StakingKeeper) notBondedTokensToBonded(ctx sdk.Context, tokens sdk.Int) error {

	//coins:= sdk.NewCoins(sdk.NewCoin(tokens))
	coin := sdk.NewCoin(tokens)
	err := k.SupplyKeeper.SendCoinsFromModuleToModule(ctx, types.NotBondedPoolName, types.BondedPoolName, coin)
	if err != nil {
		return err
	}
	return nil
}

// GetBondedPool returns the bonded tokens pool's module account
func (k StakingKeeper) GetBondedPool(ctx sdk.Context) (bondedPool exported.ModuleAccountI) {
	return k.SupplyKeeper.GetModuleAccount(ctx, types.BondedPoolName)
}


// GetNotBondedPool returns the not bonded tokens pool's module account
func (k StakingKeeper) GetNotBondedPool(ctx sdk.Context) (notBondedPool exported.ModuleAccountI) {
	return k.SupplyKeeper.GetModuleAccount(ctx, types.NotBondedPoolName)
}

// TotalBondedTokens total staking tokens supply which is bonded
func (k StakingKeeper) TotalBondedTokens(ctx sdk.Context) sdk.Int {
	bondedPool := k.GetBondedPool(ctx)
	return k.AccountKeeper.GetBalance(ctx, bondedPool.GetAddress()).Amount
}

// StakingTokenSupply staking tokens from the total supply
func (k StakingKeeper) StakingTokenSupply(ctx sdk.Context) sdk.Int {
	return k.SupplyKeeper.GetSupply(ctx).GetTotal().AmountOf(k.BondDenom(ctx))
}

// BondedRatio the fraction of the staking tokens which are currently bonded
func (k StakingKeeper) BondedRatio(ctx sdk.Context) sdk.Dec {
	stakeSupply := k.StakingTokenSupply(ctx)
	if stakeSupply.IsPositive() {
		return k.TotalBondedTokens(ctx).ToDec().QuoInt(stakeSupply)
	}

	return sdk.ZeroDec()
}