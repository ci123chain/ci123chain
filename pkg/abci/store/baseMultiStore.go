package store

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/app/types"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tm-db"
	"io"
)

type baseMultiStore struct {
	db           dbm.DB
	lastCommitID CommitID
	pruning      sdk.PruningStrategy
	storesParams map[StoreKey]storeParams
	stores       map[StoreKey]CommitStore
	keysByName   map[string]StoreKey
}

func NewBaseMultiStore(db dbm.DB) *baseMultiStore {
	return &baseMultiStore{
		db:           db,
		storesParams: make(map[StoreKey]storeParams),
		stores:       make(map[StoreKey]CommitStore),
		keysByName:   make(map[string]StoreKey),
	}
}

func(bs *baseMultiStore) SetPruning(pruning sdk.PruningStrategy) {
	bs.pruning = pruning
	for _, substore := range bs.stores {
		substore.SetPruning(pruning)
	}
}

func (bs *baseMultiStore) MountStoreWithDB(key StoreKey, typ StoreType, db dbm.DB) {
	if key == nil {
		panic("MountIAVLStore() types cannot be nil")
	}
	if _, ok := bs.storesParams[key]; ok {
		panic(fmt.Sprintf("rootMultiStore duplicate store types %v", key))
	}
	if _, ok := bs.keysByName[key.Name()]; ok {
		panic(fmt.Sprintf("rootMultiStore duplicate store types name %v", key))
	}
	bs.storesParams[key] = storeParams{
		key: key,
		typ: typ,
		db:  db,
	}
	bs.keysByName[key.Name()] = key
}

func (bs *baseMultiStore) GetStoreType() StoreType {
	return sdk.StoreTypeMulti
}

func (bs *baseMultiStore) GetCommitStore(key StoreKey) CommitStore {
	return nil
}

func (bs *baseMultiStore) GetCommitKVStore(key StoreKey) CommitKVStore {
	return nil
}

func (bs *baseMultiStore) GetStore(key StoreKey) Store {
	return nil
}

func (bs *baseMultiStore) GetKVStore(key StoreKey) KVStore {
	store := bs.stores[key].(*baseKVStore)

	return store
}

func (bs *baseMultiStore) LoadLatestVersion() error {

	ver := getLatestVersion(dbStoreAdapter{bs.db})
	err := bs.LoadVersion(ver)
	return err
}

func (bs *baseMultiStore) LoadVersion(ver int64) error {
	if ver == 0 {
		for key, storeParams := range bs.storesParams {
			id := CommitID{}
			err := bs.loadCommitStoreFromParams(key, id, storeParams)
			if err != nil {
				return fmt.Errorf("failed to load rootMultiStore: %v", err)
			}
		}
		bs.lastCommitID = CommitID{}
		return nil
	}

	for key, storeParams := range bs.storesParams {
		id := CommitID{}
		err := bs.loadCommitStoreFromParams(key, id, storeParams)
		if err != nil {
			return fmt.Errorf("failed to load rootMultiStore: %v", err)
		}
	}
	cInfo, err := getCommitInfo(dbStoreAdapter{bs.db}, ver)
	if err != nil {
		return err
	}
	bs.lastCommitID = cInfo.CommitID()
	return err
}

func (bs *baseMultiStore) LastCommitID() CommitID {
	return bs.lastCommitID
}

func (bs *baseMultiStore) Commit() CommitID {
	var CommitInfo commitInfo
	version := bs.lastCommitID.Version + 1
	cInfoKey := fmt.Sprintf(types.CommitInfoKeyFmt, version)
	cInfoBytes := bs.db.Get([]byte(cInfoKey))
	if cInfoBytes == nil {
		// Commit stores.
		CommitInfo = commitBaseStores(version, bs.stores)
		// Need to update atomically.

		//batch := bs.db.NewBatch()
		//setCommitInfo(batch, version, CommitInfo)
		//setLatestVersion(batch, version)
		//batch.Write()
		SetCommitInfo(bs.db, version, CommitInfo)
		SetLatestVersion(bs.db, version)
	}else{
		cdc.MustUnmarshalBinaryLengthPrefixed(cInfoBytes, &CommitInfo)
	}
	// Prepare for next version.
	commitID := CommitID{
		Version: version,
		Hash:    CommitInfo.Hash(),
	}
	bs.lastCommitID = commitID
	return commitID
}

func (bs *baseMultiStore) Write() {
	return
}

func commitBaseStores(version int64, storeMap map[StoreKey]CommitStore) commitInfo {
	storeInfos := make([]storeInfo, 0, len(storeMap))

	for key, store := range storeMap {
		// Commit
		commitID := store.Commit()

		if store.GetStoreType() == sdk.StoreTypeTransient {
			continue
		}

		// Record CommitID
		si := storeInfo{}
		si.Name = key.Name()
		si.Core.CommitID = commitID
		// si.Core.StoreType = store.GetStoreType()
		storeInfos = append(storeInfos, si)
	}

	ci := commitInfo{
		Version:    version,
		StoreInfos: storeInfos,
	}
	return ci
}

func (bs *baseMultiStore) WithTracer(w io.Writer) MultiStore {
	return nil
}

// WithTracingContext updates the tracing context for the MultiStore by merging
// the given context with the existing context by types. Any existing keys will
// be overwritten. It is implied that the caller should update the context when
// necessary between tracing operations. It returns a modified MultiStore.
func (bs *baseMultiStore) WithTracingContext(tc TraceContext) MultiStore {
	return nil
}

func (bs *baseMultiStore) TracingEnabled() bool {
	return false
}

// ResetTraceContext resets the current tracing context.
func (bs *baseMultiStore) ResetTraceContext() MultiStore {
	return nil
}

// Implements CacheWrapper/Store/CommitStore.
func (bs *baseMultiStore) CacheWrap() CacheWrap {
	return bs.CacheMultiStore().(CacheWrap)
}

// CacheWrapWithTrace implements the CacheWrapper interface.
func (bs *baseMultiStore) CacheWrapWithTrace(_ io.Writer, _ TraceContext) CacheWrap {
	return bs.CacheWrap()
}

//----------------------------------------
// +MultiStore

// Implements MultiStore.
func (bs *baseMultiStore) CacheMultiStore() CacheMultiStore {
	nbs := cacheMultiStore{
		db:           NewCacheKVStore(dbStoreAdapter{bs.db}),
		stores:       make(map[StoreKey]CacheWrap, len(bs.stores)),
		keysByName:   bs.keysByName,
	}

	for key, store := range bs.stores {
		nbs.stores[key] = store.CacheWrap()
	}
	return nbs
}

func (bs *baseMultiStore) loadCommitStoreFromParams(key sdk.StoreKey, id CommitID, params storeParams) error {
	_, ok := bs.stores[key]
	if !ok {
		store := NewBaseKVStore(dbStoreAdapter{bs.db}, int64(0), int64(0), key)
		store.SetPruning(bs.pruning)
		bs.stores[key] = store
	}

	return nil
}


//----------------------------------------
//query
func (bs *baseMultiStore) Query(req abci.RequestQuery) abci.ResponseQuery {
	// Query just routes this to a substore.
	path := req.Path
	storeName, subpath, err := parsePath(path)
	if err != nil {
		return err.QueryResult()
	}

	store := bs.getStoreByName(storeName)
	if store == nil {
		msg := fmt.Sprintf("no such store: %s", storeName)
		return sdk.ErrUnknownRequest(msg).QueryResult()
	}
	queryable, ok := store.(Queryable)
	if !ok {
		msg := fmt.Sprintf("store %s doesn't support queries", storeName)
		return sdk.ErrUnknownRequest(msg).QueryResult()
	}

	// trim the path and make the query
	req.Path = subpath
	res := queryable.Query(req)

	return res
}

func (bs *baseMultiStore) getStoreByName(name string) Store {
	key := bs.keysByName[name]
	if key == nil {
		return nil
	}
	return bs.stores[key]
}

func SetCommitInfo(db dbm.DB, version int64, info commitInfo) {
	infoByte, _ := cdc.MarshalBinaryLengthPrefixed(info)
	db.Set([]byte(fmt.Sprintf(types.CommitInfoKeyFmt, version)),infoByte)
}

func SetLatestVersion(db dbm.DB, version int64) {
	versionByte, _ := cdc.MarshalBinaryLengthPrefixed(version)
	db.Set([]byte(types.LatestVersionKey),versionByte)
}