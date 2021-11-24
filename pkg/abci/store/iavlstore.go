package store

import (
	"fmt"
	ics23 "github.com/confio/ics23/go"
	"github.com/cosmos/iavl"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	tmcrypto "github.com/tendermint/tendermint/proto/tendermint/crypto"
	"io"
	"reflect"
	"sync"

	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/logger"
	abci "github.com/tendermint/tendermint/abci/types"

	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	dbm "github.com/tendermint/tm-db"
)

const (
	FlagRunMode = "mode"
	ModeSingle = "single"
	ModeMulti = "multi"
	defaultIAVLCacheSize = 10000
)

// load the iavl store
func LoadIAVLStore(ldb, cdb dbm.DB, id CommitID, pruning sdk.PruningStrategy, key sdk.StoreKey) (CommitStore, error) {
	tree, err := iavl.NewMutableTree(ldb, defaultIAVLCacheSize)
	if err != nil {
		return nil, err
	}
	_, err = tree.LoadVersion(id.Version)
	if err != nil {
		return nil, err
	}
	iavl := newIAVLStore(cdb, tree, int64(0), int64(0), key)
	iavl.SetPruning(pruning)
	return iavl, nil
}

//----------------------------------------

var _ KVStore = (*IavlStore)(nil)
var _ CommitStore = (*IavlStore)(nil)
var _ Queryable = (*IavlStore)(nil)

// iavlStore Implements KVStore and CommitStore.
type IavlStore struct {

	// The underlying tree.
	tree *iavl.MutableTree

	// How many old versions we hold onto.
	// A value of 0 means keep no recent states.
	numRecent int64

	// This is the distance between state-sync waypoint states to be stored.
	// See https://github.com/tendermint/tendermint/issues/828
	// A value of 1 means store every state.
	// A value of 0 means store no waypoints. (node cannot assist in state-sync)
	// By default this value should be set the same across all nodes,
	// so that nodes can know the waypoints their peers store.
	storeEvery int64

	// KVStore save to shared DB
	parent CommitStore

	key   sdk.StoreKey
	lg logger.Logger

	mode string
}

// CONTRACT: tree should be fully loaded.
// nolint: unparam
func newIAVLStore(remoteDB dbm.DB, tree *iavl.MutableTree, numRecent int64, storeEvery int64, key sdk.StoreKey) *IavlStore {
	logger := logger.GetLogger()
	st := &IavlStore{
		tree:       tree,
		numRecent:  numRecent,
		storeEvery: storeEvery,
		lg:         logger,
		key:        key,
		mode: 		viper.GetString(FlagRunMode),
	}
	if remoteDB != nil {
		st.parent = NewBaseKVStore(dbStoreAdapter{remoteDB}, storeEvery, numRecent, key)
	}
	return st
}

func (st *IavlStore) localMode() bool {
	return st.mode == ModeSingle || st.mode == ""
}


// Implements Committer.
func (st *IavlStore) Commit() CommitID {
	// Save a new version.
	//st.parent.Commit()
	hash, version, err := st.tree.SaveVersion()
	if err != nil {
		// TODO: Do we want to extend Commit to allow returning errors?
		panic(err)
	}

	// Release an old version of history, if not a sync waypoint.
	previous := version - 1
	if st.numRecent < previous {
		toRelease := previous - st.numRecent
		if st.storeEvery == 0 || toRelease%st.storeEvery != 0 {
			err := st.tree.DeleteVersion(toRelease)
			if errCause := errors.Cause(err); errCause != nil && errCause != iavl.ErrVersionDoesNotExist {
				panic(err)
			}
		}
	}

	return CommitID{
		Version: version,
		Hash:    hash,
	}
}

// Implements Committer.
func (st *IavlStore) LastCommitID() CommitID {
	return CommitID{
		Version: st.tree.Version(),
		Hash:    st.tree.Hash(),
	}
}

// Implements Committer.
func (st *IavlStore) SetPruning(pruning sdk.PruningStrategy) {
	switch pruning {
	case sdk.PruneEverything:
		st.numRecent = 0
		st.storeEvery = 0
	case sdk.PruneNothing:
		st.storeEvery = 1
	case sdk.PruneSyncable:
		st.numRecent = 100
		st.storeEvery = 10000
	}
}

func (st *IavlStore) DeleteVersions(from int64, to int64) error {
	err := st.tree.DeleteVersionsRange(from, to)
	return err
}

// VersionExists returns whether or not a given version is stored.
func (st *IavlStore) VersionExists(version int64) bool {
	return st.tree.VersionExists(version)
}

// Implements Store.
func (st *IavlStore) GetStoreType() StoreType {
	return sdk.StoreTypeIAVL
}

// Implements Store.
func (st *IavlStore) CacheWrap() CacheWrap {
	return NewCacheKVStore(st)
}

// CacheWrapWithTrace implements the Store interface.
func (st *IavlStore) CacheWrapWithTrace(w io.Writer, tc TraceContext) CacheWrap {
	return NewCacheKVStore(NewTraceKVStore(st, w, tc))
}

// Implements KVStore.
func (st *IavlStore) Set(key, value []byte) {
	if st.parent != nil {
		st.parent.(KVStore).Set(key, value)
	}
	st.tree.Set(key, value)
}

// Implements KVStore.
func (st *IavlStore) Get(key []byte) []byte {
	_, localValue := st.tree.Get(key)

	if !st.localMode() {
		remoteValue := st.parent.(KVStore).Get(key)
		return remoteValue
	}

	return localValue
}


// Implements KVStore.
func (st *IavlStore) Has(key []byte) (exists bool) {
	if st.parent != nil {
		return st.parent.(KVStore).Has(key)
	}
	return false
}

// Implements KVStore.
func (st *IavlStore) Delete(key []byte) {
	if st.parent != nil {
		st.parent.(KVStore).Delete(key)
	}
	st.tree.Remove(key)
}

// Implements KVStore
func (st *IavlStore) Prefix(prefix []byte) KVStore {
	return prefixStore{st, prefix}
}

// Implements KVStore
func (st *IavlStore) Gas(meter GasMeter, config GasConfig) KVStore {
	return NewGasKVStore(meter, config, st)
}

// Implements KVStore
func (st *IavlStore) Latest(keys []string) KVStore {
	return NewlatestStore(st, keys)
}

// Implements KVStore
func (st *IavlStore) Parent() KVStore {
	if st.parent == nil {
		return nil
	}
	return st.parent.(KVStore)
}

// Implements KVStore.
func (st *IavlStore) RemoteIterator(start, end []byte) Iterator {
	p := st.Parent()
	if st.localMode() || p == nil {
		return st.Iterator(start, end)
	}
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
func (st *IavlStore) Iterator(start, end []byte) Iterator {
	return newIAVLIterator(st.tree.ImmutableTree, start, end, true)
}

// Implements KVStore.
func (st *IavlStore) ReverseIterator(start, end []byte) Iterator {
	return newIAVLIterator(st.tree.ImmutableTree, start, end, false)
}

// Handle gatest the latest height, if height is 0
func getHeight(tree *iavl.MutableTree, req abci.RequestQuery) int64 {
	height := req.Height
	if height == 0 {
		latest := tree.Version()
		//if tree.VersionExists(latest - 1) {
		//	height = latest - 1
		//} else {
		//	height = latest
		//}
		height = latest
	}
	return height
}

// Query implements ABCI interface, allows queries
//
// by default we will return from (latest height -1),
// as we will have merkle proofs immediately (header height = data height + 1)
// If latest-1 is not present, use latest (which must be present)
// if you care to have the latest data to see a tx results, you must
// explicitly set the height you want to see
func (st *IavlStore) Query(req abci.RequestQuery) (res abci.ResponseQuery) {
	if len(req.Data) == 0 {
		msg := "Query cannot be zero length"
		return sdkerrors.QueryResult(sdkerrors.Wrap(sdkerrors.ErrParams, msg))
	}

	tree := st.tree

	// store the height we chose in the response, with 0 being changed to the
	// latest height
	res.Height = getHeight(tree, req)

	switch req.Path {
	case "/key": // get by key
		key := req.Data // data holds the key bytes

		res.Key = key
		if !st.VersionExists(res.Height) {
			res.Log = iavl.ErrVersionDoesNotExist.Error()
			break
		}

		_, res.Value = tree.GetVersioned(key, res.Height)
		if !req.Prove {
			break
		}

		// Continue to prove existence/absence of value
		// Must convert store.Tree to iavl.MutableTree with given version to use in CreateProof
		iTree, err := tree.GetImmutable(res.Height)
		if err != nil {
			// sanity check: If value for given version was retrieved, immutable tree must also be retrievable
			panic(fmt.Sprintf("version exists in store but could not retrieve corresponding versioned tree in store, %s", err.Error()))
		}
		mtree := &iavl.MutableTree{
			ImmutableTree: iTree,
		}

		// get proof from tree and convert to merkle.Proof before adding to result
		res.ProofOps = getProofFromTree(mtree, req.Data, res.Value != nil)

	case "/subspace":
		var KVs []KVPair

		subspace := req.Data
		res.Key = subspace

		iterator := sdk.KVStorePrefixIterator(st, subspace)
		for ; iterator.Valid(); iterator.Next() {
			KVs = append(KVs, KVPair{Key: iterator.Key(), Value: iterator.Value()})
		}

		iterator.Close()
		res.Value = cdc.MustMarshalBinaryLengthPrefixed(KVs)

	default:
		msg := fmt.Sprintf("Unexpected Query path: %v", req.Path)
		return sdkerrors.QueryResult(sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, msg))
	}

	return
}
// Takes a MutableTree, a key, and a flag for creating existence or absence proof and returns the
// appropriate merkle.Proof. Since this must be called after querying for the value, this function should never error
// Thus, it will panic on error rather than returning it
func getProofFromTree(tree *iavl.MutableTree, key []byte, exists bool) *tmcrypto.ProofOps {
	var (
		commitmentProof *ics23.CommitmentProof
		err             error
	)

	if exists {
		// value was found
		commitmentProof, err = tree.GetMembershipProof(key)
		if err != nil {
			// sanity check: If value was found, membership proof must be creatable
			panic(fmt.Sprintf("unexpected value for empty proof: %s", err.Error()))
		}
	} else {
		// value wasn't found
		commitmentProof, err = tree.GetNonMembershipProof(key)
		if err != nil {
			// sanity check: If value wasn't found, nonmembership proof must be creatable
			panic(fmt.Sprintf("unexpected error for nonexistence proof: %s", err.Error()))
		}
	}

	op := NewIavlCommitmentOp(key, commitmentProof)
	return &tmcrypto.ProofOps{Ops: []tmcrypto.ProofOp{op.ProofOp()}}
}
//----------------------------------------

// Implements Iterator.
type iavlIterator struct {
	// Underlying store
	tree *iavl.ImmutableTree

	// Domain
	start, end []byte

	// Iteration order
	ascending bool

	// Channel to push iteration values.
	iterCh chan abci.EventAttribute

	// Close this to release goroutine.
	quitCh chan struct{}

	// Close this to signal that state is initialized.
	initCh chan struct{}

	//----------------------------------------
	// What follows are mutable state.
	mtx sync.Mutex

	invalid bool   // True once, true forever
	key     []byte // The current key
	value   []byte // The current value
}

func (iter *iavlIterator) Error() error {
	return nil
}

var _ Iterator = (*iavlIterator)(nil)

// newIAVLIterator will create a new iavlIterator.
// CONTRACT: Caller must release the iavlIterator, as each one creates a new
// goroutine.
func newIAVLIterator(tree *iavl.ImmutableTree, start, end []byte, ascending bool) *iavlIterator {
	iter := &iavlIterator{
		tree:      tree,
		start:     cp(start),
		end:       cp(end),
		ascending: ascending,
		iterCh:    make(chan abci.EventAttribute, 0), // Set capacity > 0?
		quitCh:    make(chan struct{}),
		initCh:    make(chan struct{}),
	}
	go iter.iterateRoutine()
	go iter.initRoutine()
	return iter
}

// Run this to funnel items from the tree to iterCh.
func (iter *iavlIterator) iterateRoutine() {
	iter.tree.IterateRange(
		iter.start, iter.end, iter.ascending,
		func(key, value []byte) bool {
			select {
			case <-iter.quitCh:
				return true // done with iteration.
			case iter.iterCh <- abci.EventAttribute{Key: key, Value: value}:
				return false // yay.
			}
		},
	)
	close(iter.iterCh) // done.
}

// Run this to fetch the first item.
func (iter *iavlIterator) initRoutine() {
	iter.receiveNext()
	close(iter.initCh)
}

// Implements Iterator.
func (iter *iavlIterator) Domain() (start, end []byte) {
	return iter.start, iter.end
}

// Implements Iterator.
func (iter *iavlIterator) Valid() bool {
	iter.waitInit()
	iter.mtx.Lock()

	validity := !iter.invalid
	iter.mtx.Unlock()
	return validity
}

// Implements Iterator.
func (iter *iavlIterator) Next() {
	iter.waitInit()
	iter.mtx.Lock()
	iter.assertIsValid(true)

	iter.receiveNext()
	iter.mtx.Unlock()
}

// Implements Iterator.
func (iter *iavlIterator) Key() []byte {
	iter.waitInit()
	iter.mtx.Lock()
	iter.assertIsValid(true)

	key := iter.key
	iter.mtx.Unlock()
	return key
}

// Implements Iterator.
func (iter *iavlIterator) Value() []byte {
	iter.waitInit()
	iter.mtx.Lock()
	iter.assertIsValid(true)

	val := iter.value
	iter.mtx.Unlock()
	return val
}

// Implements Iterator.
func (iter *iavlIterator) Close() error {
	close(iter.quitCh)
	return nil
}

//----------------------------------------

func (iter *iavlIterator) setNext(key, value []byte) {
	iter.assertIsValid(false)

	iter.key = key
	iter.value = value
}

func (iter *iavlIterator) setInvalid() {
	iter.assertIsValid(false)

	iter.invalid = true
}

func (iter *iavlIterator) waitInit() {
	<-iter.initCh
}

func (iter *iavlIterator) receiveNext() {
	kvPair, ok := <-iter.iterCh
	if ok {
		iter.setNext(kvPair.Key, kvPair.Value)
	} else {
		iter.setInvalid()
	}
}

// assertIsValid panics if the iterator is invalid. If unlockMutex is true,
// it also unlocks the mutex before panicing, to prevent deadlocks in code that
// recovers from panics
func (iter *iavlIterator) assertIsValid(unlockMutex bool) {
	if iter.invalid {
		if unlockMutex {
			iter.mtx.Unlock()
		}
		panic("invalid iterator")
	}
}

//----------------------------------------

func cp(bz []byte) (ret []byte) {
	if bz == nil {
		return nil
	}
	ret = make([]byte, len(bz))
	copy(ret, bz)
	return ret
}
