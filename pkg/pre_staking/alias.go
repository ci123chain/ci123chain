package pre_staking

import (
	"github.com/ci123chain/ci123chain/pkg/pre_staking/types"
	k "github.com/ci123chain/ci123chain/pkg/pre_staking/keeper"
	h "github.com/ci123chain/ci123chain/pkg/pre_staking/handler"
)

const (
	RouteKey = types.RouteKey
	StoreKey = types.StoreKey
	ModuleName = types.ModuleName
	DefaultCodespace = types.DefaultCodespace
)


var (
	NewQuerier = k.NewQuerier
	NewHandler = h.NewHandler
	NewKeeper = k.NewPreStakingKeeper
)

type (
	GenesisState = types.GenesisState
)