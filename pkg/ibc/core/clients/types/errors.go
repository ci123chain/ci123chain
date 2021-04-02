package types

import "errors"

var (
	ErrConsensusStateNotFound = errors.New("consensus state not found !")
	ErrClientFrozen = errors.New("client been frozen !")
	ErrClientNotFound = errors.New("err client not found !")
)
