package keeper

import (
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	auth_types "github.com/ci123chain/ci123chain/pkg/auth/types"
	"github.com/ci123chain/ci123chain/pkg/params/subspace"
)

type AuthKeeper struct {
	key 	types.StoreKey

	cdc 	*codec.Codec
	paramSubspace subspace.Subspace
}

func NewAuthKeeper(cdc *codec.Codec, key types.StoreKey, paramstore subspace.Subspace) AuthKeeper {
	return AuthKeeper{
		key:           key,
		cdc:           cdc,
		paramSubspace: paramstore.WithKeyTable(auth_types.ParamKeyTable()),
	}
}


