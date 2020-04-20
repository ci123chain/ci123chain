package order

import (
	"github.com/ci123chain/ci123chain/pkg/order/client/rest"
	"github.com/ci123chain/ci123chain/pkg/order/keeper"
	"github.com/ci123chain/ci123chain/pkg/order/types"
)

var (
	NewQuerier 				= keeper.NewQuerier
	NewKeeper  				= keeper.NewOrderKeeper
	RegisterTxRoutes 		= rest.RegisterTxRoutes
	ErrQueryTx				= types.ErrQueryTx
	NewAddShardTx           = types.NewUpgradeTx
	StoreKey				= types.StoreKey
)

const (
	ModuleName = types.ModuleName
	RouteKey = types.ModuleName
)