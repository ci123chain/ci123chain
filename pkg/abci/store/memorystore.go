package store

import (
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	dbm "github.com/tendermint/tm-db"
)

var (
	_ types.KVStore   = (*MemoryStore)(nil)
	_ types.Committer = (*MemoryStore)(nil)
)

// Store implements an in-memory only KVStore. Entries are persisted between
// commits and thus between blocks. State in Memory store is not committed as part of app state but maintained privately by each node
type MemoryStore struct {
	dbStoreAdapter
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		dbStoreAdapter{dbm.NewMemDB()},
	}
}


func (s MemoryStore) LastCommitID() (id types.CommitID) { return }

func (s *MemoryStore) Commit() (id types.CommitID) { return }

func (s *MemoryStore) SetPruning(pruning PruningStrategy) {}

// GetPruning is a no-op as pruning options cannot be directly set on this store.
// They must be set on the root commit multi-store.
//func (s *MemoryStore) GetPruning() PruningStrategy { return PruningStrategy{} }