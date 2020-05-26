package baseapp

import (
	"fmt"
	"io"
	"runtime/debug"
	"strings"

	"github.com/pkg/errors"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	"github.com/ci123chain/ci123chain/pkg/abci/store"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/abci/version"
	"github.com/ci123chain/ci123chain/pkg/transaction"
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
	commitDB    dbm.DB               //commit info DB
	cms         sdk.CommitMultiStore // Main (uncached) state
	cfcms       sdk.CommitMultiStore  //commitInfo
	keys        []*sdk.KVStoreKey        //normal key
	configKeys  []*sdk.KVStoreKey        //config key
	queryRouter QueryRouter          // router for redirecting query calls
	//handler     sdk.Handler
	router 		Router

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
func NewBaseApp(name string, logger log.Logger, db, commitDB dbm.DB, keys, configKeys []*sdk.KVStoreKey, txDecoder sdk.TxDecoder, options ...func(*BaseApp)) *BaseApp {
	app := &BaseApp{
		Logger:      logger,
		name:        name,
		db:          db,
		commitDB:    commitDB,
		keys:        keys,
		configKeys:  configKeys,
		//cms:         store.NewCommitMultiStore(db),
		cms:         store.NewBaseMultiStore(db),
		cfcms:       store.NewBaseMultiStore(commitDB),
		queryRouter: NewQueryRouter(),
		router: 	 NewRouter(),
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

func (app *BaseApp) MountStoreIAVLWithDB(isConfigKey bool,keys ...*sdk.KVStoreKey) {
	var db dbm.DB
	if isConfigKey {
		db = app.commitDB
		for _, key := range keys {
			app.MountStoreWithConfigDB(key, sdk.StoreTypeIAVL, db)
		}
	}else {
		db = app.db
		for _, key := range keys {
			app.MountStoreWithDB(key, sdk.StoreTypeIAVL, db)
		}
	}
}

// Mount stores to the provided keys in the BaseApp multistore
func (app *BaseApp) MountStoresTransient(keys ...*sdk.TransientStoreKey) {
	for _, key := range keys {
		//app.MountStore(key, sdk.StoreTypeTransient)
		app.MountStoreWithDB(key, sdk.StoreTypeTransient, app.db)
	}
}

// Mount a store to the provided types in the BaseApp multistore, using a specified DB
func (app *BaseApp) MountStoreWithDB(key sdk.StoreKey, typ sdk.StoreType, db dbm.DB) {
	app.cms.MountStoreWithDB(key, typ, db)
}

func (app *BaseApp) MountStoreWithConfigDB(key sdk.StoreKey, typ sdk.StoreType, db dbm.DB) {
	app.cfcms.MountStoreWithDB(key, typ, db)
}

// Mount a store to the provided types in the BaseApp multistore, using the default DB
func (app *BaseApp) MountStore(key sdk.StoreKey, typ sdk.StoreType) {
	app.cms.MountStoreWithDB(key, typ, nil)
}

// load latest application version
func (app *BaseApp) LoadLatestVersion(mainKey sdk.StoreKey, isConfigKey bool) error {
	var err error
	if isConfigKey {
		err = app.cfcms.LoadLatestVersion()
	}else {
		err = app.cms.LoadLatestVersion()
	}
	if err != nil {
		return err
	}
	return app.initFromStore(mainKey, isConfigKey)
}

// load application version
func (app *BaseApp) LoadVersion(version int64, mainKey sdk.StoreKey, isConfigKey bool) error {
	var err error
	if isConfigKey {
		err = app.cfcms.LoadVersion(version)
	}else {
		err = app.cms.LoadVersion(version)
	}
	if err != nil {
		return err
	}
	return app.initFromStore(mainKey, isConfigKey)
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
func (app *BaseApp) initFromStore(mainKey sdk.StoreKey, isConfigKey bool) error {
	// main store should exist.
	// TODO: we don't actually need the main store here
	var main sdk.KVStore
	if isConfigKey {
		main = app.cfcms.GetKVStore(mainKey)
	}else {
		main = app.cms.GetKVStore(mainKey)
	}
	if main == nil {
		return errors.New("baseapp expects MultiStore with 'main' KVStore")
	}
	// Needed for `gaiad export`, which inits from store but never calls initchain
	app.setCheckState(abci.Header{})

	app.Seal()

	return nil
}

// NewContext returns a new Context with the correct store, the given header, and nil txBytes.
func (app *BaseApp) NewContext(isCheckTx bool, header abci.Header) sdk.Context {
	if isCheckTx {
		return sdk.NewContext(app.checkState.ms, app.checkState.cfms, app.keys, app.configKeys, header, true, app.Logger)
	}
	return sdk.NewContext(app.deliverState.ms, app.checkState.cfms, app.keys, app.configKeys, header, false, app.Logger)
}

type state struct {
	ms  sdk.CacheMultiStore
	cfms  sdk.CacheMultiStore
	ctx sdk.Context
}

func (st *state) CacheMultiStore() sdk.CacheMultiStore {
	return st.ms.CacheMultiStore()
}

func (st *state) Context() sdk.Context {
	return st.ctx
}

func (app *BaseApp) setCheckState(header abci.Header) {
	ms := app.cms.CacheMultiStore()
	cfms := app.cfcms.CacheMultiStore()
	app.checkState = &state{
		ms:  ms,
		cfms: cfms,
		ctx: sdk.NewContext(ms, cfms, app.keys, app.configKeys, header, true, app.Logger),
	}
}

func (app *BaseApp) setDeliverState(header abci.Header) {
	ms := app.cms.(sdk.CacheMultiStore)
	cfms := app.cfcms.(sdk.CacheMultiStore)
	app.deliverState = &state{
		ms:  ms,
		cfms: cfms,
		ctx: sdk.NewContext(ms, cfms, app.keys, app.configKeys, header, false, app.Logger),
	}
}

//______________________________________________________________________________

// ABCI

// Implements ABCI
func (app *BaseApp) Info(req abci.RequestInfo) abci.ResponseInfo {
	lastCommitID := app.cfcms.LastCommitID()

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
	app.setDeliverState(abci.Header{ChainID: req.ChainId})
	app.setCheckState(abci.Header{ChainID: req.ChainId})

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
	ctx := sdk.NewContext(app.cms, app.cfcms, app.keys, app.configKeys, app.checkState.ctx.BlockHeader(), true, app.Logger)

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
		Data:      result.Data,
		Log:       result.Log,
		GasWanted: int64(result.GasWanted), // TODO: Should type accept unsigned ints?
		GasUsed:   int64(result.GasUsed),   // TODO: Should type accept unsigned ints?
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
		GasWanted: int64(result.GasWanted), // TODO: Should type accept unsigned ints?
		GasUsed:   int64(result.GasUsed),   // TODO: Should type accept unsigned ints?
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

func (app *BaseApp) runMsgs(ctx sdk.Context, tx sdk.Tx, mode runTxMode) sdk.Result {
	var code      sdk.CodeType
	var codespace sdk.CodespaceType

	txRoute := tx.Route()

	handler := app.router.Route(txRoute)

	if handler == nil {
		return sdk.ErrUnknownRequest("Unrecognized Msg type: " + txRoute).Result()
	}
	var msgResult sdk.Result

	fmt.Println("-------- enter handler ----------")
	fmt.Println(ctx.GasMeter().GasConsumed())

	// check 不实际执行
	if mode != runTxModeCheck {
		msgResult = handler(ctx, tx)
	}

	fmt.Println("-------- out handler ----------")
 	fmt.Println(ctx.GasMeter().GasConsumed())

	if !msgResult.IsOK() {
		code = msgResult.Code
		codespace = msgResult.Codespace
	}

	gasUsed := ctx.GasMeter().GasConsumed()
	if msgResult.GasUsed != 0 {
		gasUsed = msgResult.GasUsed
	}


	return sdk.Result{
		Code: 		code,
		Codespace: 	codespace,
		GasUsed:   	gasUsed,
		Log: 		strings.TrimSpace(msgResult.Log),
		Data: 		msgResult.Data,
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

	stdTx, ok := tx.(transaction.Transaction)
	gasWanted = stdTx.GetGas()
	if !ok {
		return transaction.ErrInvalidTx(sdk.CodespaceRoot, "tx must be StdTx").Result()
	}
	defer func() {

		if r := recover(); r != nil {
			switch rType := r.(type) {
			case sdk.ErrorOutOfGas:
				app.deferHandler(ctx, tx, true, 0)
				log := fmt.Sprintf("out of gas in location: %v", rType.Descriptor)
				result = sdk.ErrOutOfGas(log).Result()
				result.GasUsed = gasWanted
			default:
				app.deferHandler(ctx, tx, false, 0)
				log := fmt.Sprintf("recovered: %v\nstack:\n%v", r, string(debug.Stack()))
				result = sdk.ErrInternal(log).Result()
				result.GasUsed = ctx.GasMeter().GasConsumed()
			}
		} else if result.GasUsed != 0{
			app.deferHandler(ctx, tx, false, result.GasUsed)
		} else {
			app.deferHandler(ctx, tx, false, 0)
		}
		result.GasWanted = gasWanted
		//result.GasUsed = gasUsed
	}()

		if err := tx.ValidateBasic(); err != nil {
		return err.Result()
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

	runMsgCtx, msCache := app.cacheTxContext(ctx, txBytes)

	result = app.runMsgs(runMsgCtx, tx, mode)

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
	//app.deliverState.ms.Write()
	//commitID := app.cms.Commit()
	infoByte := app.cms.CommitStore()
	commitID := app.cfcms.CommitConfigStore(infoByte)


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
