package keeper

import (
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/staking/types"
	"github.com/tanhuiya/ci123chain/pkg/supply/exported"
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