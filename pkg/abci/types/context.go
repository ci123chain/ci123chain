// nolint
package types

import (
	"context"
	"github.com/tendermint/tendermint/config"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/proto/tendermint/types"
)

/*
The intent of Context is for it to be an immutable object that can be
cloned and updated cheaply with WithValue() and passed forward to the
next decorator or handler. For example,

 func MsgHandler(ctx Context, tx Tx) Result {
 	...
 	ctx = ctx.WithValue(types, value)
 	...
 }
*/
type Context struct {
	context.Context
	pst *thePast
	gen int
	eventManager *EventManager
	// Don't add any other fields here,
	// it's probably not what you want to do.
}

// create a new context
func NewContext(ms MultiStore, header types.Header, isCheckTx bool, logger log.Logger) Context {
	c := Context{
		Context: context.Background(),
		pst:     newThePast(),
		gen:     0,
	}
	c = c.WithMultiStore(ms)
	c = c.WithBlockHeader(header)
	c = c.WithBlockHeight(header.Height)
	c = c.WithChainID(header.ChainID)
	c = c.WithIsCheckTx(isCheckTx)
	c = c.WithTxBytes(nil)
	c = c.WithLogger(logger)
	c = c.WithVoteInfos(nil)
	c = c.WithGasMeter(NewInfiniteGasMeter())
	c = c.WithEventManager(NewEventManager())
	return c
}

// is context nil
func (c Context) IsZero() bool {
	return c.Context == nil
}

//----------------------------------------
// Getting a value

// context value for the provided types
func (c Context) Value(key interface{}) interface{} {
	value := c.Context.Value(key)
	if cloner, ok := value.(cloner); ok {
		return cloner.Clone()
	}
	if message, ok := value.(proto.Message); ok {
		return proto.Clone(message)
	}
	return value
}

// KVStore fetches a KVStore from the MultiStore.
func (c Context) KVStore(key StoreKey) KVStore {
	return c.MultiStore().GetKVStore(key).Gas(c.GasMeter(), cachedKVGasConfig)
	//return c.MultiStore().GetKVStore(key).Prefix([]byte(key.Name() + "//")).Gas(c.GasMeter(), cachedKVGasConfig)
}

// KVStore fetches a KVStore from the MultiStore.
func (c Context) LoadLatestVersion() {
	err := c.MultiStore().(CommitMultiStore).LoadLatestVersion()
	if err != nil {
		panic(err)
	}
	return
}

// TransientStore fetches a TransientStore from the MultiStore.
func (c Context) TransientStore(key StoreKey) KVStore {
	return c.MultiStore().GetKVStore(key).Gas(c.GasMeter(), cachedTransientGasConfig)
}
func (c Context) EventManager() *EventManager { return c.eventManager }

//----------------------------------------
// With* (setting a value)

// nolint
func (c Context) WithValue(key interface{}, value interface{}) Context {
	return c.withValue(key, value)
}
func (c Context) WithCloner(key interface{}, value cloner) Context {
	return c.withValue(key, value)
}
func (c Context) WithCacheWrapper(key interface{}, value CacheWrapper) Context {
	return c.withValue(key, value)
}
func (c Context) WithProtoMsg(key interface{}, value proto.Message) Context {
	return c.withValue(key, value)
}
func (c Context) WithString(key interface{}, value string) Context {
	return c.withValue(key, value)
}
func (c Context) WithInt32(key interface{}, value int32) Context {
	return c.withValue(key, value)
}
func (c Context) WithUint32(key interface{}, value uint32) Context {
	return c.withValue(key, value)
}
func (c Context) WithUint64(key interface{}, value uint64) Context {
	return c.withValue(key, value)
}

func (c Context) withValue(key interface{}, value interface{}) Context {
	c.pst.bump(Op{
		gen:   c.gen + 1,
		key:   key,
		value: value,
	}) // increment version for all relatives.

	return Context{
		Context: context.WithValue(c.Context, key, value),
		pst:     c.pst,
		gen:     c.gen + 1,
	}
}


//----------------------------------------
// Values that require no types.

type contextKey int // local to the context module

const (
	contextKeyMultiStore contextKey = iota
	contextKeyBlockHeader
	contextKeyBlockHeight
	contextKeyConsensusParams
	contextKeyChainID
	contextKeyIsCheckTx
	contextKeyTxBytes
	contextKeyLogger
	contextKeyVoteInfos
	contextKeyGasMeter
	contextKeyMinimumFees
	contextKeyTxIndex
	contextKeyGasLimit
	contextKeyNonce
	contextKeyGasPrice
	contextKeyIsRootMsg
)

func (c Context) MultiStore() MultiStore {
	return c.Value(contextKeyMultiStore).(MultiStore)
}

func (c Context) BlockTime() time.Time        { return c.BlockHeader().Time }

func (c Context) BlockHeader() types.Header { return c.Value(contextKeyBlockHeader).(types.Header) }

func (c Context) BlockHeight() int64 { return c.Value(contextKeyBlockHeight).(int64) }

func (c Context) ConsensusParams() *abci.ConsensusParams {
	return c.Value(contextKeyConsensusParams).(*abci.ConsensusParams)
}

func (c Context) ChainID() string { return c.Value(contextKeyChainID).(string) }

func (c Context) TxBytes() []byte { return c.Value(contextKeyTxBytes).([]byte) }

func (c Context) Logger() log.Logger { return c.Value(contextKeyLogger).(log.Logger) }

func (c Context) VoteInfos() []abci.VoteInfo {
	return c.Value(contextKeyVoteInfos).([]abci.VoteInfo)
}

func (c Context) GasMeter() GasMeter { return c.Value(contextKeyGasMeter).(GasMeter) }

func (c Context) GasLimit() uint64 { return c.Value(contextKeyGasLimit).(uint64) }

func (c Context) Nonce() uint64 { return c.Value(contextKeyNonce).(uint64) }

func (c Context) IsCheckTx() bool { return c.Value(contextKeyIsCheckTx).(bool) }

func (c Context) IsRootMsg() bool {
	if c.Value(contextKeyIsRootMsg) != nil {
		return c.Value(contextKeyIsRootMsg).(bool)
	}
	return false
}

func (c Context) TxIndex() uint32 {
	return c.Value(contextKeyTxIndex).(uint32)
}

func (c Context) GasPriceConfig() *config.GasPriceConfig {
	return c.Value(contextKeyGasPrice).(*config.GasPriceConfig)
}

func (c Context) WithMultiStore(ms MultiStore) Context {
	return c.withValue(contextKeyMultiStore, ms)
}

func (c Context) WithBlockHeader(header types.Header) Context {
	//var _ proto.Message = &header // for cloning.
	return c.withValue(contextKeyBlockHeader, header)
}

func (c Context) WithBlockTime(newTime time.Time) Context {
	newHeader := c.BlockHeader()
	newHeader.Time = newTime
	return c.WithBlockHeader(newHeader)
}

func (c Context) WithProposer(addr []byte) Context {
	newHeader := c.BlockHeader()
	newHeader.ProposerAddress = addr
	return c.WithBlockHeader(newHeader)
}

func (c Context) WithBlockHeight(height int64) Context {
	newHeader := c.BlockHeader()
	newHeader.Height = height
	return c.withValue(contextKeyBlockHeight, height).withValue(contextKeyBlockHeader, newHeader)
}

func (c Context) WithConsensusParams(params *abci.ConsensusParams) Context {
	if params == nil {
		return c
	}

	// TODO: Do we need to handle invalid MaxGas values?
	return c.withValue(contextKeyConsensusParams, params).
		WithGasMeter(NewGasMeter(uint64(params.Block.MaxGas)))
}

func (c Context) WithEventManager(em *EventManager) Context {
	c.eventManager = em
	return c
}

func (c Context) WithChainID(chainID string) Context { return c.withValue(contextKeyChainID, chainID) }

func (c Context) WithTxBytes(txBytes []byte) Context { return c.withValue(contextKeyTxBytes, txBytes) }

func (c Context) WithLogger(logger log.Logger) Context { return c.withValue(contextKeyLogger, logger) }

func (c Context) WithVoteInfos(VoteInfos []abci.VoteInfo) Context {
	return c.withValue(contextKeyVoteInfos, VoteInfos)
}

func (c Context) WithGasMeter(meter GasMeter) Context { return c.withValue(contextKeyGasMeter, meter) }

func (c Context) WithIsRootMsg(isRootMsg bool) Context { return c.withValue(contextKeyIsRootMsg, isRootMsg) }

func (c Context) WithIsCheckTx(isCheckTx bool) Context {
	return c.withValue(contextKeyIsCheckTx, isCheckTx)
}

func (c Context) WithTxIndex(idx uint32) Context {
	return c.withValue(contextKeyTxIndex, idx)
}

func (c Context) WithGasLimit(gas uint64) Context {
	return c.withValue(contextKeyGasLimit, gas)
}

func (c Context) WithNonce(nonce uint64) Context {
	return c.withValue(contextKeyNonce, nonce)
}

func (c Context) WithGasPrice(gasPriceConfig *config.GasPriceConfig) Context {
	return c.withValue(contextKeyGasPrice, gasPriceConfig)
}

// Cache the multistore and return a new cached context. The cached context is
// written to the context when writeCache is called.
func (c Context) CacheContext() (cc Context, writeCache func()) {
	cms := c.MultiStore().CacheMultiStore()
	cc = c.WithMultiStore(cms)
	return cc, cms.Write
}

//----------------------------------------
// thePast

// Returns false if ver <= 0 || ver > len(c.pst.ops).
// The first operation is version 1.
func (c Context) GetOp(ver int64) (Op, bool) {
	return c.pst.getOp(ver)
}

//----------------------------------------
// Misc.

type cloner interface {
	Clone() interface{} // deep copy
}

// XXX add description
type Op struct {
	// types is always 'with'
	gen   int
	key   interface{}
	value interface{}
}

type thePast struct {
	mtx sync.RWMutex
	ver int
	ops []Op
}

func newThePast() *thePast {
	return &thePast{
		ver: 0,
		ops: nil,
	}
}

func (pst *thePast) bump(op Op) {
	pst.mtx.Lock()
	pst.ver++
	pst.ops = append(pst.ops, op)
	pst.mtx.Unlock()
}

func (pst *thePast) version() int {
	pst.mtx.RLock()
	defer pst.mtx.RUnlock()
	return pst.ver
}

// Returns false if ver <= 0 || ver > len(pst.ops).
// The first operation is version 1.
func (pst *thePast) getOp(ver int64) (Op, bool) {
	pst.mtx.RLock()
	defer pst.mtx.RUnlock()
	l := int64(len(pst.ops))
	if l < ver || ver <= 0 {
		return Op{}, false
	}
	return pst.ops[ver-1], true
}

type ContextKey string

const SdkContextKey ContextKey = "sdk-context"

func WrapSDKContext(ctx Context) context.Context {
	return context.WithValue(ctx.Context, SdkContextKey, ctx)
}

// UnwrapSDKContext retrieves a Context from a context.Context instance
// attached with WrapSDKContext. It panics if a Context was not properly
// attached
func UnwrapSDKContext(ctx context.Context) Context {
	return ctx.Value(SdkContextKey).(Context)
}
