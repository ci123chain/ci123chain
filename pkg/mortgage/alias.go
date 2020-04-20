package mortgage

import (
	"github.com/ci123chain/ci123chain/pkg/mortgage/keeper"
	"github.com/ci123chain/ci123chain/pkg/mortgage/types"
)

type MortgageKeeper = keeper.MortgageKeeper

var (
	RegisterCodec = types.RegisterCodec
	StoreKey 	  = types.ModuleName
	ModuleName 	  = types.ModuleName

	NewKeeper 	  =  keeper.NewMortgageKeeper

	RouterKey 	  = types.RouterKey

	NewMortgageMsg = types.NewMsgMortgage
	NewMsgMortgageCancel = types.NewMsgMortgageCancel
	NewMsgMortgageDone  = types.NewMsgMortgageDone
)

