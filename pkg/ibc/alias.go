package ibc

import (
	"github.com/ci123chain/ci123chain/pkg/ibc/application/transfer"
	"github.com/ci123chain/ci123chain/pkg/ibc/application/transfer/client/rest"
	transferkeeper "github.com/ci123chain/ci123chain/pkg/ibc/application/transfer/keeper"
	"github.com/ci123chain/ci123chain/pkg/ibc/core"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/keeper"
)

const ModuleName = "ibc"

const RouterKey = ModuleName


var NewHandler = core.NewHandler
var NewQuerier = keeper.NewQuerier

var NewCoreModule = core.NewAppModule
var NewTransferModule = transfer.NewAppModule
var NewTransferQuerier = transferkeeper.NewQuerier
var RegisterRoutes = rest.RegisterRoutes

type Keeper = keeper.Keeper
