package account

import (
	"github.com/ci123chain/ci123chain/pkg/account/keeper"
	"github.com/ci123chain/ci123chain/pkg/account/types"
)

const (
	ModuleName = types.ModuleName
)

var (
	SetGenesisStateInAppState 	= types.SetGenesisStateInAppState
	NewGenesisAccountRaw 		= types.NewGenesisAccountRaw
	ModuleCdc 					= types.ModuleCdc

	ErrSetAccount				= types.ErrSetAccount
)

type (
	GenesisState 	= types.GenesisState
	BaseAccount 	= types.BaseAccount
	AccountKeeper 	= keeper.AccountKeeper
)
