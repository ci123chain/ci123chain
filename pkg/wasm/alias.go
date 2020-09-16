package wasm

import (
	"github.com/ci123chain/ci123chain/pkg/wasm/keeper"
	"github.com/ci123chain/ci123chain/pkg/wasm/types"
)

const (
	ModuleName = types.ModuleName
	DefaultCodespace = types.DefaultCodespace
	StoreKey = types.StoreKey
	RouteKey = types.RouteKey
)

var (
	Modulecdc = types.WasmCodec
	NewKeeper = keeper.NewKeeper
	NewHandler = keeper.NewHandler
	NewQuerier = keeper.NewQuerier

	NewInstantiateTx = types.NewMsgInstantiateContract
	NewExecuteTx = types.NewMsgExecuteContract
	NewMigrateTx = types.NewMsgMigrateContract
)

type (
	GenesisState = types.GenesisState

)