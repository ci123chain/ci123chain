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
	return dsa.DB.Iterator(start, end)
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