package gravity

import (
	"github.com/ci123chain/ci123chain/pkg/gravity/types"
	k "github.com/ci123chain/ci123chain/pkg/gravity/keeper"
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
