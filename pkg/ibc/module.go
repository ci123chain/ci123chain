package ibc

import (
	"github.com/ci123chain/ci123chain/pkg/ibc/application/transfer"
	"github.com/ci123chain/ci123chain/pkg/ibc/core"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/keeper"
	transferkeeper "github.com/ci123chain/ci123chain/pkg/ibc/application/transfer/keeper"
)

const ModuleName = "ibc"

const RouterKey = ModuleName

var NewHandler = core.NewHandler
var NewQuerier = keeper.NewQuerier

var NewCoreModule = core.NewAppModule
var NewTransferModule = transfer.NewAppModule
var NewTransferQuerier = transferkeeper.NewQuerier