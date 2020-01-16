package store

import (
	"io"

	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	dbm "github.com/tendermint/tm-db"
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
