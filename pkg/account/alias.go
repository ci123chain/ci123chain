package account

import (
	"github.com/ci123chain/ci123chain/pkg/account/keeper"
	"github.com/ci123chain/ci123chain/pkg/account/types"
)

const (
	ModuleName = types.ModuleName
	RouteKey = types.RouteKey
)

var (
	SetGenesisStateInAppState 	= types.SetGenesisStateInAppState
	NewGenesisAccountRaw 		= types.NewGenesisAccountRaw
	ModuleCdc 					= types.ModuleCdc
	NewQuerier					= keeper.NewQuerier
	ErrSetAccount				= types.ErrSetAccount
	ErrGetAccount				= types.ErrGetAccount
)

type (
	GenesisState 	= types.GenesisState
	BaseAccount 	= types.BaseAccount
	AccountKeeper 	= keeper.AccountKeeper
)
