package ibc

import (
	"github.com/ci123chain/ci123chain/pkg/ibc/core"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/host"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/keeper"
)

const RouterKey = host.RouterKey


var NewQuerier = keeper.NewQuerier
var NewHandler = core.NewHandler