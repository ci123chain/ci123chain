package registry

import (
	"github.com/ci123chain/ci123chain/pkg/registry/keeper"
	"github.com/ci123chain/ci123chain/pkg/registry/types"
)

const (
	RouteKey = types.RouteKey
	StoreKey = types.StoreKey
	ModuleName = types.ModuleName
	DefaultCodespace = types.DefaultCodespace
)


var (
	//NewQuerier = keeper.NewQuerier
	//NewHandler = h.NewHandler
	NewKeeper = keeper.NewKeeper
)

type Keeper = keeper.Keeper
