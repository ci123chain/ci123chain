package wasm

import (
	"github.com/tanhuiya/ci123chain/pkg/wasm/keeper"
	"github.com/tanhuiya/ci123chain/pkg/wasm/types"
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

	NewStoreCodeTx = types.NewStoreCodeTx
	NewInstantiateTx = types.NewInstantiateContractTx
	NewExecuteTx = types.NewExecuteContractTx
)