package baseapp

import (
	"encoding/hex"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/store"
	"github.com/ci123chain/ci123chain/pkg/transfer"
	"io"
	"runtime/debug"
	"strings"

	"github.com/pkg/errors"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/abci/version"
	distypes "github.com/ci123chain/ci123chain/pkg/distribution/types"
	iftypes "github.com/ci123chain/ci123chain/pkg/infrastructure/types"
	ordertypes "github.com/ci123chain/ci123chain/pkg/order/types"
	"github.com/ci123chain/ci123chain/pkg/snapshots"
	snapshottypes "github.com/ci123chain/ci123chain/pkg/snapshots/types"
	staktypes "github.com/ci123chain/ci123chain/pkg/staking/types"
	wasmtypes "github.com/ci123chain/ci123chain/pkg/vm/wasmtypes"
	"github.com/tendermint/tendermint/proto/tendermint/types"
)

// Key to store the header in the DB itself.
// Use the db directly instead of a store to avoid
// conflicts with handlers writing to the store
// and to avoid affecting the Merkle root.
var dbHeaderKey = []byte("header")

// Enum mode for app.runTx
type runTxMode uint8

const (
	// Check a transfer
	runTxModeCheck runTxMode = iota
	// Simulate a transfer
	runTxModeSimulate runTxMode = iota
	// Deliver a transfer
	runTxModeDeliver runTxMode = iota
)

type Committer func(ctx sdk.Context) abci.ResponseCommit

// BaseApp reflects the ABCI application implementation.
type BaseApp struct {
	// initialized on creation
	Logger      log.Logger
	name        string               // application name from abci.Info
	db          dbm.DB               // common DB backend
	cms         sdk.CommitMultiStore // Main (uncached) state
	queryRouter QueryRouter          // router for redirecting query calls
	//handler     sdk.Handler
	router 		sdk.Router

	txDecoder   sdk.TxDecoder // unmarshal []byte into sdk.Tx

	anteHandler sdk.AnteHandler // ante handler for fee and auth
	deferHandler sdk.DeferHandler // defer handler for fee and auth

	// may be nil
	initChainer      sdk.InitChainer  // initialize state with validators and state blob
	beginBlocker     sdk.BeginBlocker // logic to run before any txs
	endBlocker       sdk.EndBlocker   // logic to run after all txs, and to determine valset changes
	addrPeerFilter   sdk.PeerFilter   // filter peers by address and port
	pubkeyPeerFilter sdk.PeerFilter   // filter peers by public types
	committer 		 Committer
	//--------------------
	// Volatile
	// checkState is set on initialization and reset on Commit.
	// deliverState is set in InitChain and BeginBlock and cleared on Commit.
	// See methods setCheckState and setDeliverState.
	checkState   *state          // for CheckTx
	deliverState *state          // for DeliverTx
	voteInfos    []abci.VoteInfo // absent validators from begin block

	// flag for sealing
	sealed bool

	// manages snapshots, i.e. dumps of app state at certain intervals
	snapshotManager    *snapshots.Manager
	snapshotInterval   uint64 // block interval between state sync snapshots
	snapshotKeepRecent uint32 // recent state sync snapshots to keep
}

func (app *BaseApp) ListSnapshots(req abci.RequestListSnapshots) abci.ResponseListSnapshots {
	resp := abci.ResponseListSnapshots{Snapshots: []*abci.Snapshot{}}
	if app.snapshotManager == nil {
		return resp
	}

	snapshots, err := app.snapshotManager.List()
	if err != nil {
		app.Logger.Error("failed to list snapshots", "err", err)
		return resp
	}

	for _, snapshot := range snapshots {
		abciSnapshot, err := snapshot.ToABCI()
		if err != nil {
			app.Logger.Error("failed to list snapshots", "err", err)
			return resp
		}
		resp.Snapshots = append(resp.Snapshots, &abciSnapshot)
	}

	return resp
}

func (app *BaseApp) OfferSnapshot(req abci.RequestOfferSnapshot) abci.ResponseOfferSnapshot {
	if app.snapshotManager == nil {
		app.Logger.Error("snapshot manager not configured")
		return abci.ResponseOfferSnapshot{Result: abci.ResponseOfferSnapshot_ABORT}
	}

	if req.Snapshot == nil {
		app.Logger.Error("received nil snapshot")
		return abci.ResponseOfferSnapshot{Result: abci.ResponseOfferSnapshot_REJECT}
	}

	snapshot, err := snapshottypes.SnapshotFromABCI(req.Snapshot)
	if err != nil {
		app.Logger.Error("failed to decode snapshot metadata", "err", err)
		return abci.ResponseOfferSnapshot{Result: abci.ResponseOfferSnapshot_REJECT}
	}

	err = app.snapshotManager.Restore(snapshot)
	switch {
	case err == nil:
		return abci.ResponseOfferSnapshot{Result: abci.ResponseOfferSnapshot_ACCEPT}

	case errors.Is(err, snapshottypes.ErrUnknownFormat):
		return abci.ResponseOfferSnapshot{Result: abci.ResponseOfferSnapshot_REJECT_FORMAT}

	case errors.Is(err, snapshottypes.ErrInvalidMetadata):
		app.Logger.Error(
			"rejecting invalid snapshot",
			"height", req.Snapshot.Height,
			"format", req.Snapshot.Format,
			"err", err,
		)
		return abci.ResponseOfferSnapshot{Result: abci.ResponseOfferSnapshot_REJECT}

	default:
		app.Logger.Error(
			"failed to restore snapshot",
			"height", req.Snapshot.Height,
			"format", req.Snapshot.Format,
			"err", err,
		)

		// We currently don't support resetting the IAVL stores and retrying a different snapshot,
		// so we ask Tendermint to abort all snapshot restoration.
		return abci.ResponseOfferSnapshot{Result: abci.ResponseOfferSnapshot_ABORT}
	}
}

func (app *BaseApp) LoadSnapshotChunk(req abci.RequestLoadSnapshotChunk) abci.ResponseLoadSnapshotChunk {
	if app.snapshotManager == nil {
		return abci.ResponseLoadSnapshotChunk{}
	}
	chunk, err := app.snapshotManager.LoadChunk(req.Height, req.Format, req.Chunk)
	if err != nil {
		app.Logger.Error(
			"failed to load snapshot chunk",
			"height", req.Height,
			"format", req.Format,
			"chunk", req.Chunk,
			"err", err,
		)
		return abci.ResponseLoadSnapshotChunk{}
	}
	return abci.ResponseLoadSnapshotChunk{Chunk: chunk}
}

func (app *BaseApp) ApplySnapshotChunk(req abci.RequestApplySnapshotChunk) abci.ResponseApplySnapshotChunk {
	if app.snapshotManager == nil {
		app.Logger.Error("snapshot manager not configured")
		return abci.ResponseApplySnapshotChunk{Result: abci.ResponseApplySnapshotChunk_ABORT}
	}

	_, err := app.snapshotManager.RestoreChunk(req.Chunk)
	switch {
	case err == nil:
		return abci.ResponseApplySnapshotChunk{Result: abci.ResponseApplySnapshotChunk_ACCEPT}

	case errors.Is(err, snapshottypes.ErrChunkHashMismatch):
		app.Logger.Error(
			"chunk checksum mismatch; rejecting sender and requesting refetch",
			"chunk", req.Index,
			"sender", req.Sender,
			"err", err,
		)
		return abci.ResponseApplySnapshotChunk{
			Result:        abci.ResponseApplySnapshotChunk_RETRY,
			RefetchChunks: []uint32{req.Index},
			RejectSenders: []string{req.Sender},
		}

	default:
		app.Logger.Error("failed to restore snapshot", "err", err)
		return abci.ResponseApplySnapshotChunk{Result: abci.ResponseApplySnapshotChunk_ABORT}
	}
}

var _ abci.Application = (*BaseApp)(nil)

// NewBaseApp returns a reference to an initialized BaseApp.
//
// TODO: Determine how to use a flexible and robust configuration paradigm that
// allows for sensible defaults while being highly configurable
// (e.g. functional options).
//
// NOTE: The db is used to store the version number for now.
// Accepts a user-defined txDecoder
// Accepts variable number of option functions, which act on the BaseApp to set configuration choices
func NewBaseApp(name string, logger log.Logger, ldb dbm.DB, cdb dbm.DB, cacheDir string, txDecoder sdk.TxDecoder, options ...func(*BaseApp)) *BaseApp {
	app := &BaseApp{
		Logger:      logger,
		name:        name,
		cms:         store.NewCommitMultiStore(ldb, cdb, cacheDir),
		//cms:         store.NewBaseMultiStore(db),
		queryRouter: NewQueryRouter(),
		router:      sdk.NewRouter(),
		txDecoder:   txDecoder,
	}

	for _, option := range options {
		option(app)
	}
	return app
}

// BaseApp Name
func (app *BaseApp) Name() string {
	return app.name
}

// SetCommitMultiStoreTracer sets the store tracer on the BaseApp's underlying
// CommitMultiStore.
func (app *BaseApp) SetCommitMultiStoreTracer(w io.Writer) {
	app.cms.WithTracer(w)
}

// Mount IAVL stores to the provided keys in the BaseApp multistore
func (app *BaseApp) MountStoresIAVL(keys ...*sdk.KVStoreKey) {
	for _, key := range keys {
		app.MountStore(key, sdk.StoreTypeIAVL)
	}
}

// Mount stores to the provided keys in the BaseApp multistore
func (app *BaseApp) MountStoresTransient(keys ...*sdk.TransientStoreKey) {
	for _, key := range keys {
		app.MountStore(key, sdk.StoreTypeTransient)
	}
}

// Mount a store to the provided types in the BaseApp multistore, using a specified DB
func (app *BaseApp) MountStoreWithDB(key sdk.StoreKey, typ sdk.StoreType, db dbm.DB) {
	app.cms.MountStoreWithDB(key, typ, db)
}

// Mount a store to the provided types in the BaseApp multistore, using the default DB
func (app *BaseApp) MountStore(key sdk.StoreKey, typ sdk.StoreType) {
	app.cms.MountStoreWithDB(key, typ, nil)
}

func (app *BaseApp) GetLatestVersion() int64 {
	return app.cms.GetLatestVersion()
}

// load latest application version
func (app *BaseApp) LoadLatestVersion(mainKey sdk.StoreKey) error {
	err := app.cms.LoadLatestVersion()
	if err != nil {
		return err
	}
	return app.initFromStore(mainKey)
}

// load application version
func (app *BaseApp) LoadVersion(version int64, mainKey sdk.StoreKey) error {
	err := app.cms.LoadVersion(version)
	if err != nil {
		return err
	}
	return app.initFromStore(mainKey)
}

// the last CommitID of the multistore
func (app *BaseApp) LastCommitID() sdk.CommitID {
	return app.cms.LastCommitID()
}

// the last committed block height
func (app *BaseApp) LastBlockHeight() int64 {
	return app.cms.LastCommitID().Version
}

// SetHandler sets tx handler
//func (app *BaseApp) SetHandler(h sdk.Handler) {
//	app.handler = h
//}

// initializes the remaining logic from app.cms
func (app *BaseApp) initFromStore(mainKey sdk.StoreKey) error {
	// main store should exist.
	// TODO: we don't actually need the main store here
	main := app.cms.GetKVStore(mainKey)
		if main == nil {
		return errors.New("baseapp expects MultiStore with 'main' KVStore")
	}
	// Needed for `gaiad export`, which inits from store but never calls initchain
	app.setCheckState(types.Header{})

	app.Seal()

	return nil
}

// NewContext returns a new Context with the correct store, the given header, and nil txBytes.
func (app *BaseApp) NewContext(isCheckTx bool, header types.Header) sdk.Context {
	if isCheckTx {
		return sdk.NewContext(app.checkState.ms, header, true, app.Logger)
	}
	return sdk.NewContext(app.deliverState.ms, header, false, app.Logger)
}

type state struct {
	ms  sdk.CacheMultiStore
	ctx sdk.Context
}

func (st *state) CacheMultiStore() sdk.CacheMultiStore {
	return st.ms.CacheMultiStore()
}

func (st *state) Context() sdk.Context {
	return st.ctx
}

func (app *BaseApp) setCheckState(header types.Header) {
	ms := app.cms.CacheMultiStore()
	app.checkState = &state{
		ms:  ms,
		ctx: sdk.NewContext(ms, header, true, app.Logger),
	}
}

func (app *BaseApp) setDeliverState(header types.Header) {
	ms := app.cms.CacheMultiStore()
	app.deliverState = &state{
		ms:  ms,
		ctx: sdk.NewContext(ms, header, false, app.Logger),
	}
}

//______________________________________________________________________________

// ABCI

// Implements ABCI
func (app *BaseApp) Info(req abci.RequestInfo) abci.ResponseInfo {
	lastCommitID := app.cms.LastCommitID()

	return abci.ResponseInfo{
		Data:             app.name,
		LastBlockHeight:  lastCommitID.Version,
		LastBlockAppHash: lastCommitID.Hash,
	}
}

// Implements ABCI
func (app *BaseApp) SetOption(req abci.RequestSetOption) (res abci.ResponseSetOption) {
	// TODO: Implement
	return
}

// Implements ABCI
// InitChain runs the initialization logic directly on the CommitMultiStore and commits it.
func (app *BaseApp) InitChain(req abci.RequestInitChain) (res abci.ResponseInitChain) {
	// Initialize the deliver state and check state with ChainID and run initChain
	app.setDeliverState(types.Header{ChainID: req.ChainId})
	app.setCheckState(types.Header{ChainID: req.ChainId})

	if app.initChainer == nil {
		return
	}
	res = app.initChainer(app.deliverState.ctx, req)

	// NOTE: we don't commit, but BeginBlock for block 1
	// starts from this deliverState
	return
}

// Filter peers by address / port
func (app *BaseApp) FilterPeerByAddrPort(info string) abci.ResponseQuery {
	if app.addrPeerFilter != nil {
		return app.addrPeerFilter(info)
	}
	return abci.ResponseQuery{}
}

// Filter peers by public types
func (app *BaseApp) FilterPeerByPubKey(info string) abci.ResponseQuery {
	if app.pubkeyPeerFilter != nil {
		return app.pubkeyPeerFilter(info)
	}
	return abci.ResponseQuery{}
}

// Splits a string path using the delimter '/'.  i.e. "this/is/funny" becomes []string{"this", "is", "funny"}
func splitPath(requestPath string) (path []string) {
	path = strings.Split(requestPath, "/")
	// first element is empty string
	if len(path) > 0 && path[0] == "" {
		path = path[1:]
	}
	return path
}

// Implements ABCI.
// Delegates to CommitMultiStore if it implements Queryable
func (app *BaseApp) Query(req abci.RequestQuery) (res abci.ResponseQuery) {
	path := splitPath(req.Path)
	if len(path) == 0 {
		msg := "no query path provided"
		return sdk.ErrUnknownRequest(msg).QueryResult()
	}
	switch path[0] {
	// "/app" prefix for special application queries
	case "app":
		return handleQueryApp(app, path, req)
	case "store":
		return handleQueryStore(app, path, req)
	case "p2p":
		return handleQueryP2P(app, path, req)
	case "custom":
		return handleQueryCustom(app, path, req)
	}

	msg := "unknown query path"
	return sdk.ErrUnknownRequest(msg).QueryResult()
}

func handleQueryApp(app *BaseApp, path []string, req abci.RequestQuery) (res abci.ResponseQuery) {
	if len(path) >= 2 {
		var result sdk.Result
		switch path[1] {
		case "simulate":
			txBytes := req.Data
			tx, err := app.txDecoder(txBytes)
			if err != nil {
				result = err.Result()
			} else {
				result = app.Simulate(tx)
			}
			simulationResp := sdk.QureyAppResponse{
				Code:       uint32(result.Code),
				FormatData: string(result.Data),
				Data:       strings.ToUpper(hex.EncodeToString(result.Data)),
				Log:        result.Log,
				GasWanted:  result.GasWanted,
				GasUsed:    result.GasUsed,
				Codespace:  string(result.Codespace),
			}
			value := codec.Cdc.MustMarshalBinaryBare(simulationResp)
			return abci.ResponseQuery{
				Code:      uint32(sdk.CodeOK),
				Codespace: string(sdk.CodespaceRoot),
				Value:     value,
			}
		case "version":
			return abci.ResponseQuery{
				Code:      uint32(sdk.CodeOK),
				Codespace: string(sdk.CodespaceRoot),
				Value:     []byte(version.GetVersion()),
			}
		default:
			result = sdk.ErrUnknownRequest(fmt.Sprintf("Unknown query: %s", path)).Result()
		}

		// Encode with json
		value := codec.Cdc.MustMarshalBinaryLengthPrefixed(result)
		return abci.ResponseQuery{
			Code:      uint32(sdk.CodeOK),
			Codespace: string(sdk.CodespaceRoot),
			Value:     value,
		}
	}
	msg := "Expected second parameter to be either simulate or version, neither was present"
	return sdk.ErrUnknownRequest(msg).QueryResult()
}

func handleQueryStore(app *BaseApp, path []string, req abci.RequestQuery) (res abci.ResponseQuery) {
	// "/store" prefix for store queries
	queryable, ok := app.cms.(sdk.Queryable)
	if !ok {
		msg := "multistore doesn't support queries"
		return sdk.ErrUnknownRequest(msg).QueryResult()
	}
	req.Path = "/" + strings.Join(path[1:], "/")
	return queryable.Query(req)
}

// nolint: unparam
func handleQueryP2P(app *BaseApp, path []string, req abci.RequestQuery) (res abci.ResponseQuery) {
	// "/p2p" prefix for p2p queries
	if len(path) >= 4 {
		if path[1] == "filter" {
			if path[2] == "addr" {
				return app.FilterPeerByAddrPort(path[3])
			}
			if path[2] == "pubkey" {
				// TODO: this should be changed to `id`
				// NOTE: this changed in tendermint and we didn't notice...
				return app.FilterPeerByPubKey(path[3])
			}
		} else {
			msg := "Expected second parameter to be filter"
			return sdk.ErrUnknownRequest(msg).QueryResult()
		}
	}

	msg := "Expected path is p2p filter <addr|pubkey> <parameter>"
	return sdk.ErrUnknownRequest(msg).QueryResult()
}

func handleQueryCustom(app *BaseApp, path []string, req abci.RequestQuery) (res abci.ResponseQuery) {
	// path[0] should be "custom" because "/custom" prefix is required for keeper queries.
	// the queryRouter routes using path[1]. For example, in the path "custom/gov/proposal", queryRouter routes using "gov"
	if len(path) < 2 || path[1] == "" {
		return sdk.ErrUnknownRequest("No route for custom query specified").QueryResult()
	}
	querier := app.queryRouter.Route(path[1])
	if querier == nil {
		return sdk.ErrUnknownRequest(fmt.Sprintf("no custom querier found for route %s", path[1])).QueryResult()
	}

	// Cache wrap the commit-multistore for safety.
	ctx := sdk.NewContext(app.cms, app.checkState.ctx.BlockHeader(), true, app.Logger)

	// Passes the rest of the path as an argument to the querier.
	// For example, in the path "custom/gov/proposal/test", the gov querier gets []string{"proposal", "test"} as the path
	resBytes, err := querier(ctx, path[2:], req)
	if err != nil {
		return abci.ResponseQuery{
			Code:      uint32(err.Code()),
			Codespace: string(err.Codespace()),
			Log:       err.ABCILog(),
		}
	}
	return abci.ResponseQuery{
		Code:  uint32(sdk.CodeOK),
		Value: resBytes,
	}
}

// BeginBlock implements the ABCI application interface.
func (app *BaseApp) BeginBlock(req abci.RequestBeginBlock) (res abci.ResponseBeginBlock) {
	if app.cms.TracingEnabled() {
		app.cms.ResetTraceContext()
		app.cms.WithTracingContext(sdk.TraceContext(
			map[string]interface{}{"blockHeight": req.Header.Height},
		))
	}

	// Initialize the DeliverTx state. If this is the first block, it should
	// already be initialized in InitChain. Otherwise app.deliverState will be
	// nil, since it is reset on Commit.
	if app.deliverState == nil {
		app.setDeliverState(req.Header)
	} else {
		// In the first block, app.deliverState.ctx will already be initialized
		// by InitChain. Context is now updated with Header information.
		app.deliverState.ctx = app.deliverState.ctx.WithBlockHeader(req.Header).WithBlockHeight(req.Header.Height)
	}

	if app.beginBlocker != nil {
		res = app.beginBlocker(app.deliverState.ctx, req)
	}

	// set the signed validators for addition to context in deliverTx
	// TODO: communicate this result to the address to pubkey map in slashing
	app.voteInfos = req.LastCommitInfo.GetVotes()
	return
}

// CheckTx implements the ABCI interface. It runs the "basic checks" to see
// whether or not a transfer can possibly be executed, first decoding, then
// the ante handler (which checks signatures/fees/ValidateBasic), then finally
// the route match to see whether a handler exists.
//
// NOTE:CheckTx does not run the actual Tx handler function(s).
func (app *BaseApp) CheckTx(req abci.RequestCheckTx) (res abci.ResponseCheckTx) {
	var result sdk.Result

	tx, err := app.txDecoder(req.Tx)
	if err != nil {
		result = err.Result()
	} else {
		result = app.runTx(runTxModeCheck, req.Tx, tx)
	}
	return abci.ResponseCheckTx{
		Code:      uint32(result.Code),
		Codespace: string(result.Codespace),
		Data:      result.Data,
		Log:       result.Log,
		GasWanted: int64(result.GasWanted), // TODO: Should types accept unsigned ints?
		GasUsed:   int64(result.GasUsed),   // TODO: Should types accept unsigned ints?
		Events:    result.Events.ToABCIEvents(),
	}
}

// DeliverTx implements the ABCI interface.
func (app *BaseApp) DeliverTx(req abci.RequestDeliverTx) (res abci.ResponseDeliverTx) {
	var result sdk.Result

	tx, err := app.txDecoder(req.Tx)
	if err != nil {
		result = err.Result()
	} else {
		result = app.runTx(runTxModeDeliver, req.Tx, tx)
	}

	return abci.ResponseDeliverTx{
		Code:      uint32(result.Code),
		Codespace: string(result.Codespace),
		Data:      result.Data,
		Log:       result.Log,
		GasWanted: int64(result.GasWanted), // TODO: Should types accept unsigned ints?
		GasUsed:   int64(result.GasUsed),   // TODO: Should types accept unsigned ints?
		Events:    result.Events.ToABCIEvents(),
	}
}

// retrieve the context for the tx w/ txBytes and other memoized values.
func (app *BaseApp) getContextForTx(mode runTxMode, txBytes []byte) (ctx sdk.Context) {
	ctx = app.getState(mode).ctx.
		WithTxBytes(txBytes).
		WithVoteInfos(app.voteInfos)
	if mode == runTxModeSimulate {
		ctx, _ = ctx.CacheContext()
	}
	return
}

func (app *BaseApp) runMsgs(ctx sdk.Context, msgs []sdk.Msg, mode runTxMode) sdk.Result {
	idxLogs := make([]sdk.ABCIMessageLog, 0, len(msgs))
	var code      sdk.CodeType
	var codespace sdk.CodespaceType
	var data []byte
	var events sdk.Events

	for msgIdx, msg := range msgs{
		msgRoute := msg.Route()
		handler := app.router.Route(msgRoute)
		if handler == nil {
			return sdk.ErrUnknownRequest("Unrecognized Msg types: " + msgRoute).Result()
		}

		var msgResult sdk.Result

		// check 不实际执行
		if mode != runTxModeCheck {
			msgResult = handler(ctx, msg)
		}

		// append date and log and events
		data = append(data, msgResult.Data...)
		idxLog := sdk.ABCIMessageLog{MsgIndex: uint16(msgIdx), Log: msgResult.Log}
		events = append(events, msgResult.Events...)

		if !msgResult.IsOK() {
			idxLog.Success = false
			idxLogs = append(idxLogs, idxLog)

			code = msgResult.Code
			codespace = msgResult.Codespace
			break
		}

		idxLog.Success = true
		idxLogs = append(idxLogs, idxLog)
	}

	logJSON := codec.Cdc.MustMarshalJSON(idxLogs)
	return sdk.Result{
		Code: 		code,
		Codespace: 	codespace,
		GasUsed:   	ctx.GasMeter().GasConsumed(),
		Log: 		strings.TrimSpace(string(logJSON)),
		Data: 		data,
		Events: 	events,
	}
}

// Returns the applicantion's deliverState if app is in runTxModeDeliver,
// otherwise it returns the application's checkstate.
func (app *BaseApp) getState(mode runTxMode) *state {
	if mode == runTxModeCheck || mode == runTxModeSimulate {
		return app.checkState
	}

	return app.deliverState
}

// cacheTxContext returns a new context based off of the provided context with
// a cache wrapped multi-store.
func (app *BaseApp) cacheTxContext(ctx sdk.Context, txBytes []byte) (
	sdk.Context, sdk.CacheMultiStore) {

	ms := ctx.MultiStore()
	// TODO: https://github.com/bluele/hypermint/pkg/abci/issues/2824
	msCache := ms.CacheMultiStore()
	if msCache.TracingEnabled() {
		msCache = msCache.WithTracingContext(
			sdk.TraceContext(
				map[string]interface{}{
					"txHash": fmt.Sprintf("%X", tmhash.Sum(txBytes)),
				},
			),
		).(sdk.CacheMultiStore)
	}

	return ctx.WithMultiStore(msCache), msCache
}

// runTx processes a transfer. The transactions is proccessed via an
// anteHandler. The provided txBytes may be nil in some cases, eg. in tests. For
// further details on transfer execution, reference the BaseApp SDK
// documentation.
func (app *BaseApp) runTx(mode runTxMode, txBytes []byte, tx sdk.Tx) (result sdk.Result) {
	// NOTE: GasWanted should be returned by the AnteHandler. GasUsed is
	// determined by the GasMeter. We need access to the context to get the gas
	// meter so we initialize upfront.
	var gasWanted uint64
	//var gasUsed uint64
	ctx := app.getContextForTx(mode, txBytes)
	ms := ctx.MultiStore()
	gasWanted = tx.GetGas()
	ctx = ctx.WithGasLimit(gasWanted)
	ctx = ctx.WithNonce(tx.GetNonce())

	var all_attributes = make([][]sdk.Attribute, 0)
	defer func() {
		if r := recover(); r != nil {
			switch rType := r.(type) {
			case sdk.ErrorOutOfGas:
				app.deferHandler(ctx, tx, true, mode == runTxModeSimulate)
				newLog := fmt.Sprintf("out of gas in location: %v", rType.Descriptor)
				result = sdk.ErrOutOfGas(newLog).Result()
				result.GasUsed = gasWanted
			default:
				res := app.deferHandler(ctx, tx, false, mode == runTxModeSimulate)
				newLog := fmt.Sprintf("recovered: %v\nstack:%v\n", r, string(debug.Stack()))
				result = sdk.ErrInternal(newLog).Result()
				result.GasUsed = res.GasUsed
			}
			for _, v := range all_attributes {
				v = append(v, sdk.NewAttribute([]byte(sdk.EventTypeType), []byte(sdk.AttributeKeyInvalidTx)))
				event := sdk.NewEvent(sdk.AttributeKeyTx, v...)
				result.Events = append(result.Events, event)
			}
		} else {
			res := app.deferHandler(ctx, tx, false, mode == runTxModeSimulate)
			result.GasUsed = res.GasUsed
			for _, v := range all_attributes {
				v = append(v, sdk.NewAttribute([]byte(sdk.EventTypeType), []byte(sdk.AttributeKeyValidTx)))
				event := sdk.NewEvent(sdk.AttributeKeyTx, v...)
				result.Events = append(result.Events, event)
			}
		}
		result.GasWanted = gasWanted
	}()
	msgs := tx.GetMsgs()
	if mode != runTxModeSimulate {
		signer := tx.GetFromAddress()
		err := validateBasicTxMsgs(msgs, signer)
		if err != nil {
			return err.Result()
		}
		if err := tx.ValidateBasic(); err != nil {
			return err.Result()
		}
	}

	// Execute the ante handler if one is defined.
	if app.anteHandler != nil {
		var anteCtx sdk.Context
		var msCache sdk.CacheMultiStore

		// Cache wrap context before anteHandler call in case it aborts.
		// This is required for both CheckTx and DeliverTx.
		// https://github.com/bluele/hypermint/pkg/abci/issues/2772
		// NOTE: Alternatively, we could require that anteHandler ensures that
		// writes do not happen if aborted/failed.  This may have some
		// performance benefits, but it'll be more difficult to get right.
		anteCtx, msCache = app.cacheTxContext(ctx, txBytes)
		newCtx, result, abort := app.anteHandler(anteCtx, tx, (mode == runTxModeSimulate))

		if !newCtx.IsZero() {
			// At this point, newCtx.MultiStore() is cache wrapped,
			// or something else replaced by anteHandler.
			// We want the original ms, not one which was cache-wrapped
			// for the ante handler.
			ctx = newCtx.WithMultiStore(ms)
		}
		//gasWanted = result.GasWanted
		//gasUsed = result.GasUsed
		if abort {
			return result
		}

		msCache.Write()
	}

	if mode == runTxModeCheck { // XXX
		return
	}

	// Create a new context based off of the existing context with a cache wrapped
	// multi-store in case message processing fails.

	all_attributes = allMsgAttributes(msgs)

	runMsgCtx, msCache := app.cacheTxContext(ctx, txBytes)
	result = app.runMsgs(runMsgCtx, msgs, mode)

	if mode == runTxModeSimulate  { // XXX
		return
	}

	// only update state if all messages pass
	if result.IsOK() {
		msCache.Write()
	}
	return result
}

// EndBlock implements the ABCI application interface.
func (app *BaseApp) EndBlock(req abci.RequestEndBlock) (res abci.ResponseEndBlock) {
	if app.deliverState.ms.TracingEnabled() {
		app.deliverState.ms = app.deliverState.ms.ResetTraceContext().(sdk.CacheMultiStore)
	}

	if app.endBlocker != nil {
		res = app.endBlocker(app.deliverState.ctx, req)
	}

	if app.committer != nil {
		app.committer(app.deliverState.ctx)
	}
	return res
}

// Implements ABCI
func (app *BaseApp) Commit() (res abci.ResponseCommit) {
	header := app.deliverState.ctx.BlockHeader()
	/*
		// Write the latest Header to the store
			headerBytes, err := proto.Marshal(&header)


			if err != nil {
				panic(err)
			}
			app.db.SetSync(dbHeaderKey, headerBytes)
	*/
	// Write the Deliver state and commit the MultiStore
	app.deliverState.ms.Write()
	commitID := app.cms.Commit()

	// TODO: this is missing a module identifier and dumps byte array
	app.Logger.Debug("Commit synced",
		"commit", commitID,
	)

	// Reset the Check state to the latest committed
	// NOTE: safe because Tendermint holds a lock on the mempool for Commit.
	// Use the header from this latest block.
	app.setCheckState(header)

	// Empty the Deliver state
	app.deliverState = nil
	return abci.ResponseCommit{
		Data: commitID.Hash,
	}
}

func validateBasicTxMsgs(msgs []sdk.Msg, signer sdk.AccAddress) sdk.Error {
	if msgs == nil || len(msgs) == 0 {
		return sdk.ErrUnknownRequest("Tx.GetMsgs() must return at least one message in list")
	}

	for _, msg := range msgs {
		// Validate the Msg.
		if signer != msg.GetFromAddress(){
			return sdk.ErrInvalidSign("Signer is different from msg.from")
		}
		err := msg.ValidateBasic()
		if err != nil {
			return err
		}
	}

	return nil
}

func allMsgAttributes(msgs []sdk.Msg) [][]sdk.Attribute {
	var attributes = make([][]sdk.Attribute, 0)
	var multiMsg string
	if len(msgs) == 1{
		multiMsg = "false"
	}else {
		multiMsg = "true"
	}
	for _, v := range msgs {
		var operation = ""
		var amount = ""
		var receiver = "0x0000000000000000000000000000000000000000"
		var module = ""
		var attrs = make([]sdk.Attribute, 0)

		switch vt := v.(type) {
		case *transfer.MsgTransfer:
			operation = "transfer"
			amount = vt.Amount.Amount.String()
			receiver = vt.To.String()
			module = transfer.AttributeValueCategory
		case *distypes.MsgSetWithdrawAddress:
			operation = "modify_withdraw_address"
			module = distypes.AttributeValueCategory
		case *distypes.MsgFundCommunityPool:
			operation = "fund_community_pool"
			amount = vt.Amount.Amount.String()
			module = distypes.AttributeValueCategory
		case *distypes.MsgWithdrawDelegatorReward:
			operation = "withdraw_rewards"
			module = distypes.AttributeValueCategory
		case *distypes.MsgWithdrawValidatorCommission:
			operation = "withdraw_commission"
			module = distypes.AttributeValueCategory
		case *staktypes.MsgEditValidator:
			operation = "edit_validator"
			module = staktypes.AttributeValueCategory
		case *staktypes.MsgCreateValidator:
			operation = "create_validator"
			amount = vt.Value.Amount.String()
			module = staktypes.AttributeValueCategory
		case *staktypes.MsgDelegate:
			operation = "delegate"
			amount = vt.Amount.Amount.String()
			module = staktypes.AttributeValueCategory
		case *staktypes.MsgRedelegate:
			operation = "redelegate"
			amount = vt.Amount.Amount.String()
			module = staktypes.AttributeValueCategory
		case *staktypes.MsgUndelegate:
			operation = "undelegate"
			amount = vt.Amount.Amount.String()
			module = staktypes.AttributeValueCategory
		case *wasmtypes.MsgExecuteContract:
			operation = "invoke_contract"
			receiver = vt.Contract.String()
			module = staktypes.AttributeValueCategory
		case *wasmtypes.MsgInstantiateContract:
			operation = "init_contract"
			module = staktypes.AttributeValueCategory
		case *wasmtypes.MsgMigrateContract:
			operation = "migrate_contract"
			receiver = vt.Contract.String()
			module = staktypes.AttributeValueCategory
		case *wasmtypes.MsgUploadContract:
			operation = "upload_contract"
			module = staktypes.AttributeValueCategory
		case *ordertypes.MsgUpgrade:
			operation = "add_shard"
			module = ordertypes.AttributeValueCategory
		//case *ibctypes.MsgApplyIBC:
		//	operation = "apply_ibc"
		//	module = ibctypes.ModuleName
		case *iftypes.MsgStoreContent:
			operation = "store_content"
			module = iftypes.ModuleName
		}

		attrs = sdk.NewAttributes(attrs,
			sdk.NewAttribute([]byte(sdk.AttributeKeySender), []byte(v.GetFromAddress().String())),
			sdk.NewAttribute([]byte(sdk.AttributeKeyReceiver), []byte(receiver)),
			sdk.NewAttribute([]byte(sdk.AttributeKeyMethod),[]byte(operation)),
			sdk.NewAttribute([]byte(sdk.AttributeKeyAmount), []byte(amount)),
			sdk.NewAttribute([]byte(sdk.AttributeKeyModule), []byte(module)),
			sdk.NewAttribute([]byte(sdk.EventTypeMultiMsg), []byte(multiMsg)),
		)

		attributes = append(attributes, attrs)
	}
	return attributes
}