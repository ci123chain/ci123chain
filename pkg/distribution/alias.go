package distribution

import (
	"github.com/tanhuiya/ci123chain/pkg/distribution/types"
	"github.com/tanhuiya/ci123chain/pkg/distribution/client/rest"
	"github.com/tanhuiya/ci123chain/pkg/distribution/keeper"
	"strconv"
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
)

const (
	ModuleName = types.ModuleName
	RouteKey = types.ModuleName
	block int64 = 100
	ModuleHeight int64 = 1
)

var (
	ModuleCdc 	= types.DistributionCdc


	RegisterRoutes = rest.RegisterTxRoutes
	NewQuerier = keeper.NewQuerier
)

type (
	GenesisState 	= types.GenesisState
)

type lastCommitValidatorsAddr struct {
	Address []string `json:"address"`
}

func getKey(addr string, height int64) sdk.AccAddr {

	tkey := strconv.FormatInt(height, 10)
	key := sdk.AccAddr([]byte(addr + tkey))
	return key
}