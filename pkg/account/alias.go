package account

import (
	"github.com/ci123chain/ci123chain/pkg/account/keeper"
	"github.com/ci123chain/ci123chain/pkg/account/types"
)

const (
	ModuleName = types.ModuleName
	RouteKey = types.RouteKey
	StoreKey = types.StoreKey
)

var (
	//SetGenesisStateInAppState 	= types.SetGenesisStateInAppState
	//NewGenesisAccountRaw 		= types.NewGenesisAccountRaw
	ModuleCdc 					= types.ModuleCdc
	NewQuerier					= keeper.NewQuerier
	NewGensisState              = types.NewGenesisState
	//ErrSetAccount				= types.ErrSetAccount
	//ErrGetAccount				= types.ErrGetAccount
)

type (
	GenesisState 	= types.GenesisState
	BaseAccount 	= types.BaseAccount
	AccountKeeper 	= keeper.AccountKeeper
	GenesisAccounts = types.GenesisAccounts
	//GenesisAccount  = types.GenesisAccount
)
