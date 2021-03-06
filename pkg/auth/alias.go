package auth

import (
	"github.com/ci123chain/ci123chain/pkg/auth/keeper"
	"github.com/ci123chain/ci123chain/pkg/auth/types"
)

const (
	DefaultCodespace = types.DefaultParamspace
	StoreKey 		 = types.StoreKey

	FeeCollectorName              = types.FeeCollectorName
)

var (
	ErrTxValidateBasic = types.ErrTxValidateBasic
	NewAuthKeeper = keeper.NewAuthKeeper
)

type (
	AuthKeeper 	  = keeper.AuthKeeper
	GenesisState  = types.GenesisState
	Params        = types.Params
)