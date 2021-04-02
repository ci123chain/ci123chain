package types

import "github.com/pkg/errors"

var (
	ErrConnectionNotFound = errors.New("connection not found !")
	ErrInvalidConnectionState = errors.New("connection state invalid !")
)