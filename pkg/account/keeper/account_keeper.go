package keeper

import (
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/account/exported"
	acc_types "github.com/ci123chain/ci123chain/pkg/account/types"
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
	if acc.String() == "0x4feA76427B8345861e80A3540a8a9D936FD39391" {
		fmt.Println(acc.GetIsModule())
	}
	addr := acc.GetAddress()
	store := ctx.KVStore(ak.key)
	//bz, err := ak.cdc.MarshalBinaryBare(acc)
	bz, err := ak.cdc.MarshalBinaryLengthPrefixed(acc)
	if err != nil {
		panic(err)
	}
	store.Set(acc_types.AddressStoreKey(addr), bz)

	//coins := acc.GetCoins()
	//chain_id := ctx.ChainID()
	//update := util.HeightUpdate{
	//	Shard: chain_id,
	//	Coins: coins,
	//}
	//h := ctx.BlockHeight()
	//v, err := ak.cdc.MarshalBinaryLengthPrefixed(update)
	//if err != nil {
	//	panic(err)
	//}
	//
	//var heights util.Heights
	//b := store.Get(acc_types.HeightsUpdateKey(addr))
	//if b == nil {
	//	heights = make(util.Heights, 0)
	//	heights = append(heights, h)
	//}else {
	//	err := ak.cdc.UnmarshalBinaryLengthPrefixed(b, &heights)
	//	if err != nil {
	//		panic(err)
	//	}
	//	sort.Sort(heights)
	//	if heights[len(heights)-1] < h{
	//		heights = append(heights, h)
	//	}
	//}
	//hbz, err := ak.cdc.MarshalBinaryLengthPrefixed(heights)
	//if err != nil {
	//	panic(err)
	//}
	//
	//store.Set(acc_types.HeightUpdateKey(addr, h), v)
	//store.Set(acc_types.HeightsUpdateKey(addr), hbz)
}

//func SearchHeight(ctx types.Context, ak AccountKeeper, acc exported.Account, i int64) int64 {
//	addr := acc.GetAddress()
//	store := ctx.KVStore(ak.key)
//	var heights util.Heights
//	b := store.Get(acc_types.HeightsUpdateKey(addr))
//	err := ak.cdc.UnmarshalBinaryLengthPrefixed(b, &heights)
//	if err != nil {
//		return -3
//	}else {
//		return heights.Search(i)
//	}
//}

//func GetHistoryBalance(ctx types.Context, ak AccountKeeper, acc exported.Account, i int64) util.HeightUpdate {
//	addr := acc.GetAddress()
//	store := ctx.KVStore(ak.key)
//	var update util.HeightUpdate
//	b := store.Get(acc_types.HeightUpdateKey(addr, i))
//	err := ak.cdc.UnmarshalBinaryLengthPrefixed(b, & update)
//	if err != nil {
//		panic(err)
//	}
//	return update
//}

func (ak AccountKeeper) GetAccount(ctx types.Context, addr types.AccAddress) exported.Account {
	store := ctx.KVStore(ak.key)
	bz := store.Get(acc_types.AddressStoreKey(addr))
	if bz == nil {
		return nil
	}
	acc := ak.decodeAccount(bz)
	return acc
}

// GetAllAccounts returns all accounts in the accountKeeper.
func (ak AccountKeeper) GetAllAccounts(ctx types.Context) (accounts []exported.Account) {
	ak.IterateAccounts(ctx,
		func(acc exported.Account) (stop bool) {
			accounts = append(accounts, acc)
			return false
		})
	return accounts
}

func (ak AccountKeeper) RemoveAccount(ctx types.Context, acc exported.Account) {
	addr := acc.GetAddress()
	store := ctx.KVStore(ak.key)
	store.Delete(acc_types.AddressStoreKey(addr))
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

// IterateAccounts iterates over all the stored accounts and performs a callback function
func (ak AccountKeeper) IterateAccounts(ctx types.Context, cb func(account exported.Account) (stop bool)) {
	store := ctx.KVStore(ak.key)
	//
	//bz := store.Get(acc_types.AddressStoreKey(types.HexToAddress("0x3F43E75Aaba2c2fD6E227C10C6E7DC125A93DE3c")))
	//acc := ak.decodeAccount(bz)
	//fmt.Println(acc.GetAddress().String())
	prefix := acc_types.AddressStoreKeyPrefix
	//iterator := store.RemoteIterator(prefix, types.PrefixEndBytes(prefix))
	iterator := types.KVStorePrefixIterator(store, prefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		account := ak.decodeAccount(iterator.Value())

		if cb(account) {
			break
		}
	}
}