package keeper

import (
	"github.com/tanhuiya/ci123chain/pkg/abci/codec"
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
	acc_types "github.com/tanhuiya/ci123chain/pkg/account/types"
)

type AccountKeeper struct {
	key types.StoreKey

	cdc *codec.Codec
}

func NewAccountKeeper(cdc *codec.Codec, key types.StoreKey) AccountKeeper {
	return AccountKeeper{
		key: 	key,
		cdc: 	cdc,
	}
}

func (ak AccountKeeper) NewAccount(ctx types.Context, acc acc_types.BaseAccount) acc_types.BaseAccount {
	if err := acc.SetAccountNumber(ak.GetNextAccountNumber(ctx)); err != nil {
		panic(err)
	}
	return acc
}

func (ak AccountKeeper) SetAccount(ctx types.Context, acc acc_types.BaseAccount) {
	addr := acc.GetAddress()
	store := ctx.KVStore(ak.key)
	bz, err := ak.cdc.MarshalBinaryBare(acc)
	if err != nil {
		panic(err)
	}
	store.Set(acc_types.AddressStoreKey(addr), bz)
}

func (ak AccountKeeper) GetNextAccountNumber(ctx types.Context) uint64 {
	var accNumber uint64
	store := ctx.KVStore(ak.key)
	bz := store.Get(acc_types.GlobalAccountNumberKey)
	if bz == nil {
		accNumber = 0
	} else {
		err := ak.cdc.UnmarshalBinaryLengthPrefixed(bz, &accNumber)
		if err != nil {
			panic(err)
		}
	}
	bz = ak.cdc.MustMarshalBinaryLengthPrefixed(accNumber + 1)
	store.Set(acc_types.GlobalAccountNumberKey, bz)
	return accNumber
}