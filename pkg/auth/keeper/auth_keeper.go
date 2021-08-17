package keeper

import (
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	auth_types "github.com/ci123chain/ci123chain/pkg/auth/types"
	subspace "github.com/ci123chain/ci123chain/pkg/params/types"
	"math/big"
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

func (ak AuthKeeper)SetNumTxs(ctx types.Context) {
	ctx.KVStore(ak.key).Set(auth_types.KeyNumTxs, big.NewInt(ctx.BlockHeader().NumTxs).Bytes())
}

func (ak AuthKeeper)GetNumTxs(ctx types.Context) int64{
	numTxs := new(big.Int)
	numTxs.SetBytes(ctx.KVStore(ak.key).Get(auth_types.KeyNumTxs))
	return numTxs.Int64()
}