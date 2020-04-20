package distribution

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/distribution/client/rest"
	"github.com/ci123chain/ci123chain/pkg/distribution/keeper"
	"github.com/ci123chain/ci123chain/pkg/distribution/types"
	"strconv"
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