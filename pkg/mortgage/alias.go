package mortgage

import (
	"github.com/tanhuiya/ci123chain/pkg/mortgage/keeper"
	"github.com/tanhuiya/ci123chain/pkg/mortgage/types"
)

type MortgageKeeper = keeper.MortgageKeeper

var (
	RegisterCodec = types.RegisterCodec
	StoreKey 	  = types.ModuleName
	ModuleName 	  = types.ModuleName

	NewKeeper 	  =  keeper.NewMortgageKeeper

	RouterKey 	  = types.RouterKey
)