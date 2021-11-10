package vm

import (
	"github.com/ci123chain/ci123chain/pkg/vm/evmtypes"
	"github.com/ci123chain/ci123chain/pkg/vm/keeper"
	"github.com/ci123chain/ci123chain/pkg/vm/moduletypes"
	"github.com/ci123chain/ci123chain/pkg/vm/wasmtypes"
)

const (
	ModuleName = moduletypes.ModuleName
	DefaultCodespace = moduletypes.DefaultCodespace
	StoreKey = moduletypes.StoreKey
	RouteKey = moduletypes.RouteKey
)

var (
	NewKeeper = keeper.NewKeeper
	NewHandler = keeper.NewHandler
	NewQuerier = keeper.NewQuerier

	NewMsgEvmTx = evmtypes.NewMsgEvmTx
	NewUploadTx = types.NewMsgUploadContract
	NewInstantiateTx = types.NewMsgInstantiateContract
	NewExecuteTx = types.NewMsgExecuteContract
	NewMigrateTx = types.NewMsgMigrateContract
)

type (
	Keeper = keeper.Keeper
)
