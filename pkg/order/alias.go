package order

import (
	"github.com/tanhuiya/ci123chain/pkg/order/client/rest"
	"github.com/tanhuiya/ci123chain/pkg/order/keeper"
	"github.com/tanhuiya/ci123chain/pkg/order/types"
)

var (
	NewQuerier 				= keeper.NewQuerier
	NewKeeper  				= keeper.NewOrderKeeper
	RegisterTxRoutes 		= rest.RegisterTxRoutes

	ErrQueryTx				= types.ErrQueryTx

	StoreKey				= types.StoreKey
)

const (
	ModuleName = types.ModuleName
	RouteKey = types.ModuleName
)