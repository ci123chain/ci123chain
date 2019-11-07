package keeper

import (
	"github.com/tanhuiya/ci123chain/pkg/abci/codec"
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/params/subspace"
	auth_types "github.com/tanhuiya/ci123chain/pkg/auth/types"
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


