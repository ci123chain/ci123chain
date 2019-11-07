package auth

import (
	"github.com/tanhuiya/ci123chain/pkg/auth/keeper"
	"github.com/tanhuiya/ci123chain/pkg/auth/types"
)

const (
	DefaultCodespace = types.DefaultParamspace
	StoreKey 		 = types.StoreKey
)

var (
	NewAuthKeeper = keeper.NewAuthKeeper
)

type (
	AuthKeeper 	  = keeper.AuthKeeper
	GenesisState  = types.GenesisState
	Params        = types.Params
)