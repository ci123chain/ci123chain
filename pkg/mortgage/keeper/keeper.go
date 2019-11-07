package keeper

import (
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/supply"
)

type MortgageKeeper struct {
	StoreKey 	sdk.StoreKey

	SupplyKeeper supply.Keeper
}

func NewMortgageKeeper(key sdk.StoreKey, supplyKeeper supply.Keeper) MortgageKeeper {
	return MortgageKeeper{
		StoreKey: 	key,
		SupplyKeeper: supplyKeeper,
	}
}