package params

import (
	"github.com/tanhuiya/ci123chain/pkg/params/subspace"
	"github.com/tanhuiya/ci123chain/pkg/params/types"
)

const (
	StoreKey 		= subspace.StoreKey
	TStoreKey 		= subspace.TStoreKey
	DefaultCodespace= types.DefaultCodespace
	ModuleName 		= types.ModuleName
)

type (
	Subspace 		= subspace.Subspace
)