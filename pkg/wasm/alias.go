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
	NewKeeper = keeper.NewKeeper
	NewHandler = keeper.NewHandler
	NewQuerier = keeper.NewQuerier

	NewInstantiateTx = types.NewInstantiateContractTx
	NewExecuteTx = types.NewExecuteContractTx
	NewMigrateTx = types.NewMigrateContractTx
)