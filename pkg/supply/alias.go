package supply

import (
	"github.com/ci123chain/ci123chain/pkg/supply/keeper"
	"github.com/ci123chain/ci123chain/pkg/supply/types"
)

type (
	Keeper = keeper.Keeper
	GenesisState = types.GenesisState
)

var (
	NewKeeper = keeper.NewKeeper
	StoreKey  = types.ModuleName

	ModuleName = types.ModuleName

	RegisterCodec = types.RegisterCodec

	ModuleCdc   = types.ModuleCdc
)

const (
	Minter   = types.Minter
	Burner   = types.Burner
	Staking  = types.Staking
)