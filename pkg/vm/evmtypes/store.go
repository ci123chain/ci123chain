package evmtypes

import (
	"bytes"
	"github.com/ci123chain/ci123chain/pkg/abci/store"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	"io"
)

type Store struct {
	parent types.KVStore
	prefix []byte
}

func (s Store) GetStoreType() types.StoreType {
	return s.parent.GetStoreType()
}

func (s Store) CacheWrap() types.CacheWrap {
	return store.NewCacheKVStore(s)
}

func (s Store) CacheWrapWithTrace(w io.Writer, tc types.TraceContext) types.CacheWrap {
	return store.NewCacheKVStore(store.NewTraceKVStore(s, w, tc))
}

func (s Store) Iterator(start, end []byte) types.Iterator {
	newstart := cloneAppend(s.prefix, start)

	var newend []byte
	if end == nil {
		newend = types.PrefixEndBytes(s.prefix)
	} else {
		newend = cloneAppend(s.prefix, end)
	}

	iter := s.parent.Iterator(newstart, newend)
	return NewStoreIterator(s.prefix, start, end, iter)
}

func (s Store) ReverseIterator(start, end []byte) types.Iterator {
	newstart := cloneAppend(s.prefix, start)

	var newend []byte
	if end == nil {
		newend = types.PrefixEndBytes(s.prefix)
	} else {
		newend = cloneAppend(s.prefix, end)
	}

	iter := s.parent.ReverseIterator(newstart, newend)

	return NewStoreIterator(s.prefix, start, end, iter)
}

func (s Store) RemoteIterator(start, end []byte) types.Iterator {
	start = s.key(start)
	if end == nil {
		end = types.PrefixEndBytes(start)
	} else {
		end = s.key(end)
	}

	return s.parent.RemoteIterator(start, end)
}

func (s Store) Prefix(prefix []byte) types.KVStore {
	return Store{parent:s, prefix:prefix}
}

func (s Store) Gas(gs types.GasMeter, gc types.GasConfig) types.KVStore {
	return store.NewGasKVStore(gs, gc, s)
}

func (s Store) Latest(Keys []string) types.KVStore {
	return nil
}

func (s Store) Parent() types.KVStore {
	return s.parent
}


var _ types.Iterator = (*StoreIterator)(nil)

type StoreIterator struct {
	prefix     []byte
	start, end []byte
	iter       types.Iterator
	valid      bool
}

func NewStoreIterator(prefix, start, end []byte, parent types.Iterator) *StoreIterator {
	v := parent.Valid()
	return &StoreIterator{
		prefix: prefix,
		start:  start,
		end:    end,
		iter:   parent,
		valid:  v,
	}
}

func (s StoreIterator) Domain() (start []byte, end []byte) {
	return s.start, s.end
}

func (s StoreIterator) Valid() bool {
	return s.valid && s.iter.Valid()
}

func (s StoreIterator) Next() {
	if !s.valid {
		panic("prefixIterator invalid, cannot call Next()")
	}
	s.iter.Next()
	if !s.iter.Valid() || !bytes.HasPrefix(s.iter.Key(), s.prefix) {
		s.valid = false
	}
}

func (s StoreIterator) Key() (key []byte) {
	if !s.valid {
		panic("prefixIterator invalid, cannot call Key()")
	}
	key = s.iter.Key()
	key = stripPrefix(key, s.prefix)
	return
}

func (s StoreIterator) Value() (value []byte) {
	if !s.valid {
		panic("prefixIterator invalid, cannot call Value()")
	}
	return s.iter.Value()
}

func (s StoreIterator) Error() error {
	return nil
}

func (s StoreIterator) Close() error {
	return s.iter.Close()
}

func stripPrefix(key []byte, prefix []byte) []byte {
	if len(key) < len(prefix) || !bytes.Equal(key[:len(prefix)], prefix) {
		panic("should not happen")
	}
	return key[len(prefix):]
}

func NewStore(parent types.KVStore, prefix []byte) Store {
	return Store{
		parent: parent,
		prefix: prefix,
	}
}

// ImplemenmiddleIns.fun["allocate"] = allocatets KVStore
func (s Store) Set(key, value []byte) {
	AssertValidKey(key)
	AssertValidValue(value)
	s.parent.Set(s.key(key), value)
}

// Implements KVStore
func (s Store) Get(key []byte) []byte {
	res := s.parent.Get(s.key(key))
	return res
}

// Implements KVStore
func (s Store) Has(key []byte) bool {
	res := s.parent.Get(s.key(key))
	return len(res) != 0
}

// Implements KVStore
func (s Store) Delete(key []byte) {
	s.parent.Delete(s.key(key))
}

func (s Store) key(key []byte) (res []byte) {
	if key == nil {
		panic("nil key on Store")
	}
	res = cloneAppend(s.prefix, key)
	return
}


func cloneAppend(bz []byte, tail []byte) (res []byte) {
	res = make([]byte, len(bz)+len(tail))
	copy(res, bz)
	copy(res[len(bz):], tail)
	return
}

// Check if the key is valid(key is not nil)
func AssertValidKey(key []byte) {
	if key == nil {
		panic("key is nil")
	}
}

// Check if the value is valid(value is not nil)
func AssertValidValue(value []byte) {
	if value == nil {
		panic("value is nil")
	}
}