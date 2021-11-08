package store

import (
	"encoding/json"
	"fmt"
	order "github.com/ci123chain/ci123chain/pkg/order/types"
	ics23 "github.com/confio/ics23/go"
	tmcrypto "github.com/tendermint/tendermint/proto/tendermint/crypto"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tm-db"

	sdkmaps "github.com/ci123chain/ci123chain/pkg/abci/store/internal/maps"
	sdkproofs "github.com/ci123chain/ci123chain/pkg/abci/store/internal/proofs"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"

	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
)

const (
	latestVersionKey = "s/latest"
	commitInfoKeyFmt = "s/%d" // s/<version>
	CacheName        = "cache"
)

// rootMultiStore is composed of many CommitStores. Name contrasts with
// cacheMultiStore which is for cache-wrapping other MultiStores. It implements
// the CommitMultiStore interface.
type rootMultiStore struct {
	ldb 	     dbm.DB
	cdb          dbm.DB
	lastCommitID CommitID
	pruning      sdk.PruningStrategy
	storesParams map[StoreKey]storeParams
	stores       map[StoreKey]CommitStore
	keysByName   map[string]StoreKey
	cacheDir     string

	traceWriter  io.Writer
	traceContext TraceContext
}

var _ CommitMultiStore = (*rootMultiStore)(nil)
var _ Queryable = (*rootMultiStore)(nil)

// nolint
func NewCommitMultiStore(ldb dbm.DB, cdb dbm.DB, cacheDir string) *rootMultiStore {
	return &rootMultiStore{
		ldb:          ldb,
		cdb:		  cdb,
		storesParams: make(map[StoreKey]storeParams),
		stores:       make(map[StoreKey]CommitStore),
		keysByName:   make(map[string]StoreKey),
		cacheDir:     cacheDir,
	}
}

// Implements CommitMultiStore
func (rs *rootMultiStore) SetPruning(pruning sdk.PruningStrategy) {
	rs.pruning = pruning
	for _, substore := range rs.stores {
		substore.SetPruning(pruning)
	}
}

// Implements Store.
func (rs *rootMultiStore) GetStoreType() StoreType {
	return sdk.StoreTypeMulti
}

// Implements CommitMultiStore.
func (rs *rootMultiStore) MountStoreWithDB(key StoreKey, typ StoreType, db dbm.DB) {
	if key == nil {
		panic("MountIAVLStore() key cannot be nil")
	}
	if _, ok := rs.storesParams[key]; ok {
		panic(fmt.Sprintf("rootMultiStore duplicate store key %v", key))
	}
	if _, ok := rs.keysByName[key.Name()]; ok {
		panic(fmt.Sprintf("rootMultiStore duplicate store key name %v", key))
	}
	rs.storesParams[key] = storeParams{
		key: key,
		typ: typ,
		db:  db,
	}
	rs.keysByName[key.Name()] = key
}

// Implements CommitMultiStore.
func (rs *rootMultiStore) GetCommitStore(key StoreKey) CommitStore {
	return rs.stores[key]
}

// Implements CommitMultiStore.
func (rs *rootMultiStore) GetCommitKVStore(key StoreKey) CommitKVStore {
	return rs.stores[key].(CommitKVStore)
}

func (rs *rootMultiStore) GetLatestVersion() int64 {
	return getLatestVersion(dbStoreAdapter{rs.ldb})
}

// Implements CommitMultiStore.
func (rs *rootMultiStore) LoadLatestVersion() error {
	ver := getLatestVersion(dbStoreAdapter{rs.ldb})
	return rs.LoadVersion(ver)
}

// Implements CommitMultiStore.
func (rs *rootMultiStore) LoadVersion(ver int64) error {

	// Special logic for version 0
	if ver == 0 {
		for key, storeParams := range rs.storesParams {
			id := CommitID{}
			store, err := rs.loadCommitStoreFromParams(key, id, storeParams)
			if err != nil {
				return fmt.Errorf("failed to load rootMultiStore: %v", err)
			}
			rs.stores[key] = store
		}

		rs.lastCommitID = CommitID{}
		return nil
	}
	// Otherwise, version is 1 or greater

	// Get commitInfo
	cInfo, err := getCommitInfo(dbStoreAdapter{rs.ldb}, ver)
	if err != nil {
		return err
	}

	// Convert StoreInfos slice to map
	infos := make(map[StoreKey]storeInfo)
	for _, storeInfo := range cInfo.StoreInfos {
		infos[rs.nameToKey(storeInfo.Name)] = storeInfo
	}

	// Load each Store
	var newStores = make(map[StoreKey]CommitStore)
	for key, storeParams := range rs.storesParams {
		var id CommitID
		info, ok := infos[key]
		if ok {
			id = info.Core.CommitID
		}

		store, err := rs.loadCommitStoreFromParams(key, id, storeParams)
		if err != nil {
			return fmt.Errorf("failed to load rootMultiStore: %v", err)
		}
		newStores[key] = store
	}

	// Success.
	rs.lastCommitID = cInfo.CommitID()
	rs.stores = newStores
	return nil
}

// WithTracer sets the tracer for the MultiStore that the underlying
// stores will utilize to trace operations. A MultiStore is returned.
func (rs *rootMultiStore) WithTracer(w io.Writer) MultiStore {
	rs.traceWriter = w
	return rs
}

func (rs *rootMultiStore) SetInitialVersion(version int64) error {
	rs.lastCommitID.Version = version - 1
	return nil
}

// WithTracingContext updates the tracing context for the MultiStore by merging
// the given context with the existing context by key. Any existing keys will
// be overwritten. It is implied that the caller should update the context when
// necessary between tracing operations. It returns a modified MultiStore.
func (rs *rootMultiStore) WithTracingContext(tc TraceContext) MultiStore {
	if rs.traceContext != nil {
		for k, v := range tc {
			rs.traceContext[k] = v
		}
	} else {
		rs.traceContext = tc
	}

	return rs
}

// TracingEnabled returns if tracing is enabled for the MultiStore.
func (rs *rootMultiStore) TracingEnabled() bool {
	return rs.traceWriter != nil
}

// ResetTraceContext resets the current tracing context.
func (rs *rootMultiStore) ResetTraceContext() MultiStore {
	rs.traceContext = nil
	return rs
}

//----------------------------------------
// +CommitStore

// Implements Committer/CommitStore.
func (rs *rootMultiStore) LastCommitID() CommitID {
	if rs.lastCommitID.Version == 0 {
		v := rs.GetLatestVersion()
		_ = rs.LoadVersion(v)
	}
	return rs.lastCommitID
}

// Implements Committer/CommitStore.
func (rs *rootMultiStore) Commit() CommitID {
	var commitInfo commitInfo
	var cacheMap []CacheMap

	version := rs.lastCommitID.Version + 1
	// check cache
	cacheName := filepath.Join(rs.cacheDir, CacheName)
	if _, err := os.Stat(cacheName); os.IsNotExist(err) {
		for _, store := range rs.stores {
			if reflect.TypeOf(store).Elem() == reflect.TypeOf(iavlStore{}){
				remote := store.(*iavlStore).Parent()
				if remote == nil {
					continue
				}
				s := remote.(*baseKVStore).GetCache()
				for k, v := range s{
					cacheMap = append(cacheMap, CacheMap{
						Key:   []byte(k),
						Value: v,
					})
				}
			}
		}
		cache, err := json.Marshal(cacheMap)
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile(cacheName, cache, os.ModePerm)
		if err != nil {
			panic(err)
		}
		//stores = rs.stores
	}
	// commitLocalStores
	commitInfo = rs.commitStores(version, rs.stores)

	// updateCommitInfo Need to update atomically.
	batch := rs.ldb.NewBatch()
	setCommitInfo(batch, version, commitInfo)
	setLatestVersion(batch, version)
	batch.Write()

	// commitSharedDB
	rs.commitSharedDB(cacheMap)

	//remove cache
	os.Remove(cacheName)

	//done
	if rs.cdb != nil {
		var orderBook order.OrderBook
		cdb := dbm.NewPrefixDB(rs.cdb, []byte("s/k:"+order.StoreKey+"/"))
		orderBytes, _ := cdb.Get([]byte(order.OrderBookKey))
		order.ModuleCdc.UnmarshalJSON(orderBytes, &orderBook)
		orderBook.Current.State = order.StateDone
		orderBytes, _ = order.ModuleCdc.MarshalJSON(orderBook)
		cdb.Set([]byte(order.OrderBookKey), orderBytes)
	}


	// Prepare for next version.
	commitID := CommitID{
		Version: version,
		Hash:    commitInfo.Hash(),
	}
	rs.lastCommitID = commitID
	return commitID
}

// Implements CacheWrapper/Store/CommitStore.
func (rs *rootMultiStore) CacheWrap() CacheWrap {
	return rs.CacheMultiStore().(CacheWrap)
}

// CacheWrapWithTrace implements the CacheWrapper interface.
func (rs *rootMultiStore) CacheWrapWithTrace(_ io.Writer, _ TraceContext) CacheWrap {
	return rs.CacheWrap()
}

//----------------------------------------
// +MultiStore

// Implements MultiStore.
func (rs *rootMultiStore) CacheMultiStore() CacheMultiStore {
	return newCacheMultiStoreFromRMS(rs)
}

// Implements MultiStore.
func (rs *rootMultiStore) GetStore(key StoreKey) Store {
	return rs.stores[key]
}

// GetKVStore implements the MultiStore interface. If tracing is enabled on the
// rootMultiStore, a wrapped TraceKVStore will be returned with the given
// tracer, otherwise, the original KVStore will be returned.
func (rs *rootMultiStore) GetKVStore(key StoreKey) KVStore {
	store := rs.stores[key].(KVStore)
	if rs.TracingEnabled() {
		store = NewTraceKVStore(store, rs.traceWriter, rs.traceContext)
	}

	return store
}

// Implements MultiStore

// getStoreByName will first convert the original name to
// a special key, before looking up the CommitStore.
// This is not exposed to the extensions (which will need the
// StoreKey), but is useful in main, and particularly app.Query,
// in order to convert human strings into CommitStores.
func (rs *rootMultiStore) getStoreByName(name string) Store {
	key := rs.keysByName[name]
	if key == nil {
		return nil
	}
	return rs.stores[key]
}

//---------------------- Query ------------------

// Query calls substore.Query with the same `req` where `req.Path` is
// modified to remove the substore prefix.
// Ie. `req.Path` here is `/<substore>/<path>`, and trimmed to `/<path>` for the substore.
// TODO: add proof for `multistore -> substore`.
func (rs *rootMultiStore) Query(req abci.RequestQuery) abci.ResponseQuery {
	// Query just routes this to a substore.
	path := req.Path
	storeName, subpath, err := parsePath(path)
	if err != nil {
		return sdkerrors.QueryResult(err)
	}

	store := rs.getStoreByName(storeName)
	if store == nil {
		msg := fmt.Sprintf("no such store: %s", storeName)
		return sdkerrors.QueryResult(sdkerrors.Wrap(sdkerrors.ErrParams, msg))
	}
	queryable, ok := store.(Queryable)
	if !ok {
		msg := fmt.Sprintf("store %s doesn't support queries", storeName)
		return sdkerrors.QueryResult(sdkerrors.Wrap(sdkerrors.ErrParams, msg))
	}

	// trim the path and make the query
	req.Path = subpath
	res := queryable.Query(req)

	if !req.Prove || !RequireProof(subpath) {
		return res
	}

	if res.ProofOps == nil || len(res.ProofOps.Ops) == 0 {
		return sdkerrors.QueryResult(sdkerrors.Wrap(sdkerrors.ErrResponse, "substore proof was nil/empty when it should never be"))
	}

	commitInfo, errMsg := getCommitInfo(dbStoreAdapter{rs.ldb}, res.Height)
	if errMsg != nil {
		return sdkerrors.QueryResult(errMsg)
	}

	// Restore origin path and append proof op.
	//res.ProofOps.Ops = append(res.ProofOps.Ops, NewMultiStoreProofOp(
	//	[]byte(storeName),
	//	NewMultiStoreProof(commitInfo.StoreInfos),
	//).ProofOp())

	res.ProofOps.Ops = append(res.ProofOps.Ops, commitInfo.ProofOp(storeName))
	return res
}

// parsePath expects a format like /<storeName>[/<subpath>]
// Must start with /, subpath may be empty
// Returns error if it doesn't start with /
func parsePath(path string) (storeName string, subpath string, err error) {
	if !strings.HasPrefix(path, "/") {
		err = sdkerrors.Wrap(sdkerrors.ErrParams, fmt.Sprintf("invalid path: %v", path))
		return
	}

	paths := strings.SplitN(path[1:], "/", 2)
	storeName = paths[0]

	if len(paths) == 2 {
		subpath = "/" + paths[1]
	}

	return
}

//----------------------------------------

func (rs *rootMultiStore) loadCommitStoreFromParams(key sdk.StoreKey, id CommitID, params storeParams) (store CommitStore, err error) {
	var ldb dbm.DB
	var cdb dbm.DB
	if params.db != nil {
		ldb = dbm.NewPrefixDB(params.db, []byte("s/_/"))
		cdb = rs.cdb
	} else {
		ldb = dbm.NewPrefixDB(rs.ldb, []byte("s/k:"+params.key.Name()+"/"))
		if rs.cdb != nil {
			cdb = dbm.NewPrefixDB(rs.cdb, []byte("s/k:"+params.key.Name()+"/"))
		}
	}

	switch params.typ {
	case sdk.StoreTypeMulti:
		panic("recursive MultiStores not yet supported")
		// TODO: id?
		// return NewCommitMultiStore(db, id)
	case sdk.StoreTypeIAVL:
		store, err = LoadIAVLStore(ldb, cdb, id, rs.pruning, key)
		return
	case sdk.StoreTypeDB:
		panic("dbm.DB is not a CommitStore")
	case sdk.StoreTypeTransient:
		_, ok := key.(*sdk.TransientStoreKey)
		if !ok {
			err = fmt.Errorf("invalid StoreKey for StoreTypeTransient: %s", key.String())
			return
		}
		store = newTransientStore()
		return
	case sdk.StoreTypeMemory:
		if _, ok := key.(*sdk.MemoryStoreKey); !ok {
			return nil, fmt.Errorf("unexpected key type for a MemoryStoreKey; got: %s", key.String())
		}
		store = NewMemoryStore()
		return
	default:
		panic(fmt.Sprintf("unrecognized store types %v", params.typ))
	}
}

func (rs *rootMultiStore) nameToKey(name string) StoreKey {
	for key := range rs.storesParams {
		if key.Name() == name {
			return key
		}
	}
	panic("Unknown name " + name)
}

//----------------------------------------
// storeParams

type storeParams struct {
	key StoreKey
	db  dbm.DB
	typ StoreType
}

//----------------------------------------
// commitInfo

// NOTE: Keep commitInfo a simple immutable struct.
type commitInfo struct {

	// Version
	Version int64

	// Store info for
	StoreInfos []storeInfo
}

// Hash returns the simple merkle root hash of the stores sorted by name.
func (ci commitInfo) Hash() []byte {
	if len(ci.StoreInfos) == 0 {
		return nil
	}

	rootHash, _, _ := sdkmaps.ProofsFromMap(ci.toMap())
	return rootHash
}

func (ci commitInfo) CommitID() CommitID {
	return CommitID{
		Version: ci.Version,
		Hash:    ci.Hash(),
	}
}


func (ci commitInfo) toMap() map[string][]byte {
	m := make(map[string][]byte, len(ci.StoreInfos))
	for _, storeInfo := range ci.StoreInfos {
		m[storeInfo.Name] = storeInfo.GetHash()
	}

	return m
}

func (ci commitInfo) ProofOp(storeName string) tmcrypto.ProofOp {
	cmap := ci.toMap()
	_, proofs, _ := sdkmaps.ProofsFromMap(cmap)

	proof := proofs[storeName]
	if proof == nil {
		panic(fmt.Sprintf("ProofOp for %s but not registered store name", storeName))
	}

	// convert merkle.SimpleProof to CommitmentProof
	existProof, err := sdkproofs.ConvertExistenceProof(proof, []byte(storeName), cmap[storeName])
	if err != nil {
		panic(fmt.Errorf("could not convert simple proof to existence proof: %w", err))
	}

	commitmentProof := &ics23.CommitmentProof{
		Proof: &ics23.CommitmentProof_Exist{
			Exist: existProof,
		},
	}

	return NewSimpleMerkleCommitmentOp([]byte(storeName), commitmentProof).ProofOp()
}

func NewSimpleMerkleCommitmentOp(key []byte, proof *ics23.CommitmentProof) CommitmentOp {
	return CommitmentOp{
		Type:  ProofOpSimpleMerkleCommitment,
		Spec:  ics23.TendermintSpec,
		Key:   key,
		Proof: proof,
	}
}
//----------------------------------------
// storeInfo

// storeInfo contains the name and collactor reference for an
// underlying store.  It is the leaf of the rootMultiStores top
// level simple merkle tree.
type storeInfo struct {
	Name string
	Core storeCore
}

func (si storeInfo) GetHash() []byte {
	return si.Core.CommitID.Hash
}

type storeCore struct {
	// StoreType StoreType
	CommitID CommitID
	// ... maybe add more state
}

// Implements merkle.Hasher.
//func (si storeInfo) Hash() []byte {
//	// Doesn't write Name, since merkle.SimpleHashFromMap() will
//	// include them via the keys.
//	bz, _ := cdc.MarshalBinaryLengthPrefixed(si.Core)
//	hasher := tmhash.New()
//
//	_, err := hasher.Write(bz)
//	if err != nil {
//		// TODO: Handle with #870
//		panic(err)
//	}
//
//	return hasher.Sum(nil)
//}

//----------------------------------------
// Misc.

func getLatestVersion(db KVStore) int64 {
	var latest int64
	latestBytes := db.Get([]byte(latestVersionKey))
	if latestBytes == nil {
		return 0
	}

	err := cdc.UnmarshalBinaryLengthPrefixed(latestBytes, &latest)
	if err != nil {
		panic(err)
	}

	return latest
}

// Set the latest version.
func setLatestVersion(batch dbm.Batch, version int64) {
	latestBytes, _ := cdc.MarshalBinaryLengthPrefixed(version)
	batch.Set([]byte(latestVersionKey), latestBytes)
}

// Commits each store and returns a new commitInfo.
func (rs *rootMultiStore) commitStores(version int64, storeMap map[StoreKey]CommitStore) commitInfo {
	storeInfos := make([]storeInfo, 0, len(storeMap))
	for key, store := range storeMap {
		// Commit localDB
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

func (rs *rootMultiStore) commitSharedDB(cacheMap []CacheMap) {
	if rs.cdb == nil {
		return
	}
	batch := rs.cdb.NewBatch()
	defer batch.Close()
	for _, cache := range cacheMap {
		batch.Set(cache.Key, cache.Value)
	}
	batch.Write()
}

// Gets commitInfo from disk.
func getCommitInfo(db KVStore, ver int64) (commitInfo, error) {

	// Get from DB.
	cInfoKey := fmt.Sprintf(commitInfoKeyFmt, ver)
	cInfoBytes := db.Get([]byte(cInfoKey))
	if cInfoBytes == nil {
		return commitInfo{}, fmt.Errorf("failed to get rootMultiStore: no data")
	}

	var cInfo commitInfo

	err := cdc.UnmarshalBinaryLengthPrefixed(cInfoBytes, &cInfo)
	if err != nil {
		return commitInfo{}, fmt.Errorf("failed to get rootMultiStore: %v", err)
	}

	return cInfo, nil
}

// Set a commitInfo for given version.
func setCommitInfo(batch dbm.Batch, version int64, cInfo commitInfo) {
	cInfoBytes := cdc.MustMarshalBinaryLengthPrefixed(cInfo)
	cInfoKey := fmt.Sprintf(commitInfoKeyFmt, version)
	batch.Set([]byte(cInfoKey), cInfoBytes)
}

type CacheMap struct {
	Key		 []byte `json:"key"`
	Value    []byte `json:"value"`
}