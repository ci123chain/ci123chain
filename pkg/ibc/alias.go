package ibc

import (
	"github.com/tanhuiya/ci123chain/pkg/ibc/client/rest"
	"github.com/tanhuiya/ci123chain/pkg/ibc/handler"
	"github.com/tanhuiya/ci123chain/pkg/ibc/keeper"
	"github.com/tanhuiya/ci123chain/pkg/ibc/types"
)

var (
	StoreKey  = types.StoreKey
	RouterKey  = types.RouterKey

	NewHandler = handler.NewHandler
	NewKeeper = keeper.NewIBCKeeper

	NewIBCTransfer = types.NewIBCTransferMsg

	RegisterCodec = types.RegisterCodec

	RegisterRoutes = rest.RegisterTxRoutes
)

type IBCMsg types.IBCMsg 

