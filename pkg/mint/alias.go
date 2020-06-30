package mint

import (
	"github.com/ci123chain/ci123chain/pkg/mint/keeper"
	"github.com/ci123chain/ci123chain/pkg/mint/types"
)


const (
	DefaultCodeSpce = types.DefaultParamspace
	ModuleName = types.ModuleName
	StoreKey   = types.ModuleName
)

var (

	DefaultGenesisState = types.DefaultGenesisState
	RegisterCodec = types.RegisterCodec

	ModuleCdc = types.MintCdc
	NewKeeper = keeper.NewMinterKeeper
)

type (
	Keeper keeper.MinterKeeper
	GenesisState types.GenesisState
)