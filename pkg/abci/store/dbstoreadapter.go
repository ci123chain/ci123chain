package store

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	dbm "github.com/tendermint/tm-db"
	"io"
)

// Wrapper type for dbm.Db with implementation of KVStore
type dbStoreAdapter struct {
	dbm.DB
}

func (dsa dbStoreAdapter) Get(key []byte) []byte {
	v, err := dsa.DB.Get(key)
	if err != nil {
		return nil
	}
	return v
}

func (dsa dbStoreAdapter) Has(key []byte) bool {
	ok, err := dsa.DB.Has(key)
	if err != nil || !ok {
		return false
	}
	return true
}

func (dsa dbStoreAdapter) Set(key, value []byte) {
	_ = dsa.DB.Set(key, value)
}

func (dsa dbStoreAdapter) Delete(key []byte) {
	_ = dsa.DB.Delete(key)
}

func (dsa dbStoreAdapter) Iterator(start, end []byte) sdk.Iterator {
	I, _ := dsa.DB.Iterator(start,end)
	return I
}

func (dsa dbStoreAdapter) ReverseIterator(start, end []byte) sdk.Iterator {
	I, _ := dsa.DB.ReverseIterator(start, end)
	return I
}

// Implements Store.
func (dbStoreAdapter) GetStoreType() StoreType {
	return sdk.StoreTypeDB
}

// Implements KVStore.
func (dsa dbStoreAdapter) CacheWrap() CacheWrap {
	return NewCacheKVStore(dsa)
}

// CacheWrapWithTrace implements the KVStore interface.
func (dsa dbStoreAdapter) CacheWrapWithTrace(w io.Writer, tc TraceContext) CacheWrap {
	return NewCacheKVStore(NewTraceKVStore(dsa, w, tc))
}

// Implements KVStore
func (dsa dbStoreAdapter) Prefix(prefix []byte) KVStore {
	return prefixStore{dsa, prefix}
}

// Implements KVStore
func (dsa dbStoreAdapter) Gas(meter GasMeter, config GasConfig) KVStore {
	return NewGasKVStore(meter, config, dsa)
}

// Implements KVStore
func (dsa dbStoreAdapter) Latest(keys []string) KVStore {
	return nil
}

// Implements KVStore
func (dsa dbStoreAdapter) Parent() KVStore {
	return nil
}

// dbm.DB implements KVStore so we can CacheKVStore it.
var _ KVStore = dbStoreAdapter{}

func (dsa dbStoreAdapter) RemoteIterator(start, end []byte) Iterator {
	i, _ := dsa.DB.Iterator(start, end)
	return i
}

//func ParentGetDbStoreAdapter(p KVStore, oriKey []byte) (db KVStore, key []byte){
//	for {
//		if reflect.TypeOf(p) != reflect.TypeOf(dbStoreAdapter{}) {
//			p = p.Parent()
//			if reflect.TypeOf(p) == reflect.TypeOf(prefixStore{}) {
//				pre := p.(prefixStore).prefix
//				oriKey = append(pre, oriKey...)
//			}
//		} else {
//			db = p
//			return
//		}
//	}
//}