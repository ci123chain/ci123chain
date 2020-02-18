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

	ModuleName = types.ModuleName

	RegisterCodec = types.RegisterCodec
)

const (
	Burner   = types.Burner
	Staking  = types.Staking
)