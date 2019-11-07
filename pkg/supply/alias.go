package supply

import (
	"github.com/tanhuiya/ci123chain/pkg/supply/keeper"
	"github.com/tanhuiya/ci123chain/pkg/supply/types"
)

type (
	Keeper = keeper.Keeper
)

var (
	NewKeeper = keeper.NewKeeper
	StoreKey  = types.ModuleName
)
