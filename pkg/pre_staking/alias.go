package pre_staking

import (
	h "github.com/ci123chain/ci123chain/pkg/pre_staking/handler"
	k "github.com/ci123chain/ci123chain/pkg/pre_staking/keeper"
	"github.com/ci123chain/ci123chain/pkg/pre_staking/types"
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
	Keeper = k.PreStakingKeeper
)