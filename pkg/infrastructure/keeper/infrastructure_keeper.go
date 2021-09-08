package keeper

import (
	"errors"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
)

type InfrastructureKeeper struct {
	cdc     		   *codec.Codec
	storeKey  		   sdk.StoreKey
}


func NewInfrastructureKeeper(cdc *codec.Codec, key sdk.StoreKey) InfrastructureKeeper {

	return InfrastructureKeeper{
		cdc:           cdc,
		storeKey:      key,
	}
}

func (k InfrastructureKeeper) GetStoreKey() sdk.StoreKey {
	return k.storeKey
}


//get method
func(k InfrastructureKeeper) GetContent(ctx sdk.Context, key []byte) ([]byte, error) {
	store := ctx.KVStore(k.storeKey)
	res := store.Get(key)
	if res == nil {
		return nil, errors.New("no content found")
	}
	return res, nil
}


//set method
func(k InfrastructureKeeper) SetContent(ctx sdk.Context, key []byte, value []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Set(key, value)
}