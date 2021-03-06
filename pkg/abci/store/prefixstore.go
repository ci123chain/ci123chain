package store

import (
	"bytes"
	db "github.com/tendermint/tm-db"
	"io"
	"reflect"

	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
)

var _ KVStore = prefixStore{}

// prefixStore is similar with tendermint/tendermint/libs/db/prefix_db
// both gives access only to the limited subset of the store
// for convinience or safety

func NewPrefixStore(parent KVStore, prefix []byte) prefixStore {
	return prefixStore{
		parent: parent,
		prefix: prefix,
	}
}

type prefixStore struct {
	parent KVStore
	prefix []byte
}

func cloneAppend(bz []byte, tail []byte) (res []byte) {
	res = make([]byte, len(bz)+len(tail))
	copy(res, bz)
	copy(res[len(bz):], tail)
	return
}

func (s prefixStore) key(key []byte) (res []byte) {
	if key == nil {
		panic("nil key on prefixStore")
	}
	prefix := []byte("s/k:" + string(s.prefix) + "/")
	res = append(prefix, key...)
	return
}

// Implements Store
func (s prefixStore) GetStoreType() StoreType {
	return s.parent.GetStoreType()
}

// Implements CacheWrap
func (s prefixStore) CacheWrap() CacheWrap {
	return NewCacheKVStore(s)
}

// CacheWrapWithTrace implements the KVStore interface.
func (s prefixStore) CacheWrapWithTrace(w io.Writer, tc TraceContext) CacheWrap {
	return NewCacheKVStore(NewTraceKVStore(s, w, tc))
}

// Implements KVStore
func (s prefixStore) Get(key []byte) []byte {
	res := s.parent.Get(s.key(key))
	return res
}

// Implements KVStore
func (s prefixStore) Has(key []byte) bool {
	return s.parent.Has(s.key(key))
}

// Implements KVStore
func (s prefixStore) Set(key, value []byte) {
	s.parent.Set(s.key(key), value)
}

// Implements KVStore
func (s prefixStore) Delete(key []byte) {
	s.parent.Delete(s.key(key))
}

// Implements KVStore
func (s prefixStore) Prefix(prefix []byte) KVStore {
	return prefixStore{s, prefix}
}

// Implements KVStore
func (s prefixStore) Latest(keys []string) KVStore {
	return nil
}

// Implements KVStore
func (s prefixStore) Parent() KVStore {
	return s.parent
}

func (s prefixStore) Commit() CommitID {
	return s.parent.(CommitStore).Commit()
}

func (s prefixStore) LastCommitID() CommitID {
	return s.parent.(CommitStore).LastCommitID()
}

func (s prefixStore) SetPruning(strategy PruningStrategy) {
	s.parent.(CommitStore).SetPruning(strategy)
}

// Implements KVStore
func (s prefixStore) BatchSet(batch db.Batch) {
	s.parent.(*baseKVStore).BatchSet(batch)
}

// Implements KVStore
func (s prefixStore) GetCache() map[string][]byte{
	return s.parent.(*baseKVStore).GetCache()
}

// Implements KVStore
func (s prefixStore) Gas(meter GasMeter, config GasConfig) KVStore {
	return NewGasKVStore(meter, config, s)
}

// Implements KVStore.
func (s prefixStore) RemoteIterator(start, end []byte) Iterator {
	start = s.key(start)
	end = s.key(end)
	p := s.Parent()
	for {
		if reflect.TypeOf(p) != reflect.TypeOf(dbStoreAdapter{}) {
			p = p.Parent()
			if reflect.TypeOf(p) == reflect.TypeOf(prefixStore{}) {
				start = p.(prefixStore).key(start)
				end = p.(prefixStore).key(end)
			}
		} else {
			break
		}
	}
	return p.RemoteIterator(start, end)
}

// Implements KVStore
// Check https://github.com/tendermint/tendermint/blob/master/libs/db/prefix_db.go#L106
func (s prefixStore) Iterator(start, end []byte) Iterator {
	newstart := cloneAppend(s.prefix, start)

	var newend []byte
	if end == nil {
		newend = cpIncr(s.prefix)
	} else {
		newend = cloneAppend(s.prefix, end)
	}

	iter := s.parent.Iterator(newstart, newend)

	return newPrefixIterator(s.prefix, start, end, iter)
}

// Implements KVStore
// Check https://github.com/tendermint/tendermint/blob/master/libs/db/prefix_db.go#L129
func (s prefixStore) ReverseIterator(start, end []byte) Iterator {
	newstart := cloneAppend(s.prefix, start)

	var newend []byte
	if end == nil {
		newend = cpIncr(s.prefix)
	} else {
		newend = cloneAppend(s.prefix, end)
	}

	iter := s.parent.ReverseIterator(newstart, newend)

	return newPrefixIterator(s.prefix, start, end, iter)
}

var _ sdk.Iterator = (*prefixIterator)(nil)

type prefixIterator struct {
	prefix     []byte
	start, end []byte
	iter       Iterator
	valid      bool
}

func newPrefixIterator(prefix, start, end []byte, parent Iterator) *prefixIterator {
	return &prefixIterator{
		prefix: prefix,
		start:  start,
		end:    end,
		iter:   parent,
		valid:  parent.Valid() && bytes.HasPrefix(parent.Key(), prefix),
	}
}

// Implements Iterator
func (iter *prefixIterator) Domain() ([]byte, []byte) {
	return iter.start, iter.end
}

// Implements Iterator
func (iter *prefixIterator) Valid() bool {
	return iter.valid && iter.iter.Valid()
}

// Implements Iterator
func (iter *prefixIterator) Next() {
	if !iter.valid {
		panic("prefixIterator invalid, cannot call Next()")
	}
	iter.iter.Next()
	if !iter.iter.Valid() || !bytes.HasPrefix(iter.iter.Key(), iter.prefix) {
		iter.valid = false
	}
}

// Implements Iterator
func (iter *prefixIterator) Key() (key []byte) {
	if !iter.valid {
		panic("prefixIterator invalid, cannot call Key()")
	}
	key = iter.iter.Key()
	key = stripPrefix(key, iter.prefix)
	return
}

// Implements Iterator
func (iter *prefixIterator) Value() []byte {
	if !iter.valid {
		panic("prefixIterator invalid, cannot call Value()")
	}
	return iter.iter.Value()
}

// Implements Iterator
func (iter *prefixIterator) Close() {
	iter.iter.Close()
}

// copied from github.com/tendermint/tm-db/prefix_db.go
func stripPrefix(key []byte, prefix []byte) []byte {
	if len(key) < len(prefix) || !bytes.Equal(key[:len(prefix)], prefix) {
		panic("should not happen")
	}
	return key[len(prefix):]
}

// wrapping sdk.PrefixEndBytes
func cpIncr(bz []byte) []byte {
	return sdk.PrefixEndBytes(bz)
}

// copied from github.com/tendermint/tm-db/util.go
func cpDecr(bz []byte) (ret []byte) {
	if len(bz) == 0 {
		panic("cpDecr expects non-zero bz length")
	}
	ret = make([]byte, len(bz))
	copy(ret, bz)
	for i := len(bz) - 1; i >= 0; i-- {
		if ret[i] > byte(0x00) {
			ret[i]--
			return
		}
		ret[i] = byte(0xFF)
		if i == 0 {
			return nil
		}
	}
	return nil
}

func skipOne(iter Iterator, skipKey []byte) {
	if iter.Valid() {
		if bytes.Equal(iter.Key(), skipKey) {
			iter.Next()
		}
	}
}

