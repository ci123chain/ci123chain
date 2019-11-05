package keeper

import (
	"github.com/tanhuiya/ci123chain/pkg/abci/codec"
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/account/exported"
	acc_types "github.com/tanhuiya/ci123chain/pkg/account/types"
)

type AccountKeeper struct {
	key types.StoreKey

	proto func() exported.Account
	cdc *codec.Codec
}

func NewAccountKeeper(cdc *codec.Codec, key types.StoreKey, proto func() exported.Account) AccountKeeper {
	return AccountKeeper{
		key: 	key,
		cdc: 	cdc,
		proto: 	proto,
	}
}

// 利用一个地址 新建一个账户
func (ak AccountKeeper) NewAccountWithAddress(ctx types.Context, addr types.AccAddress) exported.Account {
	acc := ak.proto()
	err := acc.SetAddress(addr)
	if err != nil {
		panic(err)
	}
	return ak.NewAccount(ctx, acc)
}

func (ak AccountKeeper) NewAccount(ctx types.Context, acc exported.Account) exported.Account {
	if err := acc.SetAccountNumber(ak.GetNextAccountNumber(ctx)); err != nil {
		panic(err)
	}
	return acc
}

func (ak AccountKeeper) SetAccount(ctx types.Context, acc exported.Account) {
	addr := acc.GetAddress()
	store := ctx.KVStore(ak.key)
	//bz, err := ak.cdc.MarshalBinaryBare(acc)
	bz, err := ak.cdc.MarshalBinaryLengthPrefixed(acc)
	if err != nil {
		panic(err)
	}
	store.Set(acc_types.AddressStoreKey(addr), bz)
}

func (ak AccountKeeper) GetAccount(ctx types.Context, addr types.AccAddress) exported.Account {
	store := ctx.KVStore(ak.key)
	bz := store.Get(acc_types.AddressStoreKey(addr))
	if bz == nil {
		return nil
	}
	acc := ak.decodeAccount(bz)
	return acc
}


func (ak AccountKeeper) decodeAccount(bz []byte) (acc exported.Account) {
	err := ak.cdc.UnmarshalBinaryLengthPrefixed(bz, &acc)
	if err != nil {
		panic(err)
	}
	return
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
