package account

import "github.com/tanhuiya/ci123chain/pkg/account/types"

const (
	ModuleName = types.ModuleName
)

var (
	SetGenesisStateInAppState 	= types.SetGenesisStateInAppState
	NewGenesisAccountRaw 		= types.NewGenesisAccountRaw

	ModuleCdc 					= types.ModuleCdc

)

type (
	GenesisState 	= types.GenesisState
)
