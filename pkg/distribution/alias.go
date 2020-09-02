package distribution

import (
	h "github.com/ci123chain/ci123chain/pkg/distribution/handler"
	"github.com/ci123chain/ci123chain/pkg/distribution/keeper"
	"github.com/ci123chain/ci123chain/pkg/distribution/types"
)

const (
	DefaultCodespace = types.DefaultParamspace
	ModuleName = types.ModuleName
	RouteKey = types.ModuleName
	//block int64 = 100
	ModuleHeight int64 = 1
)

var (
	ModuleCdc 	= types.DistributionCdc
	NewQuerier = keeper.NewQuerier
	NewHandler = h.NewHandler
)

type (
	GenesisState 	= types.GenesisState
)