package subspace

import (
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	"github.com/ci123chain/ci123chain/pkg/abci/store"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	"reflect"
)

const (
	// StoreKey is the string store types for the param store
	StoreKey = "params"

	// TStoreKey is the string store types for the param transient store
	TStoreKey = "transient_params"
)

// Individual parameter store for each keeper
// Transient store persists for a block, so we use it for
// recording whether the parameter has been changed or not
type Subspace struct {
	cdc  *codec.Codec
	key  types.StoreKey // []byte -> []byte, stores parameter
	tkey types.StoreKey // []byte -> bool, stores parameter change

	name []byte

	table KeyTable
}

// NewSubspace constructs a store with namestore
func NewSubspace(cdc *codec.Codec, key types.StoreKey, tkey types.StoreKey, name string) (res Subspace) {
	res = Subspace{
		cdc:  cdc,
		key:  key,
		tkey: tkey,
		name: []byte(name),
		table: KeyTable{
			m: make(map[string]attribute),
		},
	}

	return
}

// HasKeyTable returns if the Subspace has a KeyTable registered.
func (s Subspace) HasKeyTable() bool {
	return len(s.table.m) > 0
}



// WithKeyTable initializes KeyTable and returns modified Subspace
func (s Subspace) WithKeyTable(table KeyTable) Subspace {
	if table.m == nil {
		panic("SetKeyTable() called with nil KeyTable")
	}
	if len(s.table.m) != 0 {
		panic("SetKeyTable() called on already initialized Subspace")
	}

	for k, v := range table.m {
		s.table.m[k] = v
	}

	// Allocate additional capicity for Subspace.name
	// So we don't have to allocate extra space each time appending to the key
	name := s.name
	s.name = make([]byte, len(name), len(name)+table.maxKeyLength())
	copy(s.name, name)

	return s
}


// Returns a KVStore identical with ctx.KVStore(s.types).Prefix()
func (s Subspace) kvStore(ctx types.Context) types.KVStore {
	// append here is safe, appends within a function won't cause
	// weird side effects when its singlethreaded
	return store.NewPrefixStore(ctx.KVStore(s.key), append(s.name, '/'))
}


// Returns a transient store for modification
func (s Subspace) transientStore(ctx types.Context) types.KVStore {
	// append here is safe, appends within a function won't cause
	// weird side effects when its singlethreaded
	return store.NewPrefixStore(ctx.TransientStore(s.tkey), append(s.name, '/'))
}



// Get parameter from store
func (s Subspace) Get(ctx types.Context, key []byte, ptr interface{}) {
	store := s.kvStore(ctx)
	bz := store.Get(key)
	err := s.cdc.UnmarshalJSON(bz, ptr)
	if err != nil {
		panic(err)
	}
}


// Set stores the parameter. It returns error if stored parameter has different type from input.
// It also set to the transient store to record change.
func (s Subspace) Set(ctx types.Context, key []byte, param interface{}) {
	store := s.kvStore(ctx)

	s.checkType(store, key, param)

	bz, err := s.cdc.MarshalJSON(param)
	if err != nil {
		panic(err)
	}
	store.Set(key, bz)

	tstore := s.transientStore(ctx)
	tstore.Set(key, []byte{})
}

func (s Subspace) checkType(store types.KVStore, key []byte, param interface{}) {
	attr, ok := s.table.m[string(key)]
	if !ok {
		panic("Parameter not registered")
	}

	ty := attr.ty
	pty := reflect.TypeOf(param)
	if pty.Kind() == reflect.Ptr {
		pty = pty.Elem()
	}

	if pty != ty {
		panic("Type mismatch with registered table")
	}
}


// Set from ParamSet
func (s Subspace) SetParamSet(ctx types.Context, ps ParamSet) {
	for _, pair := range ps.ParamSetPairs() {
		// pair.Field is a pointer to the field, so indirecting the ptr.
		// go-amino automatically handles it but just for sure,
		// since SetStruct is meant to be used in InitGenesis
		// so this method will not be called frequently
		v := reflect.Indirect(reflect.ValueOf(pair.Value)).Interface()
		s.Set(ctx, pair.Key, v)
	}
}

// Get to ParamSet
func (s Subspace) GetParamSet(ctx types.Context, ps ParamSet) {
	for _, pair := range ps.ParamSetPairs() {
		s.Get(ctx, pair.Key, pair.Value)
	}
}