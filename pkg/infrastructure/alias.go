package infrastructure

import (
	"github.com/ci123chain/ci123chain/pkg/infrastructure/handler"
	"github.com/ci123chain/ci123chain/pkg/infrastructure/keeper"
	"github.com/ci123chain/ci123chain/pkg/infrastructure/types"
)

const (
	DefaultCodeSpce = types.DefaultParamspace
	ModuleName = types.ModuleName
	StoreKey   = types.ModuleName
	RouteKey = types.ModuleName
)


var (

	DefaultGenesisState = types.DefaultGenesisState
	RegisterCodec = types.RegisterCodec

	ModuleCdc = types.InfrastructureCdc
	NewKeeper = keeper.NewInfrastructureKeeper
	NewHandler = handler.NewHandler
	NewQuerier = keeper.NewQuerier

	NewStoreContentMsg = types.NewMsgStoreContent

	AttributeValueStoreContent       = types.EventStoreContent
)

type (
	GenesisState = types.GenesisState
	Keeper = keeper.InfrastructureKeeper
)