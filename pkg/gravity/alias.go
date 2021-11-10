package gravity

import (
	k "github.com/ci123chain/ci123chain/pkg/gravity/keeper"
	"github.com/ci123chain/ci123chain/pkg/gravity/types"
)

const (
	StoreKey = types.StoreKey
	RouteKey = types.RouterKey
	ModuleName = types.ModuleName
)

var (
	NewKeeper = k.NewKeeper
	NewQuerier = k.NewQuerier
)

type (
	Keeper = k.Keeper
)