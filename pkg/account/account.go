package account

import (
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/util"
)


type Account struct {
	Address types.AccAddress
	Amount  uint64
}

type AccountMapper interface {
	GetBalance(types.Context, types.AccAddress) (uint64, error)
	AddBalance(types.Context, types.AccAddress, uint64) (uint64, error)
	SubBalance(types.Context, types.AccAddress, uint64) (uint64, error)
	Transfer(types.Context, types.AccAddress, uint64, types.AccAddress) error
}

type accountMapper struct {
	storeKey types.StoreKey
}

func NewAccountMapper(storeKey types.StoreKey) AccountMapper {
	return &accountMapper{
		storeKey: storeKey,
	}
}

func (am *accountMapper) GetBalance(ctx types.Context, addr types.AccAddress) (uint64, error) {
	return am.getBalance(am.getStore(ctx), addr)
}



func (am *accountMapper) getBalance(kvs types.KVStore, addr types.AccAddress) (uint64, error) {
	v := kvs.Get(addr.Bytes())
	if v == nil {
		return 0, ErrAccountNotFound
	}
	return util.BytesToUint64(v)
}

func (am *accountMapper) AddBalance(ctx types.Context, addr types.AccAddress, amount uint64) (uint64, error) {
	return am.addBalance(am.getStore(ctx), addr, amount)
}

func (am *accountMapper) addBalance(kvs types.KVStore, addr types.AccAddress, amount uint64) (uint64, error) {
	bl, err := am.getBalance(kvs, addr)
	if err != nil && err != ErrAccountNotFound {
		return 0, err
	}
	total := bl + amount
	return total, setBalance(kvs, addr, total)
}

func (am *accountMapper) SubBalance(ctx types.Context, addr types.AccAddress, amount uint64) (uint64, error) {
	return am.subBalance(am.getStore(ctx), addr, amount)
}

func (am *accountMapper) subBalance(kvs types.KVStore, addr types.AccAddress, amount uint64) (uint64, error) {
	bl, err := am.getBalance(kvs, addr)
	if err != nil && err != ErrAccountNotFound {
		return 0, err
	}
	if amount > bl {
		return 0, ErrNotEnoughBalance
	}
	total := bl - amount
	return total, setBalance(kvs, addr, total)
}

func (am *accountMapper) Transfer(ctx types.Context, from types.AccAddress, amount uint64, to types.AccAddress) error {
	kvs := am.getStore(ctx)
	if _, err := am.subBalance(kvs, from, amount); err != nil {
		return err
	}
	if _, err := am.addBalance(kvs, to, amount); err != nil {
		return err
	}
	return nil
}

func (am *accountMapper) getStore(ctx types.Context) types.KVStore {
	return ctx.KVStore(am.storeKey)
}

func setBalance(kvs types.KVStore, addr types.AccAddress, amount uint64) error {
	kvs.Set(addr.Bytes(), util.Uint64ToBytes(amount))
	return nil
}
