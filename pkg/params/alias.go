package params

import (
	"github.com/ci123chain/ci123chain/pkg/params/subspace"
	"github.com/ci123chain/ci123chain/pkg/params/types"
)

const (
	StoreKey 		= subspace.StoreKey
	TStoreKey 		= subspace.TStoreKey
	DefaultCodespace= types.DefaultCodespace
	ModuleName 		= types.ModuleName
)

var (
	NewParamSetPair = subspace.NewParamSetPair
	NewKeyTable     = subspace.NewKeyTable
)

type (
	Subspace 		= subspace.Subspace
	ParamSetPair  = subspace.ParamSetPair
	ParamSetPairs = subspace.ParamSetPairs
	KeyTable      = subspace.KeyTable
)