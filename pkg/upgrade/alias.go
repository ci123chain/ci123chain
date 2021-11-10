package upgrade

import (
	"github.com/ci123chain/ci123chain/pkg/upgrade/keeper"
	"github.com/ci123chain/ci123chain/pkg/upgrade/types"
)

const (
	ModuleName                        = types.ModuleName
	RouterKey                         = types.RouterKey
	StoreKey                          = types.StoreKey
	QuerierKey                        = types.QuerierKey
	PlanByte                          = types.PlanByte
	DoneByte                          = types.DoneByte
)

var (
	// functions aliases
	RegisterCodec                    = types.RegisterCodec
	PlanKey                          = types.PlanKey

	NewKeeper                        = keeper.NewKeeper
	NewQuerier                       = keeper.NewQuerier
)


type (
	UpgradeHandler                = types.UpgradeHandler //nolint:golint

	Plan                          = types.Plan
	Keeper                        = keeper.Keeper
)

