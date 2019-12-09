package distribution

import (
	"github.com/tanhuiya/ci123chain/pkg/distribution/types"
	"github.com/tanhuiya/ci123chain/pkg/distribution/client/rest"
	"github.com/tanhuiya/ci123chain/pkg/distribution/keeper"
)

const (
	ModuleName = types.ModuleName
	RouteKey = types.ModuleName
)

var (
	ModuleCdc 	= types.DistributionCdc


	RegisterRoutes = rest.RegisterTxRoutes
	NewQuerier = keeper.NewQuerier
)

type (
	GenesisState 	= types.GenesisState
)