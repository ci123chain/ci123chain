package store

import (
	"bytes"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"io"
	"reflect"
	"sort"
	"sync"
)

type latestStore struct {
	mtx    		sync.Mutex
	cache  		map[string]cValue
	parent 		KVStore
	storeEvery 	int64
	numRecent 	int64
	preKey		sdk.StoreKey
	latestKeys  []string
}

func NewlatestStore(parent KVStore, latestKeys []string) *latestStore {
	return &latestStore{
		cache:  	make(map[string]cValue),
		parent: 	parent,
		latestKeys: latestKeys,
	}
}

func (ks *latestStore) SetPruning(pruning sdk.PruningStrategy) {
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
func (ks *latestStore) GetStoreType() StoreType {
	return sdk.StoreTypeMulti
}

// Implements KVStore.
func (ks *latestStore) Get(key []byte) (value []byte) {
	ks.mtx.Lock()
	defer ks.mtx.Unlock()
	ks.assertValidKey(key)
	for _, v := range ks.latestKeys{
		if v == string(key) {
			value = ks.parentGet(key)
			return
		}
	}

	value = ks.parent.Get(key)
	return value
}

func (ks *latestStore) parentGet(key []byte) (value []byte) {
	oriKey := key
	p := ks.parent
	for {
		if reflect.TypeOf(p) != reflect.TypeOf(dbStoreAdapter{}) {
			p = p.Parent()
			if reflect.TypeOf(p) == reflect.TypeOf(prefixStore{}) {
				key = p.(prefixStore).key(key)
			}
		} else {
			break
		}
	}

	value = p.Get(key)
	if len(value) > 0 {
		return
	} else {
		value = ks.parent.Get(oriKey)
	}
	return
}

// Implements KVStore.
func (ks *latestStore) Set(key []byte, value []byte) {
	ks.mtx.Lock()
	defer ks.mtx.Unlock()
	ks.parent.Set(key, value)
}

// Implements KVStore.
func (ks *latestStore) Has(key []byte) bool {
	value := ks.Get(key)
	return value != nil
}

// Implements KVStore.
func (ks *latestStore) Delete(key []byte) {
	ks.mtx.Lock()
	defer ks.mtx.Unlock()
	ks.assertValidKey(key)
	ks.setCacheValue([]byte(key), nil, true, true)
}

// Implements KVStore
func (ks *latestStore) Prefix(prefix []byte) KVStore {
	return prefixStore{ks, prefix}
}

// Implements KVStore
func (ks *latestStore) Gas(meter GasMeter, config GasConfig) KVStore {
	return NewGasKVStore(meter, config, ks)
}

// Implements KVStore
func (ks *latestStore) Latest(keys []string) KVStore {
	return nil
}

// Implements KVStore
func (ks *latestStore) Parent() KVStore {
	return ks.parent
}

// Implements CacheWrapper.
func (ks *latestStore) CacheWrap() CacheWrap {
	return &cacheKVStore{
		cache:  make(map[string]cValue),
		parent: ks,
	}
}

// CacheWrapWithTrace implements the CacheWrapper interface.
func (ks *latestStore) CacheWrapWithTrace(w io.Writer, tc TraceContext) CacheWrap {
	return nil
}

// Implements KVStore.
func (ks *latestStore) RemoteIterator(start, end []byte) Iterator {
	p := ks.Parent()
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

// Implements KVStore.
func (ks *latestStore) Iterator(start, end []byte) Iterator {
	return ks.iterator(start, end, true)
}

// Implements KVStore.
func (ks *latestStore) ReverseIterator(start, end []byte) Iterator {
	return ks.iterator(start, end, false)
}

func (ks *latestStore) iterator(start, end []byte, ascending bool) Iterator {
	var parent, cache Iterator

	if ascending {
		parent = ks.parent.Iterator(start, end)
	} else {
		parent = ks.parent.ReverseIterator(start, end)
	}

	items := ks.dirtyItems(ascending)
	cache = newMemIterator(start, end, items)

	return newCacheMergeIterator(parent, cache, ascending)
}

// Constructs a slice of dirty items, to use w/ memIterator.
func (ks *latestStore) dirtyItems(ascending bool) []abci.EventAttribute {
	items := make([]abci.EventAttribute, 0, len(ks.cache))

	for key, cacheValue := range ks.cache {
		if !cacheValue.dirty {
			continue
		}

		items = append(items, abci.EventAttribute{Key: []byte(key), Value: cacheValue.value})
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
func (ks *latestStore) Write() {
	return
}

func (ks *latestStore) assertValidKey(key []byte) {
	if key == nil {
		panic("types is nil")
	}
}

func (ks *latestStore) assertValidValue(value []byte) {
	if value == nil {
		panic("value is nil")
	}
}

// Only entrypoint to mutate ci.cache.
func (ks *latestStore) setCacheValue(key, value []byte, deleted bool, dirty bool) {
	ks.cache[string(key)] = cValue{
		value:   value,
		deleted: deleted,
		dirty:   dirty,
	}
}

func (ks *latestStore) Commit() CommitID {
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

func (ks *latestStore) LastCommitID() CommitID {
	version := getLatestVersion(ks.parent)
	cInfo, err := getCommitInfo(ks.parent, version)
	if err != nil {
		panic(err)
	}
	return cInfo.CommitID()
}

//func (ks *latestStore) getCombineKey(key []byte) string {
//	//var version int64
//	//
//	//version = ks.LastCommitID().Version
//	//
//	//ckey := ks.preKey.Name() + "/" + strconv.FormatInt(version,10) + "/" + string(key)
//	//ckey := ks.preKey.Name() + "/" + "/" + string(key)
//	ckey := string(key)
//	return ckey
//}

func (ks *latestStore) setCombineKey(key []byte) string {
	//var version int64
	//
	//version = ks.LastCommitID().Version + 1
	//
	//ckey := ks.preKey.Name() + "/" + strconv.FormatInt(version,10) + "/" + string(key)
	//ckey := ks.preKey.Name() + "/" + "/" + string(key)
	ckey := string(key)
	return ckey
}

//-------------------------------------
//query
func (ks *latestStore) Query(req abci.RequestQuery) (res abci.ResponseQuery) {
	if len(req.Data) == 0 {
		msg := "Query cannot be zero length"
		return sdk.ErrTxDecode(msg).QueryResult()
	}

	// store the height we chose in the response, with 0 being changed to the
	// latest height

	switch req.Path {
	case "/key": // get by key
		key := req.Data // data holds the key bytes

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