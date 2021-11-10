package supply

import (
	"github.com/ci123chain/ci123chain/pkg/supply/keeper"
	"github.com/ci123chain/ci123chain/pkg/supply/types"
)

type (
	Keeper = keeper.Keeper
	GenesisState = types.GenesisState
)

var  (
	NewKeeper = keeper.NewKeeper
	RegisterCodec = types.RegisterCodec
	ModuleCdc   = types.ModuleCdc
)

const (
	StoreKey  = types.ModuleName
	ModuleName = types.ModuleName
	Minter   = types.Minter
	Burner   = types.Burner
	Staking  = types.Staking
)