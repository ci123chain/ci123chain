package slashing

import (
	k "github.com/ci123chain/ci123chain/pkg/slashing/keeper"
	"github.com/ci123chain/ci123chain/pkg/slashing/types"
)

const (
	StoreKey = types.StoreKey
	RouteKey = types.RouterKey
	ModuleName = types.ModuleName
)

type (
	Keeper = k.Keeper
)

var (
	NewKeeper = k.NewKeeper
	NewQuerier = k.NewQuerier
)