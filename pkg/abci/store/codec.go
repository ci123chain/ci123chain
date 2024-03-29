package store

import (
	"github.com/ci123chain/ci123chain/pkg/abci/types"
)

// Import cosmos-sdk/types/store.go for convenience.
// nolint
type (
	PruningStrategy         = types.PruningStrategy
	Store                   = types.Store
	Committer               = types.Committer
	CommitStore             = types.CommitStore
	MultiStore              = types.MultiStore
	CacheMultiStore         = types.CacheMultiStore
	CommitMultiStore        = types.CommitMultiStore
	KVStore                 = types.KVStore
	KVPair                  = types.KVPair
	Iterator                = types.Iterator
	CacheKVStore            = types.CacheKVStore
	CommitKVStore           = types.CommitKVStore
	CacheWrapper            = types.CacheWrapper
	CacheWrap               = types.CacheWrap
	CommitID                = types.CommitID
	StoreKey                = types.StoreKey
	StoreType               = types.StoreType
	Queryable               = types.Queryable
	StoreWithInitialVersion = types.StoreWithInitialVersion
	TraceContext            = types.TraceContext
	Gas                     = types.Gas
	GasMeter                = types.GasMeter
	GasConfig               = types.GasConfig
)
