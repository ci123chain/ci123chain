package slashing

import (
	"github.com/ci123chain/ci123chain/pkg/slashing/types"
	k "github.com/ci123chain/ci123chain/pkg/slashing/keeper"
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