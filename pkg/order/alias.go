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
	//ErrQueryTx				= types.ErrQueryTx
	NewMsgUpgrade           = types.NewMsgUpgrade

	//EventType                = types.EventType
	//AttributeValueCategory  = types.AttributeValueCategory
	//AttributeValueAddShard  = types.AttributeValueAddShard
)

const (
	StoreKey  = types.ModuleName
	ModuleName = types.ModuleName
	RouteKey = types.ModuleName
)