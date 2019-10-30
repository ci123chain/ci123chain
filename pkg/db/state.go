package db

import (
	sdk "gitlab.oneitfarm.com/blockchain/ci123chain/pkg/abci/types"
)

type StateManager struct {
	key sdk.StoreKey
}

func NewStateManager(key sdk.StoreKey) *StateManager {
	return &StateManager{key: key}
}



