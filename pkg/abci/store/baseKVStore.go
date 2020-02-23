package store

import (
	"bytes"
	"fmt"
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	cmn "github.com/tendermint/tendermint/libs/common"
	"io"
	"sort"
	"sync"
)

type baseKVStore struct {
	mtx    		sync.Mutex
	cache  		map[string]cValue
	parent 		KVStore
	storeEvery 	int64
	numRecent 	int64
	preKey		sdk.StoreKey
}

func NewBaseKVStore(parent KVStore, storeEvery, numRecent int64, key sdk.StoreKey) *baseKVStore {
	return &baseKVStore{
		cache:  	make(map[string]cValue),
		parent: 	parent,
		storeEvery: storeEvery,
		numRecent:	numRecent,
		preKey:		key,
	}
}

func (ks *baseKVStore) SetPruning(pruning sdk.PruningStrategy) {
	switch pruning {
	case sdk.PruneEverything:
		ks.numRecent = 0
		ks.storeEvery = 0
	case sdk.PruneNothing:
		ks.storeEvery = 1
	case sdk.PruneSyncable:
		ks.numRecent = 100
		ks.storeEvery = 10000
	}
}

// Implements Store.
func (ks *baseKVStore) GetStoreType() StoreType {
	return sdk.StoreTypeMulti
}

// Implements KVStore.
func (ks *baseKVStore) Get(key []byte) (value []byte) {
	ks.mtx.Lock()
	defer ks.mtx.Unlock()
	ks.assertValidKey(key)
	ckey := ks.getCombineKey(key)
	cacheValue, ok := ks.cache[ckey]
	if !ok {
		value = ks.parent.Get([]byte(ckey))
			ks.setCacheValue([]byte(ckey), value, false, false)
	} else {
		value = cacheValue.value
	}

	return value
}

// Implements KVStore.
func (ks *baseKVStore) Set(key []byte, value []byte) {
	ks.mtx.Lock()
	defer ks.mtx.Unlock()
	ks.assertValidKey(key)
	ks.assertValidValue(value)
	ckey := ks.setCombineKey(key)
	ks.setCacheValue([]byte(ckey), value, false, true)
}

// Implements KVStore.
func (ks *baseKVStore) Has(key []byte) bool {
	value := ks.Get(key)
	return value != nil
}

// Implements KVStore.
func (ks *baseKVStore) Delete(key []byte) {
	ks.mtx.Lock()
	defer ks.mtx.Unlock()
	ks.assertValidKey(key)
	ckey := ks.getCombineKey(key)
	ks.setCacheValue([]byte(ckey), nil, true, true)
}

// Implements KVStore
func (ks *baseKVStore) Prefix(prefix []byte) KVStore {
	return prefixStore{ks, prefix}
}

// Implements KVStore
func (ks *baseKVStore) Gas(meter GasMeter, config GasConfig) KVStore {
	return NewGasKVStore(meter, config, ks)
}

// Implements CacheWrapper.
func (ks *baseKVStore) CacheWrap() CacheWrap {
	return &cacheKVStore{
		cache:  make(map[string]cValue),
		parent: ks,
	}
}

// CacheWrapWithTrace implements the CacheWrapper interface.
func (ks *baseKVStore) CacheWrapWithTrace(w io.Writer, tc TraceContext) CacheWrap {
	return nil
}

// Implements KVStore.
func (ks *baseKVStore) Iterator(start, end []byte) Iterator {
	//cstart := ks.getCombineKey(start)
	//cend := ks.getCombineKey(end)
	return ks.iterator([]byte(start), []byte(end), true)
}

// Implements KVStore.
func (ks *baseKVStore) ReverseIterator(start, end []byte) Iterator {
	cstart := ks.getCombineKey(start)
	cend := ks.getCombineKey(end)
	return ks.iterator([]byte(cstart), []byte(cend), false)
}

func (ks *baseKVStore) iterator(start, end []byte, ascending bool) Iterator {
	var parent, cache Iterator
	cstart := ks.getCombineKey(start)
	cend := ks.getCombineKey(end)

	if ascending {
		parent = ks.parent.Iterator([]byte(cstart), []byte(cend))
	} else {
		parent = ks.parent.ReverseIterator([]byte(cstart), []byte(cend))
	}

	items := ks.dirtyItems(ascending)
	cache = newMemIterator([]byte(cstart), []byte(cend), items)

	return newCacheMergeIterator(parent, cache, ascending)
}

// Constructs a slice of dirty items, to use w/ memIterator.
func (ks *baseKVStore) dirtyItems(ascending bool) []cmn.KVPair {
	items := make([]cmn.KVPair, 0, len(ks.cache))

	for key, cacheValue := range ks.cache {
		if !cacheValue.dirty {
			continue
		}

		items = append(items, cmn.KVPair{Key: []byte(key), Value: cacheValue.value})
	}

	sort.Slice(items, func(i, j int) bool {
		if ascending {
			return bytes.Compare(items[i].Key, items[j].Key) < 0
		}

		return bytes.Compare(items[i].Key, items[j].Key) > 0
	})

	return items
}

// Implements CacheKVStore.
func (ks *baseKVStore) Write() {
	return
}

func (ks *baseKVStore) assertValidKey(key []byte) {
	if key == nil {
		panic("types is nil")
	}
}

func (ks *baseKVStore) assertValidValue(value []byte) {
	if value == nil {
		panic("value is nil")
	}
}

// Only entrypoint to mutate ci.cache.
func (ks *baseKVStore) setCacheValue(key, value []byte, deleted bool, dirty bool) {
	ks.cache[string(key)] = cValue{
		value:   value,
		deleted: deleted,
		dirty:   dirty,
	}
}

func (ks *baseKVStore) Commit() CommitID {
	ks.mtx.Lock()
	defer ks.mtx.Unlock()

	// We need a copy of all of the keys.
	// Not the best, but probably not a bottleneck depending.
	keys := make([]string, 0, len(ks.cache))
	for key, dbValue := range ks.cache {
		if dbValue.dirty {
			keys = append(keys, key)
		}
	}

	sort.Strings(keys)

	var valueBytes []cValue
	// TODO: Consider allowing usage of Batch, which would allow the write to
	// at least happen atomically.
	for _, key := range keys {
		cacheValue := ks.cache[key]
		if cacheValue.deleted {
			ks.parent.Delete([]byte(key))
		} else if cacheValue.value == nil {
			// Skip, it already doesn't exist in parent.
		} else {
			ks.parent.Set([]byte(key), cacheValue.value)
			valueBytes = append(valueBytes, cacheValue)
		}
	}

	// compute commit hash
	bz, _ := cdc.MarshalBinaryLengthPrefixed(valueBytes)
	hasher := tmhash.New()

	_, err := hasher.Write(bz)
	if err != nil {
		panic(err)
	}
	hash := hasher.Sum(nil)

	// Clear the cache
	ks.cache = make(map[string]cValue)

	version := getLatestVersion(ks.parent)
	if version == 0 {
		return CommitID{
			Version: 1,
			Hash:	hash,
		}
	}else{
		cInfo, err := getCommitInfo(ks.parent, version)
		if err != nil {
			panic(err)
		}

		return CommitID{
			Version: cInfo.Version + 1,
			Hash:    hash,
		}
	}
}

func (ks *baseKVStore) LastCommitID() CommitID {
	version := getLatestVersion(ks.parent)
	cInfo, err := getCommitInfo(ks.parent, version)
	if err != nil {
		panic(err)
	}
	return cInfo.CommitID()
}

func (ks *baseKVStore) Latest(keys []string) KVStore {
	return NewlatestStore(ks, ks.storeEvery, ks.numRecent, ks.preKey, keys)
}

func (ks *baseKVStore) Parent() KVStore {
	return ks.parent
}

func (ks *baseKVStore) getCombineKey(key []byte) string {
	//var version int64
	//
	//version = ks.LastCommitID().Version
	//
	//ckey := ks.preKey.Name() + "/" + strconv.FormatInt(version,10) + "/" + string(key)
	ckey := ks.preKey.Name() + "/" + "/" + string(key)
	return ckey
}

func (ks *baseKVStore) setCombineKey(key []byte) string {
	//var version int64
    //
	//version = ks.LastCommitID().Version + 1
	//
	//ckey := ks.preKey.Name() + "/" + strconv.FormatInt(version,10) + "/" + string(key)
	ckey := ks.preKey.Name() + "/" + "/" + string(key)
	return ckey
}

//-------------------------------------
//query
func (ks *baseKVStore) Query(req abci.RequestQuery) (res abci.ResponseQuery) {
	if len(req.Data) == 0 {
		msg := "Query cannot be zero length"
		return sdk.ErrTxDecode(msg).QueryResult()
	}

	// store the height we chose in the response, with 0 being changed to the
	// latest height

	switch req.Path {
	case "/types": // get by types
		key := req.Data // data holds the types bytes

		res.Key = key
		value := ks.Get(key)
		res.Value = value
	case "/subspace":
		var KVs []KVPair

		subspace := req.Data
		res.Key = subspace
		iterator := sdk.KVStorePrefixIterator(ks, subspace)
		for ; iterator.Valid(); iterator.Next() {
			KVs = append(KVs, KVPair{Key: iterator.Key(), Value: iterator.Value()})
		}

		iterator.Close()
		res.Value = cdc.MustMarshalBinaryLengthPrefixed(KVs)

	default:
		msg := fmt.Sprintf("Unexpected Query path: %v", req.Path)
		return sdk.ErrUnknownRequest(msg).QueryResult()
	}

	return
}