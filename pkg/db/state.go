package db

import (
	sdk "CI123Chain/pkg/abci/types"
)

type StateManager struct {
	key sdk.StoreKey
}

func NewStateManager(key sdk.StoreKey) *StateManager {
	return &StateManager{key: key}
}



